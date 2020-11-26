package configmap

import (
	"sort"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/informers"
	"zmc.io/oasis/pkg/models/resources/v1alpha2"
	"zmc.io/oasis/pkg/server/params"
)

type configMapSearcher struct {
	informer informers.SharedInformerFactory
}

func NewConfigmapSearcher(informers informers.SharedInformerFactory) v1alpha2.Interface {
	return &configMapSearcher{informer: informers}
}

func (s *configMapSearcher) Get(namespace, name string) (interface{}, error) {
	return s.informer.Core().V1().ConfigMaps().Lister().ConfigMaps(namespace).Get(name)
}

func (s *configMapSearcher) match(match map[string]string, item *v1.ConfigMap) bool {
	for k, v := range match {
		if !v1alpha2.ObjectMetaExactlyMath(k, v, item.ObjectMeta) {
			return false
		}
	}
	return true
}

func (s *configMapSearcher) fuzzy(fuzzy map[string]string, item *v1.ConfigMap) bool {
	for k, v := range fuzzy {
		if !v1alpha2.ObjectMetaFuzzyMath(k, v, item.ObjectMeta) {
			return false
		}
	}
	return true
}

func (s *configMapSearcher) Search(namespace string, conditions *params.Conditions, orderBy string, reverse bool) ([]interface{}, error) {
	configMaps, err := s.informer.Core().V1().ConfigMaps().Lister().ConfigMaps(namespace).List(labels.Everything())

	if err != nil {
		return nil, err
	}

	result := make([]*v1.ConfigMap, 0)

	if len(conditions.Match) == 0 && len(conditions.Fuzzy) == 0 {
		result = configMaps
	} else {
		for _, item := range configMaps {
			if s.match(conditions.Match, item) && s.fuzzy(conditions.Fuzzy, item) {
				result = append(result, item)
			}
		}
	}
	sort.Slice(result, func(i, j int) bool {
		if reverse {
			i, j = j, i
		}
		return v1alpha2.ObjectMetaCompare(result[i].ObjectMeta, result[j].ObjectMeta, orderBy)
	})

	r := make([]interface{}, 0)
	for _, i := range result {
		r = append(r, i)
	}
	return r, nil
}

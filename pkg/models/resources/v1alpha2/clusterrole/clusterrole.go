package clusterrole

import (
	"sort"

	rbac "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/informers"
	"zmc.io/oasis/pkg/models/resources/v1alpha2"
	"zmc.io/oasis/pkg/server/params"
)

type clusterRoleSearcher struct {
	informer informers.SharedInformerFactory
}

func NewClusterRoleSearcher(informer informers.SharedInformerFactory) v1alpha2.Interface {
	return &clusterRoleSearcher{informer: informer}
}

func (s *clusterRoleSearcher) Get(namespace, name string) (interface{}, error) {
	return s.informer.Rbac().V1().ClusterRoles().Lister().Get(name)
}

func (*clusterRoleSearcher) match(match map[string]string, item *rbac.ClusterRole) bool {
	for k, v := range match {
		if !v1alpha2.ObjectMetaExactlyMath(k, v, item.ObjectMeta) {
			return false
		}
	}
	return true
}

func (s *clusterRoleSearcher) fuzzy(fuzzy map[string]string, item *rbac.ClusterRole) bool {
	for k, v := range fuzzy {
		if !v1alpha2.ObjectMetaFuzzyMath(k, v, item.ObjectMeta) {
			return false
		}
	}
	return true
}

func (s *clusterRoleSearcher) Search(namespace string, conditions *params.Conditions, orderBy string, reverse bool) ([]interface{}, error) {
	clusterRoles, err := s.informer.Rbac().V1().ClusterRoles().Lister().List(labels.Everything())

	if err != nil {
		return nil, err
	}

	result := make([]*rbac.ClusterRole, 0)

	if len(conditions.Match) == 0 && len(conditions.Fuzzy) == 0 {
		result = clusterRoles
	} else {
		for _, item := range clusterRoles {
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

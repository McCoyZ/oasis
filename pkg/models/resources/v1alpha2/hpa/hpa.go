package hpa

import (
	"sort"

	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/informers"
	"zmc.io/oasis/pkg/models/resources/v1alpha2"
	"zmc.io/oasis/pkg/server/params"
)

type hpaSearcher struct {
	informers informers.SharedInformerFactory
}

func NewHpaSearcher(informers informers.SharedInformerFactory) v1alpha2.Interface {
	return &hpaSearcher{informers: informers}
}

func (s *hpaSearcher) Get(namespace, name string) (interface{}, error) {
	return s.informers.Autoscaling().V2beta2().HorizontalPodAutoscalers().Lister().HorizontalPodAutoscalers(namespace).Get(name)
}

func hpaTargetMatch(item *autoscalingv2beta2.HorizontalPodAutoscaler, kind, name string) bool {
	return item.Spec.ScaleTargetRef.Kind == kind && item.Spec.ScaleTargetRef.Name == name
}

// exactly Match
func (*hpaSearcher) match(match map[string]string, item *autoscalingv2beta2.HorizontalPodAutoscaler) bool {
	for k, v := range match {
		switch k {
		case v1alpha2.TargetKind:
			fallthrough
		case v1alpha2.TargetName:
			kind := match[v1alpha2.TargetKind]
			name := match[v1alpha2.TargetName]
			if !hpaTargetMatch(item, kind, name) {
				return false
			}
		default:
			if !v1alpha2.ObjectMetaExactlyMath(k, v, item.ObjectMeta) {
				return false
			}
		}
	}
	return true
}

func (*hpaSearcher) fuzzy(fuzzy map[string]string, item *autoscalingv2beta2.HorizontalPodAutoscaler) bool {
	for k, v := range fuzzy {
		if !v1alpha2.ObjectMetaFuzzyMath(k, v, item.ObjectMeta) {
			return false
		}
	}
	return true
}

func (s *hpaSearcher) Search(namespace string, conditions *params.Conditions, orderBy string, reverse bool) ([]interface{}, error) {

	horizontalPodAutoscalers, err := s.informers.Autoscaling().V2beta2().HorizontalPodAutoscalers().Lister().HorizontalPodAutoscalers(namespace).List(labels.Everything())

	if err != nil {
		return nil, err
	}

	result := make([]*autoscalingv2beta2.HorizontalPodAutoscaler, 0)

	if len(conditions.Match) == 0 && len(conditions.Fuzzy) == 0 {
		result = horizontalPodAutoscalers
	} else {
		for _, item := range horizontalPodAutoscalers {
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

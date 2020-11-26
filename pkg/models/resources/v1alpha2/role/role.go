package role

import (
	"k8s.io/client-go/informers"
	"zmc.io/oasis/pkg/models/resources/v1alpha2"

	"sort"

	rbac "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/labels"
	"zmc.io/oasis/pkg/server/params"
)

type roleSearcher struct {
	informers informers.SharedInformerFactory
}

func NewRoleSearcher(informers informers.SharedInformerFactory) v1alpha2.Interface {
	return &roleSearcher{informers: informers}
}

func (s *roleSearcher) Get(namespace, name string) (interface{}, error) {
	return s.informers.Rbac().V1().Roles().Lister().Roles(namespace).Get(name)
}

func (*roleSearcher) match(match map[string]string, item *rbac.Role) bool {
	for k, v := range match {
		if !v1alpha2.ObjectMetaExactlyMath(k, v, item.ObjectMeta) {
			return false
		}
	}
	return true
}

func (*roleSearcher) fuzzy(fuzzy map[string]string, item *rbac.Role) bool {
	for k, v := range fuzzy {
		if !v1alpha2.ObjectMetaFuzzyMath(k, v, item.ObjectMeta) {
			return false
		}
	}
	return true
}

func (s *roleSearcher) Search(namespace string, conditions *params.Conditions, orderBy string, reverse bool) ([]interface{}, error) {
	roles, err := s.informers.Rbac().V1().Roles().Lister().Roles(namespace).List(labels.Everything())

	if err != nil {
		return nil, err
	}

	result := make([]*rbac.Role, 0)

	if len(conditions.Match) == 0 && len(conditions.Fuzzy) == 0 {
		result = roles
	} else {
		for _, item := range roles {
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

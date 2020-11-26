package ingress

import (
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/informers"
	"zmc.io/oasis/pkg/api"
	"zmc.io/oasis/pkg/apiserver/query"
	"zmc.io/oasis/pkg/models/resources/v1alpha3"
)

type ingressGetter struct {
	sharedInformers informers.SharedInformerFactory
}

func New(sharedInformers informers.SharedInformerFactory) v1alpha3.Interface {
	return &ingressGetter{sharedInformers: sharedInformers}
}

func (g *ingressGetter) Get(namespace, name string) (runtime.Object, error) {
	return g.sharedInformers.Extensions().V1beta1().Ingresses().Lister().Ingresses(namespace).Get(name)
}

func (g *ingressGetter) List(namespace string, query *query.Query) (*api.ListResult, error) {
	// first retrieves all deployments within given namespace
	ingresses, err := g.sharedInformers.Extensions().V1beta1().Ingresses().Lister().Ingresses(namespace).List(query.Selector())
	if err != nil {
		return nil, err
	}

	var result []runtime.Object
	for _, ingress := range ingresses {
		result = append(result, ingress)
	}

	return v1alpha3.DefaultList(result, query, g.compare, g.filter), nil
}

func (g *ingressGetter) compare(left runtime.Object, right runtime.Object, field query.Field) bool {

	leftIngress, ok := left.(*v1beta1.Ingress)
	if !ok {
		return false
	}

	rightIngress, ok := right.(*v1beta1.Ingress)
	if !ok {
		return false
	}

	switch field {
	case query.FieldUpdateTime:
		fallthrough
	default:
		return v1alpha3.DefaultObjectMetaCompare(leftIngress.ObjectMeta, rightIngress.ObjectMeta, field)
	}
}

func (g *ingressGetter) filter(object runtime.Object, filter query.Filter) bool {
	deployment, ok := object.(*v1beta1.Ingress)
	if !ok {
		return false
	}

	switch filter.Field {
	default:
		return v1alpha3.DefaultObjectMetaFilter(deployment.ObjectMeta, filter)
	}
}

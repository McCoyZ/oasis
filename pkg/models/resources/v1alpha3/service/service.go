package service

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/informers"
	"zmc.io/oasis/pkg/api"
	"zmc.io/oasis/pkg/apiserver/query"
	"zmc.io/oasis/pkg/models/resources/v1alpha3"
)

type servicesGetter struct {
	sharedInformers informers.SharedInformerFactory
}

func New(sharedInformers informers.SharedInformerFactory) v1alpha3.Interface {
	return &servicesGetter{sharedInformers: sharedInformers}
}

func (d *servicesGetter) Get(namespace, name string) (runtime.Object, error) {
	return d.sharedInformers.Core().V1().Services().Lister().Services(namespace).Get(name)
}

func (d *servicesGetter) List(namespace string, query *query.Query) (*api.ListResult, error) {
	// first retrieves all services within given namespace
	services, err := d.sharedInformers.Core().V1().Services().Lister().Services(namespace).List(query.Selector())
	if err != nil {
		return nil, err
	}

	var result []runtime.Object
	for _, deployment := range services {
		result = append(result, deployment)
	}

	return v1alpha3.DefaultList(result, query, d.compare, d.filter), nil
}

func (d *servicesGetter) compare(left runtime.Object, right runtime.Object, field query.Field) bool {

	leftService, ok := left.(*corev1.Service)
	if !ok {
		return false
	}

	rightService, ok := right.(*corev1.Service)
	if !ok {
		return false
	}

	return v1alpha3.DefaultObjectMetaCompare(leftService.ObjectMeta, rightService.ObjectMeta, field)
}

func (d *servicesGetter) filter(object runtime.Object, filter query.Filter) bool {
	service, ok := object.(*corev1.Service)
	if !ok {
		return false
	}

	return v1alpha3.DefaultObjectMetaFilter(service.ObjectMeta, filter)
}

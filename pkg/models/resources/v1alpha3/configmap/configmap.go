package configmap

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/informers"
	"zmc.io/oasis/pkg/api"
	"zmc.io/oasis/pkg/apiserver/query"
	"zmc.io/oasis/pkg/models/resources/v1alpha3"
)

type configmapsGetter struct {
	informer informers.SharedInformerFactory
}

func New(sharedInformers informers.SharedInformerFactory) v1alpha3.Interface {
	return &configmapsGetter{informer: sharedInformers}
}

func (d *configmapsGetter) Get(namespace, name string) (runtime.Object, error) {
	return d.informer.Core().V1().ConfigMaps().Lister().ConfigMaps(namespace).Get(name)
}

func (d *configmapsGetter) List(namespace string, query *query.Query) (*api.ListResult, error) {
	configmaps, err := d.informer.Core().V1().ConfigMaps().Lister().ConfigMaps(namespace).List(query.Selector())
	if err != nil {
		return nil, err
	}

	var result []runtime.Object
	for _, configmap := range configmaps {
		result = append(result, configmap)
	}

	return v1alpha3.DefaultList(result, query, d.compare, d.filter), nil
}

func (d *configmapsGetter) compare(left runtime.Object, right runtime.Object, field query.Field) bool {

	leftCM, ok := left.(*corev1.ConfigMap)
	if !ok {
		return false
	}

	rightCM, ok := right.(*corev1.ConfigMap)
	if !ok {
		return false
	}

	return v1alpha3.DefaultObjectMetaCompare(leftCM.ObjectMeta, rightCM.ObjectMeta, field)
}

func (d *configmapsGetter) filter(object runtime.Object, filter query.Filter) bool {
	configMap, ok := object.(*corev1.ConfigMap)
	if !ok {
		return false
	}

	return v1alpha3.DefaultObjectMetaFilter(configMap.ObjectMeta, filter)
}

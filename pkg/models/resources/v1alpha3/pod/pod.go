package pod

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/informers"
	"zmc.io/oasis/pkg/api"
	"zmc.io/oasis/pkg/apiserver/query"
	"zmc.io/oasis/pkg/models/resources/v1alpha3"
)

const (
	filedNameName    = "nodeName"
	filedPVCName     = "pvcName"
	filedServiceName = "serviceName"
)

type podsGetter struct {
	informer informers.SharedInformerFactory
}

func New(sharedInformers informers.SharedInformerFactory) v1alpha3.Interface {
	return &podsGetter{informer: sharedInformers}
}

func (p *podsGetter) Get(namespace, name string) (runtime.Object, error) {
	return p.informer.Core().V1().Pods().Lister().Pods(namespace).Get(name)
}

func (p *podsGetter) List(namespace string, query *query.Query) (*api.ListResult, error) {

	pods, err := p.informer.Core().V1().Pods().Lister().Pods(namespace).List(query.Selector())
	if err != nil {
		return nil, err
	}

	var result []runtime.Object
	for _, pod := range pods {
		result = append(result, pod)
	}

	return v1alpha3.DefaultList(result, query, p.compare, p.filter), nil
}

func (p *podsGetter) compare(left runtime.Object, right runtime.Object, field query.Field) bool {

	leftPod, ok := left.(*corev1.Pod)
	if !ok {
		return false
	}

	rightPod, ok := right.(*corev1.Pod)
	if !ok {
		return false
	}

	return v1alpha3.DefaultObjectMetaCompare(leftPod.ObjectMeta, rightPod.ObjectMeta, field)
}

func (p *podsGetter) filter(object runtime.Object, filter query.Filter) bool {
	pod, ok := object.(*corev1.Pod)

	if !ok {
		return false
	}
	switch filter.Field {
	case filedNameName:
		return pod.Spec.NodeName == string(filter.Value)
	case filedPVCName:
		return p.podBindPVC(pod, string(filter.Value))
	case filedServiceName:
		return p.podBelongToService(pod, string(filter.Value))
	default:
		return v1alpha3.DefaultObjectMetaFilter(pod.ObjectMeta, filter)
	}
}

func (p *podsGetter) podBindPVC(item *corev1.Pod, pvcName string) bool {
	for _, v := range item.Spec.Volumes {
		if v.VolumeSource.PersistentVolumeClaim != nil &&
			v.VolumeSource.PersistentVolumeClaim.ClaimName == pvcName {
			return true
		}
	}
	return false
}

func (p *podsGetter) podBelongToService(item *corev1.Pod, serviceName string) bool {
	service, err := p.informer.Core().V1().Services().Lister().Services(item.Namespace).Get(serviceName)
	if err != nil {
		return false
	}
	selector := labels.Set(service.Spec.Selector).AsSelectorPreValidated()
	if selector.Empty() || !selector.Matches(labels.Set(item.Labels)) {
		return false
	}
	return true
}

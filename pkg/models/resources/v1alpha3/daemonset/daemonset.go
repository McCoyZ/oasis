package daemonset

import (
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/informers"
	"zmc.io/oasis/pkg/api"
	"zmc.io/oasis/pkg/apiserver/query"
	"zmc.io/oasis/pkg/models/resources/v1alpha3"
)

const (
	statusStopped  = "stopped"
	statusRunning  = "running"
	statusUpdating = "updating"
)

type daemonSetGetter struct {
	sharedInformers informers.SharedInformerFactory
}

func New(sharedInformers informers.SharedInformerFactory) v1alpha3.Interface {
	return &daemonSetGetter{sharedInformers: sharedInformers}
}

func (d *daemonSetGetter) Get(namespace, name string) (runtime.Object, error) {
	return d.sharedInformers.Apps().V1().DaemonSets().Lister().DaemonSets(namespace).Get(name)
}

func (d *daemonSetGetter) List(namespace string, query *query.Query) (*api.ListResult, error) {
	// first retrieves all daemonSets within given namespace
	daemonSets, err := d.sharedInformers.Apps().V1().DaemonSets().Lister().DaemonSets(namespace).List(query.Selector())
	if err != nil {
		return nil, err
	}

	var result []runtime.Object
	for _, daemonSet := range daemonSets {
		result = append(result, daemonSet)
	}

	return v1alpha3.DefaultList(result, query, d.compare, d.filter), nil
}

func (d *daemonSetGetter) compare(left runtime.Object, right runtime.Object, field query.Field) bool {

	leftDaemonSet, ok := left.(*appsv1.DaemonSet)
	if !ok {
		return false
	}

	rightDaemonSet, ok := right.(*appsv1.DaemonSet)
	if !ok {
		return false
	}

	return v1alpha3.DefaultObjectMetaCompare(leftDaemonSet.ObjectMeta, rightDaemonSet.ObjectMeta, field)
}

func (d *daemonSetGetter) filter(object runtime.Object, filter query.Filter) bool {
	daemonSet, ok := object.(*appsv1.DaemonSet)
	if !ok {
		return false
	}
	switch filter.Field {
	case query.FieldStatus:
		return strings.Compare(daemonsetStatus(&daemonSet.Status), string(filter.Value)) == 0
	default:
		return v1alpha3.DefaultObjectMetaFilter(daemonSet.ObjectMeta, filter)
	}
}

func daemonsetStatus(status *appsv1.DaemonSetStatus) string {
	if status.DesiredNumberScheduled == 0 && status.NumberReady == 0 {
		return statusStopped
	} else if status.DesiredNumberScheduled == status.NumberReady {
		return statusRunning
	} else {
		return statusUpdating
	}
}

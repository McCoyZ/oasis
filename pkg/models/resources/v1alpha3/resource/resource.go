package resource

import (
	"errors"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"zmc.io/oasis/pkg/api"
	"zmc.io/oasis/pkg/apiserver/query"
	"zmc.io/oasis/pkg/informers"
	"zmc.io/oasis/pkg/models/resources/v1alpha3"
	"zmc.io/oasis/pkg/models/resources/v1alpha3/configmap"
	"zmc.io/oasis/pkg/models/resources/v1alpha3/daemonset"
	"zmc.io/oasis/pkg/models/resources/v1alpha3/deployment"
	"zmc.io/oasis/pkg/models/resources/v1alpha3/ingress"
	"zmc.io/oasis/pkg/models/resources/v1alpha3/job"
	"zmc.io/oasis/pkg/models/resources/v1alpha3/namespace"
	"zmc.io/oasis/pkg/models/resources/v1alpha3/networkpolicy"
	"zmc.io/oasis/pkg/models/resources/v1alpha3/node"
	"zmc.io/oasis/pkg/models/resources/v1alpha3/pod"
	"zmc.io/oasis/pkg/models/resources/v1alpha3/service"
	"zmc.io/oasis/pkg/models/resources/v1alpha3/statefulset"
)

var ErrResourceNotSupported = errors.New("resource is not supported")

type ResourceGetter struct {
	getters map[schema.GroupVersionResource]v1alpha3.Interface
}

func NewResourceGetter(factory informers.InformerFactory) *ResourceGetter {
	getters := make(map[schema.GroupVersionResource]v1alpha3.Interface)

	getters[schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}] = deployment.New(factory.KubernetesSharedInformerFactory())
	getters[schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "daemonsets"}] = daemonset.New(factory.KubernetesSharedInformerFactory())
	getters[schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "statefulsets"}] = statefulset.New(factory.KubernetesSharedInformerFactory())
	getters[schema.GroupVersionResource{Group: "", Version: "v1", Resource: "services"}] = service.New(factory.KubernetesSharedInformerFactory())
	getters[schema.GroupVersionResource{Group: "", Version: "v1", Resource: "namespaces"}] = namespace.New(factory.KubernetesSharedInformerFactory())
	getters[schema.GroupVersionResource{Group: "", Version: "v1", Resource: "configmaps"}] = configmap.New(factory.KubernetesSharedInformerFactory())
	getters[schema.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}] = pod.New(factory.KubernetesSharedInformerFactory())
	getters[schema.GroupVersionResource{Group: "", Version: "v1", Resource: "nodes"}] = node.New(factory.KubernetesSharedInformerFactory())
	getters[schema.GroupVersionResource{Group: "extensions", Version: "v1beta1", Resource: "ingresses"}] = ingress.New(factory.KubernetesSharedInformerFactory())
	getters[schema.GroupVersionResource{Group: "networking.k8s.io", Version: "v1", Resource: "networkpolicies"}] = networkpolicy.New(factory.KubernetesSharedInformerFactory())
	getters[schema.GroupVersionResource{Group: "batch", Version: "v1", Resource: "jobs"}] = job.New(factory.KubernetesSharedInformerFactory())

	return &ResourceGetter{
		getters: getters,
	}
}

// tryResource will retrieve a getter with resource name, it doesn't guarantee find resource with correct group version
// need to refactor this use schema.GroupVersionResource
func (r *ResourceGetter) tryResource(resource string) v1alpha3.Interface {
	for k, v := range r.getters {
		if k.Resource == resource {
			return v
		}
	}
	return nil
}

func (r *ResourceGetter) Get(resource, namespace, name string) (runtime.Object, error) {
	getter := r.tryResource(resource)
	if getter == nil {
		return nil, ErrResourceNotSupported
	}
	return getter.Get(namespace, name)
}

func (r *ResourceGetter) List(resource, namespace string, query *query.Query) (*api.ListResult, error) {
	getter := r.tryResource(resource)
	if getter == nil {
		return nil, ErrResourceNotSupported
	}
	return getter.List(namespace, query)
}

package resource

import (
	"fmt"

	"k8s.io/klog"
	"zmc.io/oasis/pkg/informers"
	"zmc.io/oasis/pkg/models"
	"zmc.io/oasis/pkg/models/resources/v1alpha2"

	"zmc.io/oasis/pkg/models/resources/v1alpha2/clusterrole"
	"zmc.io/oasis/pkg/models/resources/v1alpha2/configmap"
	"zmc.io/oasis/pkg/models/resources/v1alpha2/cronjob"
	"zmc.io/oasis/pkg/models/resources/v1alpha2/daemonset"
	"zmc.io/oasis/pkg/models/resources/v1alpha2/deployment"
	"zmc.io/oasis/pkg/models/resources/v1alpha2/hpa"
	"zmc.io/oasis/pkg/models/resources/v1alpha2/ingress"
	"zmc.io/oasis/pkg/models/resources/v1alpha2/job"
	"zmc.io/oasis/pkg/models/resources/v1alpha2/namespace"
	"zmc.io/oasis/pkg/models/resources/v1alpha2/node"
	"zmc.io/oasis/pkg/models/resources/v1alpha2/pod"
	"zmc.io/oasis/pkg/models/resources/v1alpha2/role"
	"zmc.io/oasis/pkg/models/resources/v1alpha2/secret"
	"zmc.io/oasis/pkg/models/resources/v1alpha2/service"
	"zmc.io/oasis/pkg/models/resources/v1alpha2/statefulset"
	"zmc.io/oasis/pkg/server/params"
	"zmc.io/oasis/pkg/utils/sliceutil"
)

type ResourceGetter struct {
	resourcesGetters map[string]v1alpha2.Interface
}

func (r ResourceGetter) Add(resource string, getter v1alpha2.Interface) {
	if r.resourcesGetters == nil {
		r.resourcesGetters = make(map[string]v1alpha2.Interface)
	}
	r.resourcesGetters[resource] = getter
}

func NewResourceGetter(factory informers.InformerFactory) *ResourceGetter {
	resourceGetters := make(map[string]v1alpha2.Interface)

	resourceGetters[v1alpha2.ConfigMaps] = configmap.NewConfigmapSearcher(factory.KubernetesSharedInformerFactory())
	resourceGetters[v1alpha2.CronJobs] = cronjob.NewCronJobSearcher(factory.KubernetesSharedInformerFactory())
	resourceGetters[v1alpha2.DaemonSets] = daemonset.NewDaemonSetSearcher(factory.KubernetesSharedInformerFactory())
	resourceGetters[v1alpha2.Deployments] = deployment.NewDeploymentSetSearcher(factory.KubernetesSharedInformerFactory())
	resourceGetters[v1alpha2.Ingresses] = ingress.NewIngressSearcher(factory.KubernetesSharedInformerFactory())
	resourceGetters[v1alpha2.Jobs] = job.NewJobSearcher(factory.KubernetesSharedInformerFactory())
	resourceGetters[v1alpha2.Secrets] = secret.NewSecretSearcher(factory.KubernetesSharedInformerFactory())
	resourceGetters[v1alpha2.Services] = service.NewServiceSearcher(factory.KubernetesSharedInformerFactory())
	resourceGetters[v1alpha2.StatefulSets] = statefulset.NewStatefulSetSearcher(factory.KubernetesSharedInformerFactory())
	resourceGetters[v1alpha2.Pods] = pod.NewPodSearcher(factory.KubernetesSharedInformerFactory())
	resourceGetters[v1alpha2.Roles] = role.NewRoleSearcher(factory.KubernetesSharedInformerFactory())

	resourceGetters[v1alpha2.Nodes] = node.NewNodeSearcher(factory.KubernetesSharedInformerFactory())
	resourceGetters[v1alpha2.Namespaces] = namespace.NewNamespaceSearcher(factory.KubernetesSharedInformerFactory())
	resourceGetters[v1alpha2.ClusterRoles] = clusterrole.NewClusterRoleSearcher(factory.KubernetesSharedInformerFactory())
	resourceGetters[v1alpha2.HorizontalPodAutoscalers] = hpa.NewHpaSearcher(factory.KubernetesSharedInformerFactory())

	return &ResourceGetter{resourcesGetters: resourceGetters}

}

var (
	//injector         = v1alpha2.extraAnnotationInjector{}
	clusterResources = []string{v1alpha2.Nodes, v1alpha2.Namespaces, v1alpha2.ClusterRoles}
)

func (r *ResourceGetter) GetResource(namespace, resource, name string) (interface{}, error) {
	if searcher, ok := r.resourcesGetters[resource]; ok {
		resource, err := searcher.Get(namespace, name)
		if err != nil {
			klog.Errorf("resource %s.%s.%s not found: %s", namespace, resource, name, err)
			return nil, err
		}
		return resource, nil
	}
	return nil, fmt.Errorf("resource %s.%s.%s not found", namespace, resource, name)
}

func (r *ResourceGetter) ListResources(namespace, resource string, conditions *params.Conditions, orderBy string, reverse bool, limit, offset int) (*models.PageableResponse, error) {
	items := make([]interface{}, 0)
	var err error
	var result []interface{}

	// none namespace resource
	if namespace != "" && sliceutil.HasString(clusterResources, resource) {
		err = fmt.Errorf("resource %s is not supported", resource)
		klog.Errorln(err)
		return nil, err
	}

	if searcher, ok := r.resourcesGetters[resource]; ok {
		result, err = searcher.Search(namespace, conditions, orderBy, reverse)
	} else {
		err = fmt.Errorf("resource %s is not supported", resource)
		klog.Errorln(err)
		return nil, err
	}

	if err != nil {
		klog.Errorln(err)
		return nil, err
	}

	if limit == -1 || limit+offset > len(result) {
		limit = len(result) - offset
	}

	items = result[offset : offset+limit]

	return &models.PageableResponse{TotalCount: len(result), Items: items}, nil
}

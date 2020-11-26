package v1alpha2

import (
	"strings"

	"github.com/emicklei/go-restful"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog"
	"zmc.io/oasis/pkg/api"
	"zmc.io/oasis/pkg/informers"
	"zmc.io/oasis/pkg/models/quotas"

	"zmc.io/oasis/pkg/models/resources/v1alpha2"
	"zmc.io/oasis/pkg/models/resources/v1alpha2/resource"
	"zmc.io/oasis/pkg/server/params"
)

type resourceHandler struct {
	resourcesGetter     *resource.ResourceGetter
	resourceQuotaGetter quotas.ResourceQuotaGetter
}

func newResourceHandler(k8sClient kubernetes.Interface, factory informers.InformerFactory, masterURL string) *resourceHandler {

	return &resourceHandler{
		resourcesGetter:     resource.NewResourceGetter(factory),
		resourceQuotaGetter: quotas.NewResourceQuotaGetter(factory.KubernetesSharedInformerFactory()),
	}
}

func (r *resourceHandler) handleGetNamespacedResources(request *restful.Request, response *restful.Response) {
	r.handleListNamespaceResources(request, response)
}

func (r *resourceHandler) handleListNamespaceResources(request *restful.Request, response *restful.Response) {
	namespace := request.PathParameter("namespace")
	resource := request.PathParameter("resources")
	orderBy := params.GetStringValueWithDefault(request, params.OrderByParam, v1alpha2.CreateTime)
	limit, offset := params.ParsePaging(request)
	reverse := params.GetBoolValueWithDefault(request, params.ReverseParam, false)
	conditions, err := params.ParseConditions(request)

	if err != nil {
		klog.Error(err)
		api.HandleBadRequest(response, request, err)
		return
	}

	result, err := r.resourcesGetter.ListResources(namespace, resource, conditions, orderBy, reverse, limit, offset)

	if err != nil {
		klog.Error(err)
		api.HandleInternalError(response, nil, err)
		return
	}

	response.WriteEntity(result)
}

func (r *resourceHandler) handleGetClusterQuotas(_ *restful.Request, response *restful.Response) {
	result, err := r.resourceQuotaGetter.GetClusterQuota()
	if err != nil {
		api.HandleInternalError(response, nil, err)
		return
	}

	response.WriteAsJson(result)
}

func (r *resourceHandler) handleGetNamespaceQuotas(request *restful.Request, response *restful.Response) {
	namespace := request.PathParameter("namespace")
	quota, err := r.resourceQuotaGetter.GetNamespaceQuota(namespace)

	if err != nil {
		api.HandleInternalError(response, nil, err)
		return
	}

	response.WriteAsJson(quota)
}

func (r *resourceHandler) handleGetNamespacedAbnormalWorkloads(request *restful.Request, response *restful.Response) {
	namespace := request.PathParameter("namespace")

	result := api.Workloads{
		Namespace: namespace,
		Count:     make(map[string]int),
	}

	for _, workloadType := range []string{api.ResourceKindDeployment, api.ResourceKindStatefulSet, api.ResourceKindDaemonSet, api.ResourceKindJob, api.ResourceKindPersistentVolumeClaim} {
		var notReadyStatus string

		switch workloadType {
		case api.ResourceKindPersistentVolumeClaim:
			notReadyStatus = strings.Join([]string{v1alpha2.StatusPending, v1alpha2.StatusLost}, "|")
		case api.ResourceKindJob:
			notReadyStatus = v1alpha2.StatusFailed
		default:
			notReadyStatus = v1alpha2.StatusUpdating
		}

		res, err := r.resourcesGetter.ListResources(namespace, workloadType, &params.Conditions{Match: map[string]string{v1alpha2.Status: notReadyStatus}}, "", false, -1, 0)
		if err != nil {
			api.HandleInternalError(response, nil, err)
		}

		result.Count[workloadType] = len(res.Items)
	}

	response.WriteAsJson(result)

}

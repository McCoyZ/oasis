package v1alpha2

import (
	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	"zmc.io/oasis/pkg/api"
	"zmc.io/oasis/pkg/apiserver/runtime"
	"zmc.io/oasis/pkg/constants"
	"zmc.io/oasis/pkg/informers"
	"zmc.io/oasis/pkg/models"

	"net/http"

	"zmc.io/oasis/pkg/server/params"
)

const (
	GroupName = "resources"
)

var GroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1alpha2"}

func AddToContainer(c *restful.Container, k8sClient kubernetes.Interface, factory informers.InformerFactory, masterURL string) error {
	webservice := runtime.NewWebService(GroupVersion)
	handler := newResourceHandler(k8sClient, factory, masterURL)

	webservice.Route(webservice.GET("/namespaces/{namespace}/{resources}").
		To(handler.handleListNamespaceResources).
		Deprecate().
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.NamespaceResourcesTag}).
		Doc("Namespace level resource query").
		Param(webservice.PathParameter("namespace", "the name of the project")).
		Param(webservice.PathParameter("resources", "namespace level resource type, e.g. pods,jobs,configmaps,services.")).
		Param(webservice.QueryParameter(params.ConditionsParam, "query conditions,connect multiple conditions with commas, equal symbol for exact query, wave symbol for fuzzy query e.g. name~a").
			Required(false).
			DataFormat("key=%s,key~%s")).
		Param(webservice.QueryParameter(params.PagingParam, "paging query, e.g. limit=100,page=1").
			Required(false).
			DataFormat("limit=%d,page=%d").
			DefaultValue("limit=10,page=1")).
		Param(webservice.QueryParameter(params.ReverseParam, "sort parameters, e.g. reverse=true")).
		Param(webservice.QueryParameter(params.OrderByParam, "sort parameters, e.g. orderBy=createTime")).
		Returns(http.StatusOK, api.StatusOK, models.PageableResponse{}))

	webservice.Route(webservice.GET("/{resources}").
		To(handler.handleListNamespaceResources).
		Deprecate().
		Returns(http.StatusOK, api.StatusOK, models.PageableResponse{}).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.ClusterResourcesTag}).
		Doc("Cluster level resources").
		Param(webservice.PathParameter("resources", "cluster level resource type, e.g. nodes,workspaces,storageclasses,clusterrole.")).
		Param(webservice.QueryParameter(params.ConditionsParam, "query conditions, connect multiple conditions with commas, equal symbol for exact query, wave symbol for fuzzy query e.g. name~a").
			Required(false).
			DataFormat("key=value,key~value").
			DefaultValue("")).
		Param(webservice.QueryParameter(params.PagingParam, "paging query, e.g. limit=100,page=1").
			Required(false).
			DataFormat("limit=%d,page=%d").
			DefaultValue("limit=10,page=1")).
		Param(webservice.QueryParameter(params.ReverseParam, "sort parameters, e.g. reverse=true")).
		Param(webservice.QueryParameter(params.OrderByParam, "sort parameters, e.g. orderBy=createTime")))

	webservice.Route(webservice.GET("/quotas").
		To(handler.handleGetClusterQuotas).
		Doc("get whole cluster's resource usage").
		Returns(http.StatusOK, api.StatusOK, api.ResourceQuota{}).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.ClusterResourcesTag}))

	webservice.Route(webservice.GET("/namespaces/{namespace}/quotas").
		Doc("get specified namespace's resource quota and usage").
		Param(webservice.PathParameter("namespace", "the name of the project")).
		Returns(http.StatusOK, api.StatusOK, api.ResourceQuota{}).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.NamespaceResourcesTag}).
		To(handler.handleGetNamespaceQuotas))

	webservice.Route(webservice.GET("/abnormalworkloads").
		Doc("get abnormal workloads' count of whole cluster").
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.ClusterResourcesTag}).
		Returns(http.StatusOK, api.StatusOK, api.Workloads{}).
		To(handler.handleGetNamespacedAbnormalWorkloads))
	webservice.Route(webservice.GET("/namespaces/{namespace}/abnormalworkloads").
		Doc("get abnormal workloads' count of specified namespace").
		Param(webservice.PathParameter("namespace", "the name of the project")).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.NamespaceResourcesTag}).
		Returns(http.StatusOK, api.StatusOK, api.Workloads{}).
		To(handler.handleGetNamespacedAbnormalWorkloads))

	c.Add(webservice)

	return nil
}

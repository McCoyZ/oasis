package v1alpha3

import (
	"strings"

	"github.com/emicklei/go-restful"
	"k8s.io/klog"
	"zmc.io/oasis/pkg/api"
	"zmc.io/oasis/pkg/apiserver/query"
	"zmc.io/oasis/pkg/informers"
	"zmc.io/oasis/pkg/models/resources/v1alpha2"
	resourcev1alpha2 "zmc.io/oasis/pkg/models/resources/v1alpha2/resource"
	"zmc.io/oasis/pkg/models/resources/v1alpha3/resource"
	"zmc.io/oasis/pkg/server/params"
)

type Handler struct {
	resourceGetterV1alpha3  *resource.ResourceGetter
	resourcesGetterV1alpha2 *resourcev1alpha2.ResourceGetter
}

func New(factory informers.InformerFactory) *Handler {
	return &Handler{
		resourceGetterV1alpha3:  resource.NewResourceGetter(factory),
		resourcesGetterV1alpha2: resourcev1alpha2.NewResourceGetter(factory),
	}
}

func (h *Handler) handleGetResources(request *restful.Request, response *restful.Response) {
	namespace := request.PathParameter("namespace")
	resourceType := request.PathParameter("resources")
	name := request.PathParameter("name")

	result, err := h.resourceGetterV1alpha3.Get(resourceType, namespace, name)
	if err == nil {
		response.WriteEntity(result)
		return
	}

	if err != resource.ErrResourceNotSupported {
		klog.Error(err, resourceType)
		api.HandleInternalError(response, nil, err)
		return
	}

	// fallback to v1alpha2
	resultV1alpha2, err := h.resourcesGetterV1alpha2.GetResource(namespace, resourceType, name)

	if err != nil {
		klog.Error(err)
		api.HandleInternalError(response, nil, err)
		return
	}

	response.WriteEntity(resultV1alpha2)

}

// handleListResources retrieves resources
func (h *Handler) handleListResources(request *restful.Request, response *restful.Response) {
	query := query.ParseQueryParameter(request)
	resourceType := request.PathParameter("resources")
	namespace := request.PathParameter("namespace")

	result, err := h.resourceGetterV1alpha3.List(resourceType, namespace, query)

	if err == nil {
		response.WriteEntity(result)
		return
	}

	if err != resource.ErrResourceNotSupported {
		klog.Error(err, resourceType)
		api.HandleInternalError(response, nil, err)
		return
	}

	// fallback to v1alpha2
	result, err = h.fallback(resourceType, namespace, query)

	if err != nil {
		klog.Error(err)
		api.HandleInternalError(response, nil, err)
		return
	}

	response.WriteEntity(result)
}

func (h *Handler) fallback(resourceType string, namespace string, q *query.Query) (*api.ListResult, error) {
	orderBy := string(q.SortBy)
	limit, offset := q.Pagination.Limit, q.Pagination.Offset
	reverse := !q.Ascending
	conditions := &params.Conditions{Match: make(map[string]string, 0), Fuzzy: make(map[string]string, 0)}
	for field, value := range q.Filters {
		switch field {
		case query.FieldName:
			conditions.Fuzzy[v1alpha2.Name] = string(value)
			break
		case query.FieldNames:
			conditions.Match[v1alpha2.Name] = string(value)
			break
		case query.FieldCreationTimeStamp:
			conditions.Match[v1alpha2.CreateTime] = string(value)
			break
		case query.FieldLastUpdateTimestamp:
			conditions.Match[v1alpha2.UpdateTime] = string(value)
			break
		case query.FieldLabel:
			values := strings.SplitN(string(value), ":", 2)
			if len(values) == 2 {
				conditions.Match[values[0]] = values[1]
			} else {
				conditions.Match[v1alpha2.Label] = values[0]
			}
			break
		case query.FieldAnnotation:
			values := strings.SplitN(string(value), ":", 2)
			if len(values) == 2 {
				conditions.Match[v1alpha2.Annotation] = values[1]
			} else {
				conditions.Match[v1alpha2.Annotation] = values[0]
			}
			break
		case query.FieldStatus:
			conditions.Match[v1alpha2.Status] = string(value)
			break
		default:
			conditions.Match[string(field)] = string(value)
			break
		}
	}

	result, err := h.resourcesGetterV1alpha2.ListResources(namespace, resourceType, conditions, orderBy, reverse, limit, offset)

	if err != nil {
		klog.Error(err)
		return nil, err
	}

	return &api.ListResult{
		Items:      result.Items,
		TotalItems: result.TotalCount,
	}, nil
}

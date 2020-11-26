package v1alpha2

import (
	"github.com/emicklei/go-restful"
	"k8s.io/apimachinery/pkg/runtime/schema"
	kuberconfig "zmc.io/oasis/pkg/apiserver/config"
	"zmc.io/oasis/pkg/apiserver/runtime"
)

const (
	GroupName = ""
)

var GroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1alpha2"}

func AddToContainer(c *restful.Container, config *kuberconfig.Config) error {
	webservice := runtime.NewWebService(GroupVersion)

	webservice.Route(webservice.GET("/configs/configz").
		Doc("Information about the server configuration").
		To(func(request *restful.Request, response *restful.Response) {
			// response.WriteAsJson(config.ToMap())
			response.WriteAsJson(config)
		}))

	c.Add(webservice)
	return nil
}

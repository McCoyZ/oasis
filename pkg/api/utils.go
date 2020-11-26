package api

import (
	"net/http"
	"runtime"
	"strings"

	"github.com/emicklei/go-restful"
	"k8s.io/klog"
)

// Avoid emitting errors that look like valid HTML. Quotes are okay.
var sanitizer = strings.NewReplacer(`&`, "&amp;", `<`, "&lt;", `>`, "&gt;")

func HandleInternalError(response *restful.Response, req *restful.Request, err error) {
	_, fn, line, _ := runtime.Caller(1)
	klog.Errorf("%s:%d %v", fn, line, err)
	http.Error(response, sanitizer.Replace(err.Error()), http.StatusInternalServerError)
}

// HandleBadRequest writes http.StatusBadRequest and log error
func HandleBadRequest(response *restful.Response, req *restful.Request, err error) {
	_, fn, line, _ := runtime.Caller(1)
	klog.Errorf("%s:%d %v", fn, line, err)
	http.Error(response, sanitizer.Replace(err.Error()), http.StatusBadRequest)
}

func HandleNotFound(response *restful.Response, req *restful.Request, err error) {
	_, fn, line, _ := runtime.Caller(1)
	klog.Errorf("%s:%d %v", fn, line, err)
	http.Error(response, sanitizer.Replace(err.Error()), http.StatusNotFound)
}

func HandleForbidden(response *restful.Response, req *restful.Request, err error) {
	_, fn, line, _ := runtime.Caller(1)
	klog.Errorf("%s:%d %v", fn, line, err)
	http.Error(response, sanitizer.Replace(err.Error()), http.StatusForbidden)
}

func HandleConflict(response *restful.Response, req *restful.Request, err error) {
	_, fn, line, _ := runtime.Caller(1)
	klog.Errorf("%s:%d %v", fn, line, err)
	http.Error(response, sanitizer.Replace(err.Error()), http.StatusConflict)
}

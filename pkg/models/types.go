package models

import corev1 "k8s.io/api/core/v1"

// PageableResponse 分页
type PageableResponse struct {
	Items      []interface{} `json:"items" description:"paging data"`
	TotalCount int           `json:"total_count" description:"total count"`
}

// PodInfo 详情
type PodInfo struct {
	Namespace string `json:"namespace" description:"namespace"`
	Pod       string `json:"pod" description:"pod name"`
	Container string `json:"container" description:"container name"`
}

type ResourceQuota struct {
	Namespace string                     `json:"namespace" description:"namespace"`
	Data      corev1.ResourceQuotaStatus `json:"data" description:"resource quota status"`
}

package v1alpha2

import (
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"zmc.io/oasis/pkg/server/params"
	"zmc.io/oasis/pkg/utils/sliceutil"
)

const (
	Release          = "release"
	Name             = "name"
	Label            = "label"
	TargetKind       = "targetKind"
	TargetName       = "targetName"
	Role             = "role"
	CreateTime       = "createTime"
	UpdateTime       = "updateTime"
	StartTime        = "startTime"
	LastScheduleTime = "lastScheduleTime"
	Annotation       = "Annotation"
	Keyword          = "keyword"
	Status           = "status"

	StatusRunning            = "running"
	StatusPaused             = "paused"
	StatusPending            = "pending"
	StatusUpdating           = "updating"
	StatusStopped            = "stopped"
	StatusFailed             = "failed"
	StatusBound              = "bound"
	StatusLost               = "lost"
	StatusComplete           = "completed"
	StatusWarning            = "warning"
	StatusUnschedulable      = "unschedulable"
	Deployments              = "deployments"
	DaemonSets               = "daemonsets"
	Roles                    = "roles"
	CronJobs                 = "cronjobs"
	ConfigMaps               = "configmaps"
	Ingresses                = "ingresses"
	Jobs                     = "jobs"
	PersistentVolumeClaims   = "persistentvolumeclaims"
	Pods                     = "pods"
	Secrets                  = "secrets"
	Services                 = "services"
	StatefulSets             = "statefulsets"
	HorizontalPodAutoscalers = "horizontalpodautoscalers"
	// Applications             = "applications"
	Nodes          = "nodes"
	Namespaces     = "namespaces"
	StorageClasses = "storageclasses"
	ClusterRoles   = "clusterroles"
)

type Interface interface {
	Get(namespace, name string) (interface{}, error)
	Search(namespace string, conditions *params.Conditions, orderBy string, reverse bool) ([]interface{}, error)
}

func ObjectMetaExactlyMath(key, value string, item metav1.ObjectMeta) bool {
	switch key {
	case Name:
		names := strings.Split(value, ",")
		if !sliceutil.HasString(names, item.Name) {
			return false
		}
	case Keyword:
		if !strings.Contains(item.Name, value) && !FuzzyMatch(item.Labels, "", value) && !FuzzyMatch(item.Annotations, "", value) {
			return false
		}
	default:
		// label not exist or value not equal
		if val, ok := item.Labels[key]; !ok || val != value {
			return false
		}
	}
	return true
}

func ObjectMetaFuzzyMath(key, value string, item metav1.ObjectMeta) bool {
	switch key {
	case Name:
		if !strings.Contains(item.Name, value) {
			return false
		}
	case Label:
		if !FuzzyMatch(item.Labels, "", value) {
			return false
		}
	case Annotation:
		if !FuzzyMatch(item.Annotations, "", value) {
			return false
		}
		return false
	default:
		if !FuzzyMatch(item.Labels, key, value) {
			return false
		}
	}
	return true
}

func FuzzyMatch(m map[string]string, key, value string) bool {

	val, exist := m[key]

	if value == "" && (!exist || val == "") {
		return true
	} else if value != "" && strings.Contains(val, value) {
		return true
	}

	return false
}

func ObjectMetaCompare(left, right metav1.ObjectMeta, compareField string) bool {
	switch compareField {
	case CreateTime:
		if left.CreationTimestamp.Time.Equal(right.CreationTimestamp.Time) {
			if left.Namespace == right.Namespace {
				return strings.Compare(left.Name, right.Name) < 0
			}
			return strings.Compare(left.Namespace, right.Namespace) < 0
		}
		return left.CreationTimestamp.Time.Before(right.CreationTimestamp.Time)
	case Name:
		fallthrough
	default:
		return strings.Compare(left.Name, right.Name) <= 0
	}
}

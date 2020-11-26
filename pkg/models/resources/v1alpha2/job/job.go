package job

import (
	"k8s.io/client-go/informers"
	"zmc.io/oasis/pkg/models/resources/v1alpha2"

	"sort"
	"time"

	"zmc.io/oasis/pkg/server/params"
	"zmc.io/oasis/pkg/utils/k8sutil"

	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/labels"
)

const (
	includeCronJob = "includeCronJob"
	cronJobKind    = "CronJob"
	s2iRunKind     = "S2iRun"
	includeS2iRun  = "includeS2iRun"
)

type jobSearcher struct {
	informers informers.SharedInformerFactory
}

func NewJobSearcher(informers informers.SharedInformerFactory) v1alpha2.Interface {
	return &jobSearcher{informers: informers}
}

func (s *jobSearcher) Get(namespace, name string) (interface{}, error) {
	return s.informers.Batch().V1().Jobs().Lister().Jobs(namespace).Get(name)
}

func jobStatus(item *batchv1.Job) string {
	status := v1alpha2.StatusFailed
	if item.Status.Active > 0 {
		status = v1alpha2.StatusRunning
	} else if item.Status.Failed > 0 {
		status = v1alpha2.StatusFailed
	} else if item.Status.Succeeded > 0 {
		status = v1alpha2.StatusComplete
	}
	return status
}

func (*jobSearcher) match(match map[string]string, item *batchv1.Job) bool {
	for k, v := range match {
		switch k {
		case v1alpha2.Status:
			if jobStatus(item) != v {
				return false
			}
		case includeCronJob:
			if v == "false" && k8sutil.IsControlledBy(item.OwnerReferences, cronJobKind, "") {
				return false
			}
		case includeS2iRun:
			if v == "false" && k8sutil.IsControlledBy(item.OwnerReferences, s2iRunKind, "") {
				return false
			}
		default:
			if !v1alpha2.ObjectMetaExactlyMath(k, v, item.ObjectMeta) {
				return false
			}
		}
	}
	return true
}

func (*jobSearcher) fuzzy(fuzzy map[string]string, item *batchv1.Job) bool {
	for k, v := range fuzzy {
		if !v1alpha2.ObjectMetaFuzzyMath(k, v, item.ObjectMeta) {
			return false
		}
	}
	return true
}

func jobUpdateTime(item *batchv1.Job) time.Time {
	updateTime := item.CreationTimestamp.Time
	for _, condition := range item.Status.Conditions {
		if updateTime.Before(condition.LastProbeTime.Time) {
			updateTime = condition.LastProbeTime.Time
		}
		if updateTime.Before(condition.LastTransitionTime.Time) {
			updateTime = condition.LastTransitionTime.Time
		}
	}
	return updateTime
}

func (*jobSearcher) compare(left, right *batchv1.Job, orderBy string) bool {
	switch orderBy {
	case v1alpha2.UpdateTime:
		return jobUpdateTime(left).Before(jobUpdateTime(right))
	default:
		return v1alpha2.ObjectMetaCompare(left.ObjectMeta, right.ObjectMeta, orderBy)
	}
}

func (s *jobSearcher) Search(namespace string, conditions *params.Conditions, orderBy string, reverse bool) ([]interface{}, error) {
	jobs, err := s.informers.Batch().V1().Jobs().Lister().Jobs(namespace).List(labels.Everything())

	if err != nil {
		return nil, err
	}

	result := make([]*batchv1.Job, 0)

	if len(conditions.Match) == 0 && len(conditions.Fuzzy) == 0 {
		result = jobs
	} else {
		for _, item := range jobs {
			if s.match(conditions.Match, item) && s.fuzzy(conditions.Fuzzy, item) {
				result = append(result, item)
			}
		}
	}
	sort.Slice(result, func(i, j int) bool {
		if reverse {
			i, j = j, i
		}
		return s.compare(result[i], result[j], orderBy)
	})

	r := make([]interface{}, 0)
	for _, i := range result {
		r = append(r, i)
	}
	return r, nil
}

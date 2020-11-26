package cronjob

import (
	"sort"

	"k8s.io/api/batch/v1beta1"
	"k8s.io/client-go/informers"
	"zmc.io/oasis/pkg/models/resources/v1alpha2"
	"zmc.io/oasis/pkg/server/params"

	"k8s.io/apimachinery/pkg/labels"
)

type cronJobSearcher struct {
	informer informers.SharedInformerFactory
}

func NewCronJobSearcher(informer informers.SharedInformerFactory) v1alpha2.Interface {
	return &cronJobSearcher{informer: informer}
}

func (c *cronJobSearcher) Get(namespace, name string) (interface{}, error) {
	return c.informer.Batch().V1beta1().CronJobs().Lister().CronJobs(namespace).Get(name)
}

func cronJobStatus(item *v1beta1.CronJob) string {
	if item.Spec.Suspend != nil && *item.Spec.Suspend {
		return v1alpha2.StatusPaused
	}
	return v1alpha2.StatusRunning
}

func (*cronJobSearcher) match(match map[string]string, item *v1beta1.CronJob) bool {
	for k, v := range match {
		switch k {
		case v1alpha2.Status:
			if cronJobStatus(item) != v {
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

func (*cronJobSearcher) fuzzy(fuzzy map[string]string, item *v1beta1.CronJob) bool {

	for k, v := range fuzzy {
		if !v1alpha2.ObjectMetaFuzzyMath(k, v, item.ObjectMeta) {
			return false
		}
	}

	return true
}

func (*cronJobSearcher) compare(left, right *v1beta1.CronJob, orderBy string) bool {
	switch orderBy {
	case v1alpha2.LastScheduleTime:
		if left.Status.LastScheduleTime == nil {
			return true
		}
		if right.Status.LastScheduleTime == nil {
			return false
		}
		return left.Status.LastScheduleTime.Before(right.Status.LastScheduleTime)
	default:
		return v1alpha2.ObjectMetaCompare(left.ObjectMeta, right.ObjectMeta, orderBy)
	}
}

func (c *cronJobSearcher) Search(namespace string, conditions *params.Conditions, orderBy string, reverse bool) ([]interface{}, error) {
	cronJobs, err := c.informer.Batch().V1beta1().CronJobs().Lister().CronJobs(namespace).List(labels.Everything())

	if err != nil {
		return nil, err
	}

	result := make([]*v1beta1.CronJob, 0)

	if len(conditions.Match) == 0 && len(conditions.Fuzzy) == 0 {
		result = cronJobs
	} else {
		for _, item := range cronJobs {
			if c.match(conditions.Match, item) && c.fuzzy(conditions.Fuzzy, item) {
				result = append(result, item)
			}
		}
	}
	sort.Slice(result, func(i, j int) bool {
		if reverse {
			return !c.compare(result[i], result[j], orderBy)
		} else {
			return c.compare(result[i], result[j], orderBy)
		}
	})

	r := make([]interface{}, 0)
	for i := range result {
		r = append(r, result[i])
	}
	return r, nil
}

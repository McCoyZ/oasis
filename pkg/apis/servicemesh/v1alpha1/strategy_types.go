package v1alpha1

import (
	"istio.io/api/networking/v1beta1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceKindStrategy     = "Strategy"
	ResourceSingularStrategy = "strategy"
	ResourcePluralStrategy   = "strategies"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Strategy is the Schema for the strategies API
type Strategy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StrategySpec   `json:"spec,omitempty"`
	Status StrategyStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// StrategyList contains a list of Strategy
type StrategyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Strategy `json:"items"`
}

type StrategyType string

const (
	// Canary strategy type
	CanaryType StrategyType = "Canary"

	// BlueGreen strategy type
	BlueGreenType StrategyType = "BlueGreen"

	// Mirror strategy type
	Mirror StrategyType = "Mirror"
)

type StrategyPolicy string

const (
	// apply strategy only until workload is ready
	PolicyWaitForWorkloadReady StrategyPolicy = "WaitForWorkloadReady"

	// apply strategy immediately no matter workload status is
	PolicyImmediately StrategyPolicy = "Immediately"

	// pause strategy
	PolicyPause StrategyPolicy = "Paused"
)

// StrategySpec defines the desired state of Strategy
type StrategySpec struct {
	// Strategy type
	Type StrategyType `json:"type,omitempty"`

	// Principal version, the one as reference version
	// label version value
	// +optional
	PrincipalVersion string `json:"principal,omitempty"`

	// Governor version, the version takes control of all incoming traffic
	// label version value
	// +optional
	GovernorVersion string `json:"governor,omitempty"`

	// Label selector for virtual services.
	// +optional
	Selector *metav1.LabelSelector `json:"selector,omitempty"`

	// Template describes the virtual service that will be created.
	Template VirtualServiceTemplateSpec `json:"template,omitempty"`

	// strategy policy, how the strategy will be applied
	// by the strategy controller
	StrategyPolicy StrategyPolicy `json:"strategyPolicy,omitempty"`
}

// VirtualServiceTemplateSpec
type VirtualServiceTemplateSpec struct {

	// Metadata of the virtual services created from this template
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec indicates the behavior of a virtual service.
	// +optional
	Spec v1beta1.VirtualService `json:"spec,omitempty"`
}

// StrategyStatus defines the observed state of Strategy
type StrategyStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// The latest available observations of an object's current state.
	// +optional
	Conditions []StrategyCondition `json:"conditions,omitempty"`

	// Represents time when the strategy was acknowledged by the controller.
	// It is represented in RFC3339 form and is in UTC.
	// +optional
	StartTime *metav1.Time `json:"startTime,omitempty"`

	// Represents time when the strategy was completed.
	// It is represented in RFC3339 form and is in UTC.
	// +optional
	CompletionTime *metav1.Time `json:"completionTime,omitempty"`
}

type StrategyConditionType string

// These are valid conditions of a strategy.
const (
	// StrategyComplete means the strategy has been delivered to istio.
	StrategyComplete StrategyConditionType = "Complete"

	// StrategyFailed means the strategy has failed its delivery to istio.
	StrategyFailed StrategyConditionType = "Failed"
)

// StrategyCondition describes current state of a strategy.
type StrategyCondition struct {
	// Type of strategy condition, Complete or Failed.
	Type StrategyConditionType `json:"type,omitempty"`

	// Status of the condition, one of True, False, Unknown
	Status v1.ConditionStatus `json:"status,omitempty"`

	// Last time the condition was checked.
	// +optional
	LastProbeTime metav1.Time `json:"lastProbeTime,omitempty"`

	// Last time the condition transit from one status to another
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`

	// reason for the condition's last transition
	Reason string `json:"reason,omitempty"`

	// Human readable message indicating details about last transition.
	// +optinal
	Message string `json:"message,omitempty"`
}

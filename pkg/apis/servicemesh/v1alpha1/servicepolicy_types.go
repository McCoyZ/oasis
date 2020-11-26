package v1alpha1

import (
	"istio.io/api/networking/v1beta1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceKindServicePolicy     = "ServicePolicy"
	ResourceSingularServicePolicy = "servicepolicy"
	ResourcePluralServicePolicy   = "servicepolicies"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ServicePolicy is the Schema for the servicepolicies API
type ServicePolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ServicePolicySpec   `json:"spec,omitempty"`
	Status ServicePolicyStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ServicePolicyList contains a list of ServicePolicy
type ServicePolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ServicePolicy `json:"items"`
}

// ServicePolicySpec defines the desired state of ServicePolicy
type ServicePolicySpec struct {

	// Label selector for destination rules.
	// +optional
	Selector *metav1.LabelSelector `json:"selector,omitempty"`

	// Template used to create a destination rule
	// +optional
	Template DestinationRuleSpecTemplate `json:"template,omitempty"`
}

type DestinationRuleSpecTemplate struct {

	// Metadata of the virtual services created from this template
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec indicates the behavior of a destination rule.
	// +optional
	Spec v1beta1.DestinationRule `json:"spec,omitempty"`
}

type ServicePolicyConditionType string

// These are valid conditions of a strategy.
const (
	// StrategyComplete means the strategy has been delivered to istio.
	ServicePolicyComplete ServicePolicyConditionType = "Complete"

	// StrategyFailed means the strategy has failed its delivery to istio.
	ServicePolicyFailed ServicePolicyConditionType = "Failed"
)

// StrategyCondition describes current state of a strategy.
type ServicePolicyCondition struct {
	// Type of strategy condition, Complete or Failed.
	Type ServicePolicyConditionType `json:"type,omitempty"`

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

// ServicePolicyStatus defines the observed state of ServicePolicy
type ServicePolicyStatus struct {
	// The latest available observations of an object's current state.
	// +optional
	Conditions []ServicePolicyCondition `json:"conditions,omitempty"`

	// Represents time when the strategy was acknowledged by the controller.
	// It is represented in RFC3339 form and is in UTC.
	// +optional
	StartTime *metav1.Time `json:"startTime,omitempty"`

	// Represents time when the strategy was completed.
	// It is represented in RFC3339 form and is in UTC.
	// +optional
	CompletionTime *metav1.Time `json:"completionTime,omitempty"`
}

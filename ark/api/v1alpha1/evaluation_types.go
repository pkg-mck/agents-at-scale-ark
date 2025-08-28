/* Copyright 2025. McKinsey & Company */

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EvaluationEvaluatorRef references an evaluator resource for evaluation with parameters
type EvaluationEvaluatorRef struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`
	// +kubebuilder:validation:Optional
	Namespace string `json:"namespace,omitempty"`
	// +kubebuilder:validation:Optional
	Parameters []Parameter `json:"parameters,omitempty"`
}

// QueryRef references a query for post-hoc evaluation
type QueryRef struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`
	// +kubebuilder:validation:Optional
	Namespace string `json:"namespace,omitempty"`
	// +kubebuilder:validation:Optional
	// Target name to match against query responses (e.g., "weather-agent", "summary-team")
	ResponseTarget string `json:"responseTarget,omitempty"`
}

// EvaluationConfig holds type-specific configuration parameters
type EvaluationConfig struct {
	// +kubebuilder:validation:Optional
	*DirectEvaluationConfig `json:",inline"`
	// +kubebuilder:validation:Optional
	*QueryBasedEvaluationConfig `json:",inline"`
	// +kubebuilder:validation:Optional
	*BatchEvaluationConfig `json:",inline"`
	// +kubebuilder:validation:Optional
	*BaselineEvaluationConfig `json:",inline"`
	// +kubebuilder:validation:Optional
	*EventEvaluationConfig `json:",inline"`
}

// DirectEvaluationConfig contains Direct Evaluation specific parameters
type DirectEvaluationConfig struct {
	// +kubebuilder:validation:Required
	Input string `json:"input"`
	// +kubebuilder:validation:Required
	Output string `json:"output"`
}

// QueryBasedEvaluationConfig contains Query related Evaluation specific parameters
type QueryBasedEvaluationConfig struct {
	// +kubebuilder:validation:Optional
	QueryRef *QueryRef `json:"queryRef,omitempty"`
}

// EvaluationRef references an evaluation to aggregate in batch type
type EvaluationRef struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`
	// +kubebuilder:validation:Optional
	Namespace string `json:"namespace,omitempty"`
}

// BatchEvaluationConfig contains Batch Evaluation specific parameters
type BatchEvaluationConfig struct {
	// +kubebuilder:validation:Optional
	// List of evaluations to aggregate
	Evaluations []EvaluationRef `json:"evaluations,omitempty"`
}

// BaselineEvaluationConfig contains Baseline Evaluation specific parameters
type BaselineEvaluationConfig struct{}

// EventEvaluationConfig, expression based evaluations, especially for tools
type EventEvaluationConfig struct {
	// +kubebuilder:validation:Optional
	Rules []ExpressionRule `json:"rules,omitempty"`
}

// EvaluationSpec defines the desired state of Evaluation
type EvaluationSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=direct;baseline;query;batch;event
	// +kubebuilder:default=direct
	Type string `json:"type,omitempty"`
	// +kubebuilder:validation:Required
	Config EvaluationConfig `json:"config"`
	// +kubebuilder:validation:Optional
	Evaluator EvaluationEvaluatorRef `json:"evaluator,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default="720h"
	TTL *metav1.Duration `json:"ttl,omitempty"`
	// +kubebuilder:default="5m"
	// Timeout for query execution (e.g., "30s", "5m", "1h")
	Timeout *metav1.Duration `json:"timeout,omitempty"`
}

// EvaluationStatus defines the observed state of Evaluation
type EvaluationStatus struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=pending;running;error;done;canceled
	Phase string `json:"phase,omitempty"`
	// +kubebuilder:validation:Optional
	Message string `json:"message,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Pattern=^(0(\.[0-9]+)?|1(\.0+)?)$
	Score string `json:"score,omitempty"`
	// +kubebuilder:validation:Optional
	Passed bool `json:"passed"`
	// +kubebuilder:validation:Optional
	TokenUsage *TokenUsage `json:"tokenUsage,omitempty"`
	// +kubebuilder:validation:Optional
	Duration *metav1.Duration `json:"duration,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Type",type=string,JSONPath=`.spec.type`
// +kubebuilder:printcolumn:name="Phase",type=string,JSONPath=`.status.phase`
// +kubebuilder:printcolumn:name="Score",type=string,JSONPath=`.status.score`
// +kubebuilder:printcolumn:name="Passed",type=boolean,JSONPath=`.status.passed`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

type Evaluation struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EvaluationSpec   `json:"spec,omitempty"`
	Status EvaluationStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// EvaluationList contains a list of Evaluation
type EvaluationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Evaluation `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Evaluation{}, &EvaluationList{})
}

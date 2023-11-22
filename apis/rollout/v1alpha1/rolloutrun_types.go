// Copyright 2023 The KusionStack Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:storageversion
// +kubebuilder:subresource:status

type RolloutRun struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RolloutRunSpec   `json:"spec,omitempty"`
	Status RolloutRunStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

// RolloutList contains a list of Rollout
type RolloutRunList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RolloutRun `json:"items"`
}

type RolloutRunSpec struct {
	// TargetType defines the GroupVersionKind of target resource
	TargetType ObjectTypeRef `json:"targetType,omitempty"`

	// Webhooks defines rollout webhook configuration
	Webhooks []RolloutWebhook `json:"webhooks,omitempty"`

	// Batch Strategy
	Batch RolloutRunBatchStrategy `json:"batch,omitempty"`
}

type ObjectTypeRef struct {
	// APIVersion is the group/version for the resource being referenced.
	// If APIVersion is not specified, the specified Kind must be in the core API group.
	// For any other third-party types, APIVersion is required.
	// +optional
	APIVersion string `json:"apiVersion"`
	// Kind is the type of resource being referenced
	Kind string `json:"kind"`
}

type RolloutRunBatchStrategy struct {
	// Batches define the order of phases to execute release in canary release
	Batches []RolloutRunStep `json:"batches,omitempty"`
	// Toleration is the toleration policy of the canary strategy
	// +optional
	Toleration *TolerationStrategy `json:"toleration,omitempty"`
}

type RolloutRunStep struct {
	// traffic strategy
	TrafficStrategy `json:",inline"`

	// desired target replicas
	Targets []RolloutRunStepTarget `json:"targets"`

	// If true, rollout will be paused after this canary step complete.
	Pause *bool `json:"pause,omitempty"`

	// Properties contains additional information for step
	Properties map[string]string `json:"properties,omitempty"`
}

type RolloutRunStepTarget struct {
	CrossClusterObjectNameReference `json:",inline"`

	// Replicas is the replicas of the rollout task, which represents the number of pods to be upgraded
	Replicas intstr.IntOrString `json:"replicas"`
}

type RolloutRunStatus struct {
	// ObservedGeneration is the most recent generation observed for this Rollout. It corresponds to the
	// Rollout's generation, which is updated on mutation by the API Server.
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
	// Conditions is the list of conditions
	Conditions []Condition `json:"conditions,omitempty"`
	// Phase indecates the current phase of rollout
	Phase RolloutPhase `json:"phase,omitempty"`
	// BatchStatus describes the state of the active batch release
	BatchStatus *BatchStatus `json:"batchStatus,omitempty"`
	// WorkloadStatuses describes the referenced workloads status
	WorkloadStatuses []RolloutWorkloadStatus `json:"workloadStatuses,omitempty"`
	// The last time this status was updated.
	// +optional
	LastUpdateTime *metav1.Time `json:"lastUpdateTime,omitempty"`
}

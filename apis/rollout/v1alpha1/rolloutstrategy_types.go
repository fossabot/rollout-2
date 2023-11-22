/**
 * Copyright 2023 The KusionStack Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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

// RolloutStrategy is the Schema for the rolloutstrategies API
type RolloutStrategy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Batch is the batch strategy for upgrade and operation
	// +optional
	Batch *BatchStrategy `json:"batch,omitempty"`

	// Webhooks defines
	Webhooks []RolloutWebhook `json:"webhooks,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

// RolloutStrategyList contains a list of RolloutStrategy
type RolloutStrategyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RolloutStrategy `json:"items"`
}

// BatchStrategy defines the batch strategy
type BatchStrategy struct {
	// FastBatch indicates a fast way to create the batches
	FastBatch *FastBatch `json:"fastBatch,omitempty"`
	// Batches define the order of phases to execute release in canary release
	Batches []RolloutStep `json:"batches,omitempty"`
	// Toleration is the toleration policy of the canary strategy
	// +optional
	Toleration *TolerationStrategy `json:"toleration,omitempty"`
}

// FastBatch defines a fast way to create the batches
type FastBatch struct {
	// Beta defines the canary policy of batch release step, it will be the first step
	Beta *FastBetaBatch `json:"beta,omitempty"`
	// Count indicates the number of batch steps, workloads across different clusters in the same batch will be upgraded at the same time
	Count int32 `json:"count"`
	// PausedBatches is the breakpoints of batch release steps, which indicates pause points before determined batch.
	PausedBatches []int32 `json:"pausedBatches,omitempty"`
	// PauseMode is the pause mode of the strategy
	PauseMode PauseModeType `json:"pauseMode,omitempty"`
}

// PauseModeType is the type of the pause mode
type PauseModeType string

const (
	// PauseModeTypeNever is the pause mode of never
	PauseModeTypeNever PauseModeType = ""

	// PauseModeTypeFirstBatch is the pause mode of first batch
	PauseModeTypeFirstBatch PauseModeType = "FirstBatch"

	// PauseModeTypeEachBatch is the pause mode of each batch
	PauseModeTypeEachBatch PauseModeType = "EachBatch"
)

// FastBetaBatch is the first batch of fast batch
type FastBetaBatch struct {
	// traffic strategy
	TrafficStrategy `json:",inline"`

	// Replicas indicates the replicas of the beta(first) step, which will be chosen evenly from workloads
	Replicas *intstr.IntOrString `json:"replicas"`
}

// TolerationStrategy defines the toleration strategy
type TolerationStrategy struct {
	// WorkloadFailureThreshold indicates how many failed pods can be tolerated in all upgraded pods of one workload.
	// The default value is 0, which means no failed pods can be tolerated.
	// This is a workload level threshold.
	// +optional
	WorkloadFailureThreshold *intstr.IntOrString `json:"workloadTotalFailureThreshold,omitempty"`

	// FailureThreshold indicates how many failed pods can be tolerated before marking the rollout task as success
	// If not set, the default value is 0, which means no failed pods can be tolerated
	// This is a task level threshold.
	// +optional
	TaskFailureThreshold *intstr.IntOrString `json:"taskFailureThreshold,omitempty"`

	// Number of seconds after the toleration check has started before the task are initiated.
	InitialDelaySeconds int32 `json:"initialDelaySeconds,omitempty"`
}

// Custom release step
type RolloutStep struct {
	// traffic strategy
	TrafficStrategy `json:",inline"`

	// Replicas is the replicas of the rollout task, which represents the number of pods to be upgraded
	Replicas intstr.IntOrString `json:"replicas"`

	// Match defines condition used for matching resource cross clusterset
	Match *ResourceMatch `json:"matchTargets,omitempty"`

	// If true, rollout will be paused after this canary step complete.
	Pause *bool `json:"pause,omitempty"`

	// Properties contains additional information for step
	Properties map[string]string `json:"properties,omitempty"`
}

type TrafficStrategy struct {
	// Weight indicate how many percentage of traffic the canary pods should receive
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=100
	Weight       *int32               `json:"weight,omitempty"`
	HTTPStrategy *HTTPTrafficStrategy `json:"http,omitempty"`
}

type HTTPTrafficStrategy struct {
	// Matches define conditions used for matching the incoming HTTP requests to canary service.
	Matches []HTTPRouteMatch `json:"matches,omitempty"`
	// RequestHeaderModifier defines a schema for a filter that modifies request
	// headers.
	//
	// Support: Core
	//
	// +optional
	RequestHeaderModifier *HTTPHeaderFilter `json:"requestHeaderModifier,omitempty"`
}

type HTTPRouteMatch struct {
	// Headers specifies HTTP request header matchers. Multiple match values are
	// ANDed together, meaning, a request must match all the specified headers
	// to select the route.
	// +kubebuilder:validation:MaxItems=16
	Headers []HTTPHeaderMatch `json:"headers,omitempty"`
}

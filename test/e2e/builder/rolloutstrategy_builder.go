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

package builder

import (
	"fmt"
	"net/http/httptest"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	rolloutv1alpha1 "kusionstack.io/rollout/apis/rollout/v1alpha1"
)

const defaultStrategyName = DefaultName

// RolloutStrategyBuilder is a builder for RolloutStrategy
type RolloutStrategyBuilder struct {
	builder
}

// NewRolloutStrategy returns a RolloutStrategy builder
func NewRolloutStrategy() *RolloutStrategyBuilder {
	return &RolloutStrategyBuilder{}
}

func (b *RolloutStrategyBuilder) Namespace(namespace string) *RolloutStrategyBuilder {
	b.namespace = namespace
	return b
}

// Build returns a RolloutStrategy
func (b *RolloutStrategyBuilder) Build(ts *httptest.Server, labels map[string]string) *rolloutv1alpha1.RolloutStrategy {
	b.complete()

	return &rolloutv1alpha1.RolloutStrategy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      b.name,
			Namespace: b.namespace,
		},
		Batch: &rolloutv1alpha1.BatchStrategy{
			Toleration: &rolloutv1alpha1.TolerationStrategy{
				WorkloadFailureThreshold: &intstr.IntOrString{Type: intstr.String, StrVal: "10%"},
			},
			Batches: []rolloutv1alpha1.RolloutStep{
				{
					Breakpoint: true,
					Replicas:   intstr.FromInt(1),
					Match: &rolloutv1alpha1.ResourceMatch{
						Selector: &metav1.LabelSelector{MatchLabels: labels},
					},
				},
				{
					Breakpoint: true,
					Replicas:   intstr.FromInt(2),
					Match: &rolloutv1alpha1.ResourceMatch{
						Selector: &metav1.LabelSelector{MatchLabels: labels},
					},
				},
				{
					Breakpoint: true,
					Replicas:   intstr.FromString("100%"),
					Match: &rolloutv1alpha1.ResourceMatch{
						Selector: &metav1.LabelSelector{MatchLabels: labels},
					},
				},
			},
		},
		Webhooks: []rolloutv1alpha1.RolloutWebhook{
			{
				Name:             "wh-01",
				FailureThreshold: 2,
				FailurePolicy:    rolloutv1alpha1.Ignore,
				HookTypes:        []rolloutv1alpha1.HookType{rolloutv1alpha1.HookTypePreBatchStep, rolloutv1alpha1.HookTypePostBatchStep},
				ClientConfig:     rolloutv1alpha1.WebhookClientConfig{TimeoutSeconds: 5, PeriodSeconds: 3, URL: fmt.Sprintf("%s/webhook?sleepSeconds=2", ts.URL)},
				Properties:       map[string]string{"responseBody": "{\"code\":\"OK\",\"reason\":\"Success\",\"message\":\"Success\"}"},
			},
			{
				Name:             "wh-02",
				FailureThreshold: 2,
				FailurePolicy:    rolloutv1alpha1.Ignore,
				HookTypes:        []rolloutv1alpha1.HookType{rolloutv1alpha1.HookTypePreBatchStep, rolloutv1alpha1.HookTypePostBatchStep},
				ClientConfig:     rolloutv1alpha1.WebhookClientConfig{TimeoutSeconds: 5, PeriodSeconds: 3, URL: fmt.Sprintf("%s/webhook?sleepSeconds=2", ts.URL)},
				Properties:       map[string]string{"responseBody": "{\"code\":\"OK\",\"reason\":\"Success\",\"message\":\"Success\"}"},
			},
		},
	}
}

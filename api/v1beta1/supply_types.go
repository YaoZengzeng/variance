/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SupplySpec defines the desired state of Supply
type SupplySpec struct {
	// +kubebuilder:validation:Minimum=0

	// Request of the Supply.
	Request *int64 `json:"request"`
}

// SupplyStatus defines the observed state of Supply
type SupplyStatus struct {
	// Specify the phase supply.
	Phase Phase `json:"phase"`

	// Allocations to satisify supply.
	Allocations []Allocation `json:"allocations"`
}

// Phase describes the stage of supply
type Phase string

const (
	// Succeeded to get funds from the fund pool.
	Succeeded Phase = "Succeeded"
	// Failed to get funds from the fund pool.
	Failed Phase = "Failed"
)

type Allocation struct {
	Pool       string `json:"pool"`
	Shortfalls *int64 `json:"shortfalls"`
}

// +kubebuilder:object:root=true

// Supply is the Schema for the supplies API
type Supply struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SupplySpec   `json:"spec,omitempty"`
	Status SupplyStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SupplyList contains a list of Supply
type SupplyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Supply `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Supply{}, &SupplyList{})
}

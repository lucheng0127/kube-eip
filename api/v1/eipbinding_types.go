/*
Copyright 2023.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// +kubebuilder:validation:Enum=Ready;Processing;Error
type EipBindingPhase string

const (
	PhaseReady      EipBindingPhase = "Ready"
	PhaseProcessing EipBindingPhase = "Processing"
	PhaseError      EipBindingPhase = "Error"
)

// EipBindingSpec defines the desired state of EipBinding
type EipBindingSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Eip address binded to vmi
	EipAddr string `json:"eipAddr"`

	// Name of kubevirt vmi
	VmiName string `json:"vmiName"`

	// Hyper vmi pod placed on
	CurrentHyper string `json:"currentHyper,omitempty"`

	// The vmi pod ip address
	CurrentIPAddr string `json:"currentIPAddr,omitempty"`

	// Hyper the last vmi pod placed on
	LastHyper string `json:"lastHyper,omitempty"`

	// The last vmi pod ip address
	LastIPAddr string `json:"lastIPAddr,omitempty"`

	// Eip binding pahse
	Phase EipBindingPhase `json:"phase,omitempty"`
}

// EipBindingStatus defines the observed state of EipBinding
type EipBindingStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// EipBinding is the Schema for the eipbindings API
type EipBinding struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EipBindingSpec   `json:"spec,omitempty"`
	Status EipBindingStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// EipBindingList contains a list of EipBinding
type EipBindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EipBinding `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EipBinding{}, &EipBindingList{})
}

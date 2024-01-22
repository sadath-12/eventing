/*
Copyright 2020 The Knative Authors

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
	"k8s.io/apimachinery/pkg/runtime"
	eventingduckv1 "knative.dev/eventing/pkg/apis/duck/v1"
	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"
	"knative.dev/pkg/kmeta"
)

// +genclient
// +genreconciler
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ReferenceGrant struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec defines the desired state of ReferenceGrant.
	Spec ReferenceGrantSpec `json:"spec,omitempty"`

	// Note that `Status` sub-resource has been excluded at the
	// moment as it was difficult to work out the design.
	// `Status` sub-resource may be added in future.
}

// +kubebuilder:object:root=true
// ReferenceGrantList contains a list of ReferenceGrant.
type ReferenceGrantList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ReferenceGrant `json:"items"`
}

// ReferenceGrantSpec identifies a cross namespace relationship that is trusted
// for Gateway API.
type ReferenceGrantSpec struct {
	// From describes the trusted namespaces and kinds that can reference the
	// resources described in "To". Each entry in this list MUST be considered
	// to be an additional place that references can be valid from, or to put
	// this another way, entries MUST be combined using OR.
	//
	// Support: Core
	//
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=16
	From []ReferenceGrantFrom `json:"from"`

	// To describes the resources that may be referenced by the resources
	// described in "From". Each entry in this list MUST be considered to be an
	// additional place that references can be valid to, or to put this another
	// way, entries MUST be combined using OR.
	//
	// Support: Core
	//
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=16
	To []ReferenceGrantTo `json:"to"`
}

// ReferenceGrantFrom describes trusted namespaces and kinds.
type ReferenceGrantFrom struct {
	// Group is the group of the referent.
	// When empty, the Kubernetes core API group is inferred.
	//
	// Support: Core
	Group Group `json:"group"`

	// Kind is the kind of the referent. Although implementations may support
	// additional resources, the following types are part of the "Core"
	// support level for this field.
	//
	// When used to permit a SecretObjectReference:
	//
	// * Gateway
	//
	// When used to permit a BackendObjectReference:
	//
	// * GRPCRoute
	// * HTTPRoute
	// * TCPRoute
	// * TLSRoute
	// * UDPRoute
	Kind Kind `json:"kind"`

	// Namespace is the namespace of the referent.
	//
	// Support: Core
	Namespace Namespace `json:"namespace"`
}

// ReferenceGrantTo describes what Kinds are allowed as targets of the
// references.
type ReferenceGrantTo struct {
	// Group is the group of the referent.
	// When empty, the Kubernetes core API group is inferred.
	//
	// Support: Core
	Group Group `json:"group"`

	// Kind is the kind of the referent. Although implementations may support
	// additional resources, the following types are part of the "Core"
	// support level for this field:
	//
	// * Secret when used to permit a SecretObjectReference
	// * Service when used to permit a BackendObjectReference
	Kind Kind `json:"kind"`

	// Name is the name of the referent. When unspecified, this policy
	// refers to all resources of the specified Group and Kind in the local
	// namespace.
	//
	// +optional
	Name *ObjectName `json:"name,omitempty"`
}

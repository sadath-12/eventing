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
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/utils/pointer"
	"knative.dev/pkg/apis"
	v1 "knative.dev/pkg/apis/duck/v1"

	eventingduck "knative.dev/eventing/pkg/apis/duck/v1"
)

var imcCondSet = apis.NewLivingConditionSet(
	InMemoryChannelConditionServiceReady,
	InMemoryChannelConditionEndpointsReady,
	InMemoryChannelConditionAddressable,
	InMemoryChannelConditionChannelServiceReady,
	InMemoryChannelConditionDeadLetterSinkResolved,
)

const (
	// InMemoryChannelConditionReady has status True when all subconditions below have been set to True.
	InMemoryChannelConditionReady = apis.ConditionReady

	// InMemoryChannelConditionDispatcherReady has status True when a Dispatcher deployment is ready
	// Keyed off appsv1.DeploymentAvailable, which means minimum available replicas required are up
	// and running for at least minReadySeconds.
	InMemoryChannelConditionDispatcherReady apis.ConditionType = "DispatcherReady"

	// InMemoryChannelConditionServiceReady has status True when a k8s Service is ready. This
	// basically just means it exists because there's no meaningful status in Service. See Endpoints
	// below.
	InMemoryChannelConditionServiceReady apis.ConditionType = "ServiceReady"

	// InMemoryChannelConditionEndpointsReady has status True when a k8s Service Endpoints are backed
	// by at least one endpoint.
	InMemoryChannelConditionEndpointsReady apis.ConditionType = "EndpointsReady"

	// InMemoryChannelConditionAddressable has status true when this InMemoryChannel meets
	// the Addressable contract and has a non-empty hostname.
	InMemoryChannelConditionAddressable apis.ConditionType = "Addressable"

	// InMemoryChannelConditionServiceReady has status True when a k8s Service representing the channel is ready.
	// Because this uses ExternalName, there are no endpoints to check.
	InMemoryChannelConditionChannelServiceReady apis.ConditionType = "ChannelServiceReady"

	// InMemoryChannelConditionDeadLetterSinkResolved has status True when there is a Dead Letter Sink ref or URI
	// defined in the Spec.Delivery, is a valid destination and its correctly resolved into a valid URI
	InMemoryChannelConditionDeadLetterSinkResolved apis.ConditionType = "DeadLetterSinkResolved"
)

// GetConditionSet retrieves the condition set for this resource. Implements the KRShaped interface.
func (*InMemoryChannel) GetConditionSet() apis.ConditionSet {
	return imcCondSet
}

// GetGroupVersionKind returns GroupVersionKind for InMemoryChannels
func (*InMemoryChannel) GetGroupVersionKind() schema.GroupVersionKind {
	return SchemeGroupVersion.WithKind("InMemoryChannel")
}

// GetUntypedSpec returns the spec of the InMemoryChannel.
func (i *InMemoryChannel) GetUntypedSpec() interface{} {
	return i.Spec
}

// GetCondition returns the condition currently associated with the given type, or nil.
func (imcs *InMemoryChannelStatus) GetCondition(t apis.ConditionType) *apis.Condition {
	return imcCondSet.Manage(imcs).GetCondition(t)
}

// IsReady returns true if the Status condition InMemoryChannelConditionReady
// is true and the latest spec has been observed.
func (imc *InMemoryChannel) IsReady() bool {
	imcs := imc.Status
	return imcs.ObservedGeneration == imc.Generation &&
		imc.GetConditionSet().Manage(&imcs).IsHappy()
}

// InitializeConditions sets relevant unset conditions to Unknown state.
func (imcs *InMemoryChannelStatus) InitializeConditions() {
	imcCondSet.Manage(imcs).InitializeConditions()
}

func (imcs *InMemoryChannelStatus) SetAddress(addr *v1.Addressable) { 
	imcs.Address = addr
	if addr != nil && addr.URL != nil {
		imcs.Address.Name = pointer.String(addr.URL.Scheme)
		imcCondSet.Manage(imcs).MarkTrue(InMemoryChannelConditionAddressable)
	} else {
		imcCondSet.Manage(imcs).MarkFalse(InMemoryChannelConditionAddressable, "emptyHostname", "hostname is the empty string")
	}
}

func (imcs *InMemoryChannelStatus) MarkDispatcherFailed(reason, messageFormat string, messageA ...interface{}) {
	imcCondSet.Manage(imcs).MarkFalse(InMemoryChannelConditionDispatcherReady, reason, messageFormat, messageA...)
}

func (imcs *InMemoryChannelStatus) MarkDispatcherUnknown(reason, messageFormat string, messageA ...interface{}) {
	imcCondSet.Manage(imcs).MarkUnknown(InMemoryChannelConditionDispatcherReady, reason, messageFormat, messageA...)
}

// TODO: Unify this with the ones from Eventing. Say: Broker, Trigger.
func (imcs *InMemoryChannelStatus) PropagateDispatcherStatus(ds *appsv1.DeploymentStatus) {
	for _, cond := range ds.Conditions {
		if cond.Type == appsv1.DeploymentAvailable {
			if cond.Status == corev1.ConditionTrue {
				imcCondSet.Manage(imcs).MarkTrue(InMemoryChannelConditionDispatcherReady)
			} else if cond.Status == corev1.ConditionFalse {
				imcs.MarkDispatcherFailed("DispatcherDeploymentFalse", "The status of Dispatcher Deployment is False: %s : %s", cond.Reason, cond.Message)
			} else if cond.Status == corev1.ConditionUnknown {
				imcs.MarkDispatcherUnknown("DispatcherDeploymentUnknown", "The status of Dispatcher Deployment is Unknown: %s : %s", cond.Reason, cond.Message)
			}
		}
	}
}

func (imcs *InMemoryChannelStatus) MarkServiceFailed(reason, messageFormat string, messageA ...interface{}) {
	imcCondSet.Manage(imcs).MarkFalse(InMemoryChannelConditionServiceReady, reason, messageFormat, messageA...)
}

func (imcs *InMemoryChannelStatus) MarkServiceUnknown(reason, messageFormat string, messageA ...interface{}) {
	imcCondSet.Manage(imcs).MarkUnknown(InMemoryChannelConditionServiceReady, reason, messageFormat, messageA...)
}

func (imcs *InMemoryChannelStatus) MarkServiceTrue() {
	imcCondSet.Manage(imcs).MarkTrue(InMemoryChannelConditionServiceReady)
}

func (imcs *InMemoryChannelStatus) MarkChannelServiceFailed(reason, messageFormat string, messageA ...interface{}) {
	imcCondSet.Manage(imcs).MarkFalse(InMemoryChannelConditionChannelServiceReady, reason, messageFormat, messageA...)
}

func (imcs *InMemoryChannelStatus) MarkChannelServiceUnknown(reason, messageFormat string, messageA ...interface{}) {
	imcCondSet.Manage(imcs).MarkUnknown(InMemoryChannelConditionChannelServiceReady, reason, messageFormat, messageA...)
}

func (imcs *InMemoryChannelStatus) MarkChannelServiceTrue() {
	imcCondSet.Manage(imcs).MarkTrue(InMemoryChannelConditionChannelServiceReady)
}

func (imcs *InMemoryChannelStatus) MarkEndpointsFailed(reason, messageFormat string, messageA ...interface{}) {
	imcCondSet.Manage(imcs).MarkFalse(InMemoryChannelConditionEndpointsReady, reason, messageFormat, messageA...)
}

func (imcs *InMemoryChannelStatus) MarkEndpointsUnknown(reason, messageFormat string, messageA ...interface{}) {
	imcCondSet.Manage(imcs).MarkUnknown(InMemoryChannelConditionEndpointsReady, reason, messageFormat, messageA...)
}

func (imcs *InMemoryChannelStatus) MarkEndpointsTrue() {
	imcCondSet.Manage(imcs).MarkTrue(InMemoryChannelConditionEndpointsReady)
}

func (imcs *InMemoryChannelStatus) MarkDeadLetterSinkResolvedSucceeded(ds eventingduck.DeliveryStatus) {
	imcs.DeliveryStatus = ds
	imcCondSet.Manage(imcs).MarkTrue(InMemoryChannelConditionDeadLetterSinkResolved)
}

func (imcs *InMemoryChannelStatus) MarkDeadLetterSinkNotConfigured() {
	imcs.DeliveryStatus = eventingduck.DeliveryStatus{}
	imcCondSet.Manage(imcs).MarkTrueWithReason(InMemoryChannelConditionDeadLetterSinkResolved, "DeadLetterSinkNotConfigured", "No dead letter sink is configured.")
}

func (imcs *InMemoryChannelStatus) MarkDeadLetterSinkResolvedFailed(reason, messageFormat string, messageA ...interface{}) {
	imcs.DeliveryStatus = eventingduck.DeliveryStatus{}
	imcCondSet.Manage(imcs).MarkFalse(InMemoryChannelConditionDeadLetterSinkResolved, reason, messageFormat, messageA...)
}

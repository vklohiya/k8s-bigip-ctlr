/*
Copyright 2025 F5 Networks.

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
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	cisv1 "github.com/F5Networks/k8s-bigip-ctlr/v2/api/v1"
)

// nolint:unused
// log is for logging in this package.
var transportserverlog = logf.Log.WithName("transportserver-resource")

// SetupTransportServerWebhookWithManager registers the webhook for TransportServer in the manager.
func SetupTransportServerWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).For(&cisv1.TransportServer{}).
		WithValidator(&TransportServerCustomValidator{}).
		WithDefaulter(&TransportServerCustomDefaulter{}).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-cis-f5-com-v1-transportserver,mutating=true,failurePolicy=fail,sideEffects=None,groups=cis.f5.com,resources=transportservers,verbs=create;update,versions=v1,name=mtransportserver-v1.kb.io,admissionReviewVersions=v1

// TransportServerCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind TransportServer when those are created or updated.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as it is used only for temporary operations and does not need to be deeply copied.
type TransportServerCustomDefaulter struct {
	// TODO(user): Add more fields as needed for defaulting
}

var _ webhook.CustomDefaulter = &TransportServerCustomDefaulter{}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the Kind TransportServer.
func (d *TransportServerCustomDefaulter) Default(_ context.Context, obj runtime.Object) error {
	transportserver, ok := obj.(*cisv1.TransportServer)

	if !ok {
		return fmt.Errorf("expected an TransportServer object but got %T", obj)
	}
	transportserverlog.Info("Defaulting for TransportServer", "name", transportserver.GetName())

	// TODO(user): fill in your defaulting logic.

	return nil
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// NOTE: The 'path' attribute must follow a specific pattern and should not be modified directly here.
// Modifying the path for an invalid path can cause API server errors; failing to locate the webhook.
// +kubebuilder:webhook:path=/validate-cis-f5-com-v1-transportserver,mutating=false,failurePolicy=fail,sideEffects=None,groups=cis.f5.com,resources=transportservers,verbs=create;update,versions=v1,name=vtransportserver-v1.kb.io,admissionReviewVersions=v1

// TransportServerCustomValidator struct is responsible for validating the TransportServer resource
// when it is created, updated, or deleted.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
type TransportServerCustomValidator struct {
	// TODO(user): Add more fields as needed for validation
}

var _ webhook.CustomValidator = &TransportServerCustomValidator{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type TransportServer.
func (v *TransportServerCustomValidator) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	transportserver, ok := obj.(*cisv1.TransportServer)
	if !ok {
		return nil, fmt.Errorf("expected a TransportServer object but got %T", obj)
	}
	transportserverlog.Info("Validation for TransportServer upon creation", "name", transportserver.GetName())

	// TODO(user): fill in your validation logic upon object creation.

	return nil, nil
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type TransportServer.
func (v *TransportServerCustomValidator) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	transportserver, ok := newObj.(*cisv1.TransportServer)
	if !ok {
		return nil, fmt.Errorf("expected a TransportServer object for the newObj but got %T", newObj)
	}
	transportserverlog.Info("Validation for TransportServer upon update", "name", transportserver.GetName())

	// TODO(user): fill in your validation logic upon object update.

	return nil, nil
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type TransportServer.
func (v *TransportServerCustomValidator) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	transportserver, ok := obj.(*cisv1.TransportServer)
	if !ok {
		return nil, fmt.Errorf("expected a TransportServer object but got %T", obj)
	}
	transportserverlog.Info("Validation for TransportServer upon deletion", "name", transportserver.GetName())

	// TODO(user): fill in your validation logic upon object deletion.

	return nil, nil
}

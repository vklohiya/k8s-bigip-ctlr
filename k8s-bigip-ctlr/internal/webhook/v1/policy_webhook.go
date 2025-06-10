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
var policylog = logf.Log.WithName("policy-resource")

// SetupPolicyWebhookWithManager registers the webhook for Policy in the manager.
func SetupPolicyWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).For(&cisv1.Policy{}).
		WithValidator(&PolicyCustomValidator{}).
		WithDefaulter(&PolicyCustomDefaulter{}).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-cis-f5-com-v1-policy,mutating=true,failurePolicy=fail,sideEffects=None,groups=cis.f5.com,resources=policies,verbs=create;update,versions=v1,name=mpolicy-v1.kb.io,admissionReviewVersions=v1

// PolicyCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind Policy when those are created or updated.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as it is used only for temporary operations and does not need to be deeply copied.
type PolicyCustomDefaulter struct {
	// TODO(user): Add more fields as needed for defaulting
}

var _ webhook.CustomDefaulter = &PolicyCustomDefaulter{}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the Kind Policy.
func (d *PolicyCustomDefaulter) Default(_ context.Context, obj runtime.Object) error {
	policy, ok := obj.(*cisv1.Policy)

	if !ok {
		return fmt.Errorf("expected an Policy object but got %T", obj)
	}
	policylog.Info("Defaulting for Policy", "name", policy.GetName())

	// TODO(user): fill in your defaulting logic.

	return nil
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// NOTE: The 'path' attribute must follow a specific pattern and should not be modified directly here.
// Modifying the path for an invalid path can cause API server errors; failing to locate the webhook.
// +kubebuilder:webhook:path=/validate-cis-f5-com-v1-policy,mutating=false,failurePolicy=fail,sideEffects=None,groups=cis.f5.com,resources=policies,verbs=create;update,versions=v1,name=vpolicy-v1.kb.io,admissionReviewVersions=v1

// PolicyCustomValidator struct is responsible for validating the Policy resource
// when it is created, updated, or deleted.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
type PolicyCustomValidator struct {
	// TODO(user): Add more fields as needed for validation
}

var _ webhook.CustomValidator = &PolicyCustomValidator{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type Policy.
func (v *PolicyCustomValidator) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	policy, ok := obj.(*cisv1.Policy)
	if !ok {
		return nil, fmt.Errorf("expected a Policy object but got %T", obj)
	}
	policylog.Info("Validation for Policy upon creation", "name", policy.GetName())

	// TODO(user): fill in your validation logic upon object creation.

	return nil, nil
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type Policy.
func (v *PolicyCustomValidator) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	policy, ok := newObj.(*cisv1.Policy)
	if !ok {
		return nil, fmt.Errorf("expected a Policy object for the newObj but got %T", newObj)
	}
	policylog.Info("Validation for Policy upon update", "name", policy.GetName())

	// TODO(user): fill in your validation logic upon object update.

	return nil, nil
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type Policy.
func (v *PolicyCustomValidator) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	policy, ok := obj.(*cisv1.Policy)
	if !ok {
		return nil, fmt.Errorf("expected a Policy object but got %T", obj)
	}
	policylog.Info("Validation for Policy upon deletion", "name", policy.GetName())

	// TODO(user): fill in your validation logic upon object deletion.

	return nil, nil
}

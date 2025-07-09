/*
Copyright 2025.

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
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	awsauthv1beta1 "github.com/redacid/aws-auth-controller/api/v1beta1"
)

// nolint:unused
// log is for logging in this package.
var mapaccountlog = logf.Log.WithName("mapaccount-resource")

// SetupMapAccountWebhookWithManager registers the webhook for MapAccount in the manager.
func SetupMapAccountWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).For(&awsauthv1beta1.MapAccount{}).
		WithValidator(&MapAccountCustomValidator{}).
		WithDefaulter(&MapAccountCustomDefaulter{}).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-aws-auth-prozorro-sale-v1beta1-mapaccount,mutating=true,failurePolicy=fail,sideEffects=None,groups=aws-auth.prozorro.sale,resources=mapaccounts,verbs=create;update,versions=v1beta1,name=mmapaccount-v1beta1.kb.io,admissionReviewVersions=v1

// MapAccountCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind MapAccount when those are created or updated.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as it is used only for temporary operations and does not need to be deeply copied.
type MapAccountCustomDefaulter struct {
	// TODO(user): Add more fields as needed for defaulting
}

var _ webhook.CustomDefaulter = &MapAccountCustomDefaulter{}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the Kind MapAccount.
func (d *MapAccountCustomDefaulter) Default(_ context.Context, obj runtime.Object) error {
	mapaccount, ok := obj.(*awsauthv1beta1.MapAccount)

	if !ok {
		return fmt.Errorf("expected an MapAccount object but got %T", obj)
	}
	mapaccountlog.Info("Defaulting for MapAccount", "name", mapaccount.GetName(), "AccountID", mapaccount.Spec)

	// TODO(user): fill in your defaulting logic.

	return nil
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// NOTE: The 'path' attribute must follow a specific pattern and should not be modified directly here.
// Modifying the path for an invalid path can cause API server errors; failing to locate the webhook.
// +kubebuilder:webhook:path=/validate-aws-auth-prozorro-sale-v1beta1-mapaccount,mutating=false,failurePolicy=fail,sideEffects=None,groups=aws-auth.prozorro.sale,resources=mapaccounts,verbs=create;update;delete,versions=v1beta1,name=vmapaccount-v1beta1.kb.io,admissionReviewVersions=v1

// MapAccountCustomValidator struct is responsible for validating the MapAccount resource
// when it is created, updated, or deleted.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
type MapAccountCustomValidator struct {
	// TODO(user): Add more fields as needed for validation
}

var _ webhook.CustomValidator = &MapAccountCustomValidator{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type MapAccount.
func (v *MapAccountCustomValidator) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	mapaccount, ok := obj.(*awsauthv1beta1.MapAccount)
	if !ok {
		return nil, fmt.Errorf("expected a MapAccount object but got %T", obj)
	}
	mapaccountlog.Info("Validation for MapAccount upon creation", "name", mapaccount.GetName(), "AccountID", mapaccount.Spec)

	// TODO(user): fill in your validation logic upon object creation.

	return nil, nil
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type MapAccount.
func (v *MapAccountCustomValidator) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	mapaccount, ok := newObj.(*awsauthv1beta1.MapAccount)
	if !ok {
		return nil, fmt.Errorf("expected a MapAccount object for the newObj but got %T", newObj)
	}
	mapaccountlog.Info("Validation for MapAccount upon update", "name", mapaccount.GetName(), "AccountID", mapaccount.Spec)

	// TODO(user): fill in your validation logic upon object update.

	return nil, nil
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type MapAccount.
func (v *MapAccountCustomValidator) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	mapaccount, ok := obj.(*awsauthv1beta1.MapAccount)
	if !ok {
		return nil, fmt.Errorf("expected a MapAccount object but got %T", obj)
	}
	mapaccountlog.Info("Validation for MapAccount upon deletion", "name", mapaccount.GetName(), "AccountID", mapaccount.Spec)

	// TODO(user): fill in your validation logic upon object deletion.

	return nil, nil
}

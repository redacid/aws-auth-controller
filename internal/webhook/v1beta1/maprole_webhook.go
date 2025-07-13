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
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	awsauthv1beta1 "github.com/redacid/aws-auth-controller/api/v1beta1"
	"github.com/redacid/aws-auth-controller/awsauth"
)

// nolint:unused
// log is for logging in this package.
var maprolelog = logf.Log.WithName("maprole-resource")

// SetupMapRoleWebhookWithManager registers the webhook for MapRole in the manager.
func SetupMapRoleWebhookWithManager(mgr ctrl.Manager) error {
	validator := &MapRoleCustomValidator{
		Client: mgr.GetClient(),
	}

	return ctrl.NewWebhookManagedBy(mgr).For(&awsauthv1beta1.MapRole{}).
		// WithValidator(&MapRoleCustomValidator{}).
		WithValidator(validator).
		WithDefaulter(&MapRoleCustomDefaulter{}).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-aws-auth-prozorro-sale-v1beta1-maprole,mutating=true,failurePolicy=fail,sideEffects=None,groups=aws-auth.prozorro.sale,resources=maproles,verbs=create;update,versions=v1beta1,name=mmaprole-v1beta1.kb.io,admissionReviewVersions=v1

// MapRoleCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind MapRole when those are created or updated.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as it is used only for temporary operations and does not need to be deeply copied.
type MapRoleCustomDefaulter struct {
	// TODO(user): Add more fields as needed for defaulting
}

var _ webhook.CustomDefaulter = &MapRoleCustomDefaulter{}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the Kind MapRole.
func (d *MapRoleCustomDefaulter) Default(_ context.Context, obj runtime.Object) error {
	maprole, ok := obj.(*awsauthv1beta1.MapRole)

	if !ok {
		return fmt.Errorf("expected an MapRole object but got %T", obj)
	}
	maprolelog.Info("Defaulting for MapRole", "name", maprole.GetName())

	// TODO(user): fill in your defaulting logic.

	return nil
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// NOTE: The 'path' attribute must follow a specific pattern and should not be modified directly here.
// Modifying the path for an invalid path can cause API server errors; failing to locate the webhook.
// +kubebuilder:webhook:path=/validate-aws-auth-prozorro-sale-v1beta1-maprole,mutating=false,failurePolicy=fail,sideEffects=None,groups=aws-auth.prozorro.sale,resources=maproles,verbs=create;update;delete,versions=v1beta1,name=vmaprole-v1beta1.kb.io,admissionReviewVersions=v1

// MapRoleCustomValidator struct is responsible for validating the MapRole resource
// when it is created, updated, or deleted.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
type MapRoleCustomValidator struct {
	Client client.Client
}

var _ webhook.CustomValidator = &MapRoleCustomValidator{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type MapRole.
func (v *MapRoleCustomValidator) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	// This wait need for write to configmap
	time.Sleep(100 * time.Millisecond)

	maprole, ok := obj.(*awsauthv1beta1.MapRole)
	if !ok {
		return nil, fmt.Errorf("expected a MapRole object but got %T", obj)
	}
	maprolelog.Info("Validation for MapRole upon creation",
		"name", maprole.GetName(),
		"RoleARN", maprole.Spec.RoleARN,
		"Username", maprole.Spec.Username,
		"Description", maprole.Spec.Description,
		"Groups", maprole.Spec.Groups,
		"Namespace", maprole.GetNamespace(),
	)

	if awsauth.CrdItemAllowedNamespace != "" {
		if maprole.GetNamespace() != awsauth.CrdItemAllowedNamespace {
			mapuserlog.Error(nil, "Namespace "+maprole.GetNamespace()+" is NOT allowed for creation MapRole")
			return nil, fmt.Errorf("namespace %s is NOT allowed for creation MapRole", maprole.GetNamespace())
		}
	}

	// -----
	var mapRoleList awsauthv1beta1.MapRoleList
	if err := v.Client.List(ctx, &mapRoleList); err != nil {
		return nil, err
	}
	for _, existingUser := range mapRoleList.Items {
		if existingUser.Spec.Username+existingUser.Spec.RoleARN == maprole.Spec.Username+maprole.Spec.RoleARN {
			return nil, fmt.Errorf("duplicate user data found: username %s or rolearn %s already exists in another MapRole resource: %s",
				maprole.Spec.Username, maprole.Spec.RoleARN, existingUser.GetName())
		}
	}
	// -----

	/*	kubeClient, err := kube.GetClient()
		if err != nil {
			mapuserlog.Error(err, "Failure getting kube client")
			return nil, err
		}

		// Get a new aws auth service object.
		awsauthSvc, err := awsauth.NewService(&awsauth.ServiceConfig{
			KubeClient:    kubeClient,
			Log:           ctrl.Log,
			MaxRetryCount: 5,
		})
		if err != nil {
			mapuserlog.Error(err, "Failure creating new aws auth service")
			return nil, err
		}

		if err := awsauthSvc.CheckMapRoleExists(awsauth.MapRole{
			Username: maprole.Spec.Username,
			RoleARN:  maprole.Spec.RoleARN,
			Groups:   maprole.Spec.Groups,
		}); err != nil {
			mapuserlog.Info("Failure checking, username or userarn exists in aws-auth configmap")
			return nil, fmt.Errorf("failure checking, username %v or rolearn %v exists in aws-auth configmap", maprole.Spec.Username, maprole.Spec.RoleARN)
		} else {
			mapuserlog.Info("username and rolearn not exists in aws-auth configmap")
		}*/

	return nil, nil
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type MapRole.
func (v *MapRoleCustomValidator) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	maprole, ok := newObj.(*awsauthv1beta1.MapRole)
	oldmaprole, _ := oldObj.(*awsauthv1beta1.MapRole)
	if !ok {
		return nil, fmt.Errorf("expected a MapRole object for the newObj but got %T", newObj)
	}
	maprolelog.Info("Validation for MapUser upon update OLD:",
		"name", oldmaprole.GetName(),
		"RoleARN", oldmaprole.Spec.RoleARN,
		"Username", oldmaprole.Spec.Username,
		"Description", oldmaprole.Spec.Description,
		"Groups", oldmaprole.Spec.Groups,
		"Namespace", oldmaprole.GetNamespace(),
	)
	maprolelog.Info("Validation for MapUser upon update NEW:",
		"name", maprole.GetName(),
		"RoleARN", maprole.Spec.RoleARN,
		"Username", maprole.Spec.Username,
		"Description", maprole.Spec.Description,
		"Groups", maprole.Spec.Groups,
		"Namespace", maprole.GetNamespace(),
	)

	if (maprole.Spec.Username != oldmaprole.Spec.Username) && (maprole.Spec.Username != "") {
		return nil, fmt.Errorf("username cannot be changed, pls create a new MapRole with the new Username, only change Groups allowed")
	}
	if (maprole.Spec.RoleARN != oldmaprole.Spec.RoleARN) && (maprole.Spec.RoleARN != "") {
		return nil, fmt.Errorf("RoleARN cannot be changed, pls create a new MapRole with the new RoleARN, only change Groups allowed")
	}

	if err := awsauth.VerifyGroups(maprole.Spec.Groups); err != nil {
		return nil, err
	}

	return nil, nil
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type MapRole.
func (v *MapRoleCustomValidator) ValidateDelete(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	maprole, ok := obj.(*awsauthv1beta1.MapRole)
	if !ok {
		return nil, fmt.Errorf("expected a MapRole object but got %T", obj)
	}
	maprolelog.Info("Validation for MapRole upon deletion", "name", maprole.GetName())

	// TODO(user): fill in your validation logic upon object deletion.

	return nil, nil
}

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

	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	awsauthv1beta1 "github.com/redacid/aws-auth-controller/api/v1beta1"
	"github.com/redacid/aws-auth-controller/awsauth"
	"github.com/redacid/aws-auth-controller/kube"
)

// nolint:unused
// log is for logging in this package.
var (
	mapuserlog = logf.Log.WithName("mapuser-resource")
)

// SetupMapUserWebhookWithManager registers the webhook for MapUser in the manager.
func SetupMapUserWebhookWithManager(mgr ctrl.Manager) error {
	validator := &MapUserCustomValidator{
		Client: mgr.GetClient(),
	}

	return ctrl.NewWebhookManagedBy(mgr).For(&awsauthv1beta1.MapUser{}).
		// WithValidator(&MapUserCustomValidator{}).
		WithValidator(validator).
		WithDefaulter(&MapUserCustomDefaulter{}).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-aws-auth-prozorro-sale-v1beta1-mapuser,mutating=true,failurePolicy=fail,sideEffects=None,groups=aws-auth.prozorro.sale,resources=mapusers,verbs=create;update,versions=v1beta1,name=mmapuser-v1beta1.kb.io,admissionReviewVersions=v1

// MapUserCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind MapUser when those are created or updated.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as it is used only for temporary operations and does not need to be deeply copied.
type MapUserCustomDefaulter struct {
	// TODO(user): Add more fields as needed for defaulting
}

var _ webhook.CustomDefaulter = &MapUserCustomDefaulter{}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the Kind MapUser.
func (d *MapUserCustomDefaulter) Default(_ context.Context, obj runtime.Object) error {
	mapuser, ok := obj.(*awsauthv1beta1.MapUser)

	if !ok {
		return fmt.Errorf("expected an MapUser object but got %T", obj)
	}
	mapuserlog.Info("Defaulting for MapUser", "name", mapuser.GetName())

	return nil
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// NOTE: The 'path' attribute must follow a specific pattern and should not be modified directly here.
// Modifying the path for an invalid path can cause API server errors; failing to locate the webhook.
// +kubebuilder:webhook:path=/validate-aws-auth-prozorro-sale-v1beta1-mapuser,mutating=false,failurePolicy=fail,sideEffects=None,groups=aws-auth.prozorro.sale,resources=mapusers,verbs=create;update;delete,versions=v1beta1,name=vmapuser-v1beta1.kb.io,admissionReviewVersions=v1

// MapUserCustomValidator struct is responsible for validating the MapUser resource
// when it is created, updated, or deleted.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
type MapUserCustomValidator struct {
	Client client.Client
}

var _ webhook.CustomValidator = &MapUserCustomValidator{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type MapUser.
func (v *MapUserCustomValidator) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	// This wait need for write to configmap
	time.Sleep(100 * time.Millisecond)

	mapuser, ok := obj.(*awsauthv1beta1.MapUser)
	if !ok {
		return nil, fmt.Errorf("expected a MapUser object but got %T", obj)
	}
	mapuserlog.Info("Validation for MapUser upon creation",
		"name", mapuser.GetName(),
		"UserARN", mapuser.Spec.UserARN,
		"Username", mapuser.Spec.Username,
		"Description", mapuser.Spec.Description,
		"Groups", mapuser.Spec.Groups,
		"Namespace", mapuser.GetNamespace(),
	)

	if awsauth.CrdItemAllowedNamespace != "" {
		if mapuser.GetNamespace() != awsauth.CrdItemAllowedNamespace {
			mapuserlog.Error(nil, "Namespace "+mapuser.GetNamespace()+" is NOT allowed for creation MapUser")
			return nil, fmt.Errorf("namespace %s is NOT allowed for creation MapUser", mapuser.GetNamespace())
		}
	}
	// -----
	var mapUserList awsauthv1beta1.MapUserList
	if err := v.Client.List(ctx, &mapUserList); err != nil {
		return nil, err
	}
	for _, existingUser := range mapUserList.Items {
		if existingUser.Spec.Username == mapuser.Spec.Username ||
			existingUser.Spec.UserARN == mapuser.Spec.UserARN {
			return nil, fmt.Errorf("duplicate user data found: username %s or userarn %s already exists in another MapUser resource: %s",
				mapuser.Spec.Username, mapuser.Spec.UserARN, existingUser.GetName())
		}
	}
	// -----

	kubeClient, err := kube.GetClient()
	if err != nil {
		mapuserlog.Error(err, "Failure getting kube client")
		return nil, err
	}

	// Get a new aws auth service object.
	awsauthSvc, err := awsauth.NewService(&awsauth.ServiceConfig{
		KubeClient: kubeClient,
		Log:        ctrl.Log,
	})
	if err != nil {
		mapuserlog.Error(err, "Failure creating new aws auth service")
		return nil, err
	}

	if err := awsauth.VerifyUsername(mapuser.Spec.Username, awsauth.UsernameMustBeEmail); err != nil {
		return nil, err
	}

	if err := awsauthSvc.CheckMapUserExists(awsauth.MapUser{
		Username: mapuser.Spec.Username,
		UserARN:  mapuser.Spec.UserARN,
		Groups:   mapuser.Spec.Groups,
	}); err != nil {
		mapuserlog.Info("Failure checking, username or userarn exists in aws-auth configmap")
		return nil, fmt.Errorf("failure checking, username %v or userarn %v exists in aws-auth configmap", mapuser.Spec.Username, mapuser.Spec.UserARN)
	} else {
		mapuserlog.Info("username and userarn not exists in aws-auth configmap")
	}

	return nil, nil
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type MapUser.
func (v *MapUserCustomValidator) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	mapuser, ok := newObj.(*awsauthv1beta1.MapUser)
	oldmapuser, _ := oldObj.(*awsauthv1beta1.MapUser)

	if !ok {
		return nil, fmt.Errorf("expected a MapUser object for the newObj but got %T", newObj)
	}
	mapuserlog.Info("Validation for MapUser upon update OLD:",
		"name", oldmapuser.GetName(),
		"UserARN", oldmapuser.Spec.UserARN,
		"Username", oldmapuser.Spec.Username,
		"Description", oldmapuser.Spec.Description,
		"Groups", oldmapuser.Spec.Groups,
		"Namespace", mapuser.GetNamespace(),
	)
	mapuserlog.Info("Validation for MapUser upon update NEW:",
		"name", mapuser.GetName(),
		"UserARN", mapuser.Spec.UserARN,
		"Username", mapuser.Spec.Username,
		"Description", mapuser.Spec.Description,
		"Groups", mapuser.Spec.Groups,
		"Namespace", mapuser.GetNamespace(),
	)

	if (mapuser.Spec.Username != oldmapuser.Spec.Username) && (mapuser.Spec.Username != "") {
		return nil, fmt.Errorf("username cannot be changed, pls create a new MapUser with the new Username, only change Groups allowed")
	}
	if (mapuser.Spec.UserARN != oldmapuser.Spec.UserARN) && (mapuser.Spec.UserARN != "") {
		return nil, fmt.Errorf("UserARN cannot be changed, pls create a new MapUser with the new UserARN, only change Groups allowed")
	}

	if err := awsauth.VerifyGroups(mapuser.Spec.Groups); err != nil {
		return nil, err
	}

	return nil, nil
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type MapUser.
func (v *MapUserCustomValidator) ValidateDelete(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	mapuser, ok := obj.(*awsauthv1beta1.MapUser)
	if !ok {
		return nil, fmt.Errorf("expected a MapUser object but got %T", obj)
	}
	mapuserlog.Info("Validation for MapUser upon deletion",
		"name", mapuser.GetName(),
		"UserARN", mapuser.Spec.UserARN,
		"Username", mapuser.Spec.Username,
		"Description", mapuser.Spec.Description,
		"Groups", mapuser.Spec.Groups,
	)

	mapuserlog.Info("Validation for MapUser upon deletion",
		"name", mapuser.GetName(),
		"UserARN", mapuser.Spec.UserARN,
	)

	return nil, nil
}

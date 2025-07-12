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

package controller

import (
	"context"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	awsauthv1beta1 "github.com/redacid/aws-auth-controller/api/v1beta1"
	"github.com/redacid/aws-auth-controller/awsauth"
	"github.com/redacid/aws-auth-controller/kube"
)

// MapUserReconciler reconciles a MapUser object
type MapUserReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=aws-auth.prozorro.sale,resources=mapusers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=aws-auth.prozorro.sale,resources=mapusers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=aws-auth.prozorro.sale,resources=mapusers/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;update;patch;delete,resourceNames=aws-auth

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the MapUser object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/reconcile
func (r *MapUserReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = logf.FromContext(ctx)

	// MapUser objects are named by their associated AWS IAM user ARNs.
	// mapUserName := req.Name
	mapUserName := req.Name
	log := ctrl.Log.WithValues("MapUser", mapUserName)
	// log := r. Log.WithValues("MapUser", mapUserName)
	log.Info("")

	kubeClient, err := kube.GetClient()
	if err != nil {
		log.Error(err, "failure getting kube client")
		return ctrl.Result{}, err
	}

	// Get a new aws auth service object.
	awsauthSvc, err := awsauth.NewService(&awsauth.ServiceConfig{
		KubeClient: kubeClient,
		// Log:        r.Log,
		Log:           ctrl.Log,
		MaxRetryCount: 5,
	})
	if err != nil {
		log.Error(err, "failure creating new aws auth service")
		return ctrl.Result{}, err
	}

	// Load the MapUser object by name (its AWS IAM user ARN).
	mapUser := &awsauthv1beta1.MapUser{}

	if err := r.Get(ctx, req.NamespacedName, mapUser); err != nil {
		// If any error other than a "NotFound" API error, it's a problem.
		statusErr, ok := err.(*apierrors.StatusError)
		if !ok || (ok && statusErr.ErrStatus.Reason != NotFound) {
			logf.Log.Error(err, "failure getting MapUser")
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}
	if mapUser.DeletionTimestamp.IsZero() {
		logf.Log.Info("mapAccount is not being deleted")
		// Add finalizer
		if !controllerutil.ContainsFinalizer(mapUser, awsauth.CrdFinalizerName) {
			logf.Log.Info("mapAccount is not being deleted, so adding finalizer")
			controllerutil.AddFinalizer(mapUser, awsauth.CrdFinalizerName)
			if err := r.Update(ctx, mapUser); err != nil {
				return ctrl.Result{}, err
			}
		} else {
			// Ensure that any changes are synced to the kube-system:aws-auth ConfigMap.
			if err := awsauthSvc.UpsertMapUser(awsauth.MapUser{
				Username: mapUser.Spec.Username,
				UserARN:  mapUser.Spec.UserARN,
				Groups:   mapUser.Spec.Groups,
			}); err != nil {
				log.Error(err, "failure upserting MapUser")
				return ctrl.Result{}, err
			}
			log.Info("upserted MapUser")
		}
	} else {
		log.Info("mapUser is being deleted")
		if controllerutil.ContainsFinalizer(mapUser, awsauth.CrdFinalizerName) {
			if err := awsauthSvc.RemoveMapUser(awsauth.MapUser{
				Username: mapUser.Spec.Username,
				UserARN:  mapUser.Spec.UserARN,
				Groups:   mapUser.Spec.Groups,
			}); err != nil {
				log.Error(err, "failure removing mapUser data in aws-auth configmap")
				return ctrl.Result{}, nil
			}
			logf.Log.Info("mapUser is being deleted, so removing finalizer")
			controllerutil.RemoveFinalizer(mapUser, awsauth.CrdFinalizerName)
			if err := r.Update(ctx, mapUser); err != nil {
				return ctrl.Result{}, err
			}
		}
		log.Info("removed mapUser data in aws-auth configmap")
		return ctrl.Result{}, nil
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MapUserReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&awsauthv1beta1.MapUser{}).
		Named("mapuser").
		Complete(r)
}

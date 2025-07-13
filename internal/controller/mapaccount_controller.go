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
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	awsauthv1beta1 "github.com/redacid/aws-auth-controller/api/v1beta1"
	"github.com/redacid/aws-auth-controller/awsauth"
	"github.com/redacid/aws-auth-controller/kube"
)

// MapAccountReconciler reconciles a MapAccount object
type MapAccountReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=aws-auth.prozorro.sale,resources=mapaccounts,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=aws-auth.prozorro.sale,resources=mapaccounts/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=aws-auth.prozorro.sale,resources=mapaccounts/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;update;patch;delete,resourceNames=aws-auth

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the MapAccount object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/reconcile
func (r *MapAccountReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = logf.FromContext(ctx)

	// TODO(user): your logic here

	// MapAccount objects a list of AWS Accounts.
	mapAccountName := req.Name
	log := ctrl.Log.WithValues("----------MapAccount", mapAccountName)

	log.Info("---------------------------------------------------")

	kubeClient, err := kube.GetClient()
	if err != nil {
		log.Error(err, "Failure getting kube client")
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
		log.Error(err, "Failure creating new aws auth service")
		return ctrl.Result{}, err
	}

	// Load the MapUser object by name (its AWS IAM user ARN).
	mapAccount := &awsauthv1beta1.MapAccount{}

	if err := r.Get(ctx, req.NamespacedName, mapAccount); err != nil {
		// If any error other than a "NotFound" API error, it's a problem.
		statusErr, ok := err.(*apierrors.StatusError)
		if !ok || (ok && statusErr.ErrStatus.Reason != NotFound) {
			logf.Log.Error(err, "Failure getting mapAccount")
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}
	// examine DeletionTimestamp to determine if object is under deletion
	if mapAccount.DeletionTimestamp.IsZero() {
		logf.Log.Info("mapAccount is not being deleted")
		// Add finalizer
		if !controllerutil.ContainsFinalizer(mapAccount, awsauth.CrdFinalizerName) {
			logf.Log.Info("mapAccount is not being deleted, so adding finalizer")
			controllerutil.AddFinalizer(mapAccount, awsauth.CrdFinalizerName)
			if err := r.Update(ctx, mapAccount); err != nil {
				return ctrl.Result{}, err
			}
		} else {
			// Ensure that any changes are synced to the kube-system:aws-auth ConfigMap.
			if err := awsauthSvc.UpsertMapAccount(awsauth.MapAccount{
				AccountID: mapAccount.Spec.AccountID,
			}); err != nil {
				log.Error(err, "Failure upserting MapAccount")
				return ctrl.Result{}, err
			}
			log.Info("Upserted MapAccount")
		}

	} else {
		log.Info("mapAccount is being deleted")
		if controllerutil.ContainsFinalizer(mapAccount, awsauth.CrdFinalizerName) {
			if err := awsauthSvc.RemoveMapAccount(awsauth.MapAccount{
				AccountID: mapAccount.Spec.AccountID,
			}); err != nil {
				log.Error(err, "Failure removing mapAccount data in aws-auth configmap")
				return ctrl.Result{}, nil
			}
			logf.Log.Info("Removing finalizer")
			controllerutil.RemoveFinalizer(mapAccount, awsauth.CrdFinalizerName)
			if err := r.Update(ctx, mapAccount); err != nil {
				return ctrl.Result{}, err
			}
		}
		log.Info("Removed mapAccount data in aws-auth configmap")
		return ctrl.Result{}, nil
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MapAccountReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&awsauthv1beta1.MapAccount{}).
		Named("mapaccount").
		Complete(r)
}

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
	"time"

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

// MapRoleReconciler reconciles a MapRole object
type MapRoleReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=aws-auth.prozorro.sale,resources=maproles,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=aws-auth.prozorro.sale,resources=maproles/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=aws-auth.prozorro.sale,resources=maproles/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;update;patch;delete,resourceNames=aws-auth

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the MapRole object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/reconcile
func (r *MapRoleReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = logf.FromContext(ctx)

	log := ctrl.Log.WithValues("MapRole", req.Name, "namespace", req.Namespace)
	log.Info("Reconciling MapRole")

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
		WithRetries:   true,
		MaxRetryCount: 5,
		MinRetryTime:  time.Millisecond * 100,
		MaxRetryTime:  time.Second * 3,
	})
	if err != nil {
		log.Error(err, "Failure creating new aws auth service")
		return ctrl.Result{}, err
	}

	// Load the MapRole object by name (its AWS IAM role ARN).
	mapRole := &awsauthv1beta1.MapRole{}
	if err := r.Get(ctx, req.NamespacedName, mapRole); err != nil {
		// If any error other than a "NotFound" API error, it's a problem.
		statusErr, ok := err.(*apierrors.StatusError)
		if !ok || (ok && statusErr.ErrStatus.Reason != NotFound) {
			log.Error(err, "Failure getting MapRole")
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	if mapRole.DeletionTimestamp.IsZero() {
		log.V(1).Info("mapRole is not being deleted: " + mapRole.Name)
		// Add finalizer
		if !controllerutil.ContainsFinalizer(mapRole, awsauth.CrdFinalizerName) {
			log.V(1).Info("Adding finalizer: " + mapRole.Name)
			controllerutil.AddFinalizer(mapRole, awsauth.CrdFinalizerName)
			if err := r.Update(ctx, mapRole); err != nil {
				return ctrl.Result{}, err
			}
		} else {
			// Ensure that any changes are synced to the kube-system:aws-auth ConfigMap.
			if err := awsauthSvc.UpsertMapRole(awsauth.MapRole{
				Username: mapRole.Spec.Username,
				RoleARN:  mapRole.Spec.RoleARN,
				Groups:   mapRole.Spec.Groups,
			}); err != nil {
				log.Error(err, "Failure upserting MapRole")
				return ctrl.Result{}, err
			}
			log.V(1).Info("Upserted MapRole")
		}
	} else {
		log.Info("mapRole is being deleted")
		if controllerutil.ContainsFinalizer(mapRole, awsauth.CrdFinalizerName) {
			if err := awsauthSvc.RemoveMapRole(awsauth.MapRole{
				Username: mapRole.Spec.Username,
				RoleARN:  mapRole.Spec.RoleARN,
				Groups:   mapRole.Spec.Groups,
			}); err != nil {
				log.Error(err, "Failure removing mapRole data in aws-auth configmap")
				return ctrl.Result{}, nil
			}
			log.Info("Removing finalizer")
			controllerutil.RemoveFinalizer(mapRole, awsauth.CrdFinalizerName)
			if err := r.Update(ctx, mapRole); err != nil {
				return ctrl.Result{}, err
			}
		}
		log.Info("Removed mapRole data in aws-auth configmap")
		return ctrl.Result{}, nil
	}

	return ctrl.Result{RequeueAfter: awsauth.ReconcileTime}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MapRoleReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&awsauthv1beta1.MapRole{}).
		Named("maprole").
		Complete(r)
}

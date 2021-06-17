/*
Copyright 2021 Dharmjit.

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

package controllers

import (
	"context"
	"reflect"
	"time"

	"dharmjit.dev/copyimage/utils"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// CopyImageDeploymentReconciler reconciles a CopyImageDeployment object
type CopyImageDeploymentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=dsapps.dharmjit.dev,resources=copyimagedeployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=dsapps.dharmjit.dev,resources=copyimagedeployments/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=dsapps.dharmjit.dev,resources=copyimagedeployments/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the CopyImageDeployment object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *CopyImageDeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// your logic here
	if req.Namespace == "kube-system" {
		return reconcile.Result{}, nil
	}
	logger.Info("Event Received", "Info", req.NamespacedName)
	deployment := &v1.Deployment{}
	err := r.Get(ctx, req.NamespacedName, deployment)
	if err != nil {
		return reconcile.Result{}, err
	}

	podspec, err := utils.CloneImage(logger, &deployment.Spec.Template)
	// TODO better error handling
	// currently we ignore any errors and proceed with deployment creation with applied state
	if err != nil {
		logger.Error(err, "Error occured in cloneImage")
		return reconcile.Result{}, nil
	}

	// update only if there are changes in the spec
	if !reflect.DeepEqual(deployment.Spec.Template.Spec, *podspec) {
		deployment.Spec.Template = *podspec
		//TODO add rollout message
		err = r.Client.Update(ctx, deployment)
		if err != nil {
			//TODO handle errors better for optimistic concurrency
			return reconcile.Result{RequeueAfter: 1 * time.Second}, nil
		}
	} else {
		return reconcile.Result{}, nil
	}
	return ctrl.Result{RequeueAfter: 5 * time.Minute}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *CopyImageDeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1.Deployment{}).
		WithEventFilter(predicate.Funcs{
			UpdateFunc: func(e event.UpdateEvent) bool {
				// Suppress Update events to avoid filtering them out in the Reconcile function
				return false
			},
		}).
		WithEventFilter(predicate.Funcs{
			DeleteFunc: func(e event.DeleteEvent) bool {
				// Suppress Delete events to avoid filtering them out in the Reconcile function
				return false
			},
		}).
		Complete(r)
}

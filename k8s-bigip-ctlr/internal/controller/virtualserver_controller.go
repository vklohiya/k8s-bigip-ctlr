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

package controller

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	cisv1 "github.com/F5Networks/k8s-bigip-ctlr/v2/api/v1"
)

// VirtualServerReconciler reconciles a VirtualServer object
type VirtualServerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=cis.f5.com,resources=virtualservers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cis.f5.com,resources=virtualservers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cis.f5.com,resources=virtualservers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the VirtualServer object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/reconcile
func (r *VirtualServerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = logf.FromContext(ctx)
	var svc cisv1.TLSProfile
	_ = r.Client.Get(context.TODO(), client.ObjectKey{
		Namespace: req.Namespace,
		Name:      req.Name,
	}, &svc)
	// TODO(user): your logic here

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *VirtualServerReconciler) SetupWithManager(mgr ctrl.Manager, namespaceSelector labels.Selector) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cisv1.VirtualServer{}).
		WithEventFilter(NamespaceLabelPredicate(mgr.GetClient(), namespaceSelector)).
		WithEventFilter(VirtualServerPredicates(mgr.GetClient())).
		Watches(&corev1.Secret{}, handler.EnqueueRequestsFromMapFunc(
			func(ctx context.Context, obj client.Object) []reconcile.Request {
				return ReconcileVirtualServersForSecret(mgr.GetClient(), obj)
			})).
		Watches(&cisv1.TLSProfile{}, handler.EnqueueRequestsFromMapFunc(
			func(ctx context.Context, obj client.Object) []reconcile.Request {
				return ReconcileVirtualServersForTLSProfile(mgr.GetClient(), obj)
			})).
		Named("virtualserver").
		Complete(r)
}

// VirtualServerPredicates defines event filters for VirtualServer
func VirtualServerPredicates(cli client.Client) predicate.Predicate {
	return predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			// Reconcile on VirtualServer updates
			return true
		},
		CreateFunc: func(e event.CreateEvent) bool {
			return true
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			return true
		},
		GenericFunc: func(e event.GenericEvent) bool {
			return true
		},
	}
}

// ReconcileVirtualServersForSecret maps Secret updates to VirtualServer reconciliation requests
func ReconcileVirtualServersForSecret(cli client.Client, obj client.Object) []ctrl.Request {
	secret := obj.(*corev1.Secret)
	var requests []ctrl.Request

	// Find all TLSProfiles referencing this Secret
	tlsProfileList := &cisv1.TLSProfileList{} // Adjust type based on your API
	if err := cli.List(context.Background(), tlsProfileList); err != nil {
		//log.Error(err, "failed to list TLSProfiles")
		return requests
	}

	for _, tlsProfile := range tlsProfileList.Items {
		if tlsProfile.Spec.TLS.ServerSSL == secret.Name && tlsProfile.Namespace == secret.Namespace {
			// Find VirtualServers referencing this TLSProfile
			vsList := &cisv1.VirtualServerList{}
			if err := cli.List(context.Background(), vsList); err != nil {
				//log.Error(err, "failed to list VirtualServers")
				continue
			}
			for _, vs := range vsList.Items {
				if vs.Spec.TLSProfileName == tlsProfile.Name && vs.Namespace == tlsProfile.Namespace {
					requests = append(requests, ctrl.Request{
						NamespacedName: types.NamespacedName{
							Name:      vs.Name,
							Namespace: vs.Namespace,
						},
					})
				}
			}
		}
	}
	return requests
}

// ReconcileVirtualServersForTLSProfile maps Secret updates to VirtualServer reconciliation requests
func ReconcileVirtualServersForTLSProfile(cli client.Client, obj client.Object) []ctrl.Request {
	secret := obj.(*corev1.Secret)
	var requests []ctrl.Request

	// Find all TLSProfiles referencing this Secret
	tlsProfileList := &cisv1.TLSProfileList{} // Adjust type based on your API
	if err := cli.List(context.Background(), tlsProfileList); err != nil {
		//log.Error(err, "failed to list TLSProfiles")
		return requests
	}

	for _, tlsProfile := range tlsProfileList.Items {
		if tlsProfile.Spec.TLS.ServerSSL == secret.Name && tlsProfile.Namespace == secret.Namespace {
			// Find VirtualServers referencing this TLSProfile
			vsList := &cisv1.VirtualServerList{}
			if err := cli.List(context.Background(), vsList); err != nil {
				//log.Error(err, "failed to list VirtualServers")
				continue
			}
			for _, vs := range vsList.Items {
				if vs.Spec.TLSProfileName == tlsProfile.Name && vs.Namespace == tlsProfile.Namespace {
					requests = append(requests, ctrl.Request{
						NamespacedName: types.NamespacedName{
							Name:      vs.Name,
							Namespace: vs.Namespace,
						},
					})
				}
			}
		}
	}
	return requests
}

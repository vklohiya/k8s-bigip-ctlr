package controller

import (
	"context"
	"fmt"
	cisv1 "github.com/F5Networks/k8s-bigip-ctlr/v2/api/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"time"
)

// NamespaceReconciler reconciles namespaces to trigger CR reconciliation
type NamespaceReconciler struct {
	client.Client
	Scheme            *runtime.Scheme
	NamespaceSelector labels.Selector
}

func (r *NamespaceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := logf.FromContext(ctx)
	var ns corev1.Namespace
	if err := r.Get(ctx, req.NamespacedName, &ns); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if !r.NamespaceSelector.Matches(labels.Set(ns.Labels)) {
		return ctrl.Result{}, nil
	}

	crLists := []client.ObjectList{
		&cisv1.VirtualServerList{},
		&cisv1.TLSProfileList{},
		&cisv1.IngressLinkList{},
		&cisv1.TransportServerList{},
		&cisv1.ExternalDNSList{},
		&cisv1.PolicyList{},
	}

	for _, list := range crLists {
		if err := r.List(ctx, list, client.InNamespace(ns.Name)); err != nil {
			logger.Error(err, "Failed to list CRs", "type", fmt.Sprintf("%T", list), "namespace", ns.Name)
			continue
		}
		// Use TypeMeta to handle list items generically
		items, err := meta.ExtractList(list)
		if err != nil {
			logger.Error(err, "Failed to extract list", "type", fmt.Sprintf("%T", list))
			continue
		}
		for _, item := range items {
			obj := item.(client.Object)
			if obj.GetAnnotations() == nil {
				obj.SetAnnotations(make(map[string]string))
			}
			obj.GetAnnotations()["cis.f5.com/last-reconciled"] = time.Now().String()
			if err := r.Update(ctx, obj); err != nil {
				logger.Error(err, "Failed to update CR", "object", obj.GetName(), "namespace", obj.GetNamespace())
			}
		}
	}
	return ctrl.Result{}, nil
}

func (r *NamespaceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Namespace{}).
		WithEventFilter(predicate.Funcs{
			CreateFunc: func(e event.CreateEvent) bool {
				ns := e.Object.(*corev1.Namespace)
				return r.NamespaceSelector.Matches(labels.Set(ns.Labels))
			},
			UpdateFunc: func(e event.UpdateEvent) bool {
				ns := e.ObjectNew.(*corev1.Namespace)
				return r.NamespaceSelector.Matches(labels.Set(ns.Labels))
			},
			DeleteFunc: func(e event.DeleteEvent) bool {
				return false // No action on delete
			},
			GenericFunc: func(e event.GenericEvent) bool {
				ns := e.Object.(*corev1.Namespace)
				return r.NamespaceSelector.Matches(labels.Set(ns.Labels))
			},
		}).
		Complete(r)
}

// NamespaceLabelPredicate filters events based on namespace labels
func NamespaceLabelPredicate(cli client.Client, selector labels.Selector) predicate.Predicate {
	return predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			return checkNamespaceLabel(cli, e.Object.GetNamespace(), selector)
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			return checkNamespaceLabel(cli, e.ObjectNew.GetNamespace(), selector)
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			return checkNamespaceLabel(cli, e.Object.GetNamespace(), selector)
		},
		GenericFunc: func(e event.GenericEvent) bool {
			return checkNamespaceLabel(cli, e.Object.GetNamespace(), selector)
		},
	}
}

// checkNamespaceLabel checks if a namespace matches the label selector
func checkNamespaceLabel(cli client.Client, ns string, selector labels.Selector) bool {
	if ns == "" {
		return false // Skip cluster-scoped resources
	}
	var namespace corev1.Namespace
	err := cli.Get(ctrl.SetupSignalHandler(), types.NamespacedName{Name: ns}, &namespace)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return false
		}
		return false
	}
	return selector.Matches(labels.Set(namespace.Labels))
}

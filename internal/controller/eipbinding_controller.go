/*
Copyright 2023.

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

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	extensionv1 "github.com/lucheng0127/kube-eip/api/v1"
	corev1 "k8s.io/api/core/v1"
	vmv1 "kubevirt.io/api/core/v1"
)

// EipBindingReconciler reconciles a EipBinding object
type EipBindingReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=extension.my.domain,resources=eipbindings,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=extension.my.domain,resources=eipbindings/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=extension.my.domain,resources=eipbindings/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the EipBinding object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.0/pkg/reconcile
func (r *EipBindingReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Get EipBinding object
	var eb extensionv1.EipBinding
	if err := r.Get(ctx, req.NamespacedName, &eb); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}

		log.Error(err, "unable to fetch EipBinding")
		return ctrl.Result{}, err
	}

	// Init vars
	newHyper, newIPAddr, err := r.getVMIInfo(eb)
	if err != nil {
		if !errors.IsNotFound(err) {
			log.Error(err, "unable to get vmi")
			return ctrl.Result{}, err
		}

		log.Info("VMI offline")
	}

	// Add finalizer and handle delete
	finalizerName := "extension.lucheng0127/finlizer"

	if eb.ObjectMeta.DeletionTimestamp.IsZero() {
		// EipBinding is not been deleted, register finalizer
		if !controllerutil.ContainsFinalizer(&eb, finalizerName) {
			controllerutil.AddFinalizer(&eb, finalizerName)

			if err := r.Update(ctx, &eb); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		// EipBinding is being deleted
		if controllerutil.ContainsFinalizer(&eb, finalizerName) {
			// TODO(shawnlu): do clean up
			log.Info("Clen up EipBinding")
		}

		controllerutil.RemoveFinalizer(&eb, finalizerName)
		if err := r.Update(ctx, &eb); err != nil {
			log.Error(err, "remove finalizer for EipBinding failed")
			return ctrl.Result{}, err
		}

		log.Info("Deleted EipBinding")
		return ctrl.Result{}, nil
	}

	// Work in progress, skip
	if eb.Spec.Phase == extensionv1.PhaseProcessing {
		log.Info("EipBinding work in progress, skip")
		return ctrl.Result{}, nil
	}

	// Do clean up for not ready vmi
	if newIPAddr == "" || newHyper == "" {
		log.Info("Do clean up for not ready vmi")
		// TODO(shawnlu): do clean up
	}

	// Check vmi hyper and ip changed
	if eb.Spec.CurrentHyper == newHyper && eb.Spec.CurrentIPAddr == newIPAddr {
		log.Info("VMI info not changed, skip")
		return ctrl.Result{}, nil
	}

	// Handle vmi info change
	staleHyper := eb.Spec.LastHyper
	staleIPAddr := eb.Spec.LastIPAddr
	eb.Spec.LastHyper = eb.Spec.CurrentHyper
	eb.Spec.LastIPAddr = eb.Spec.CurrentIPAddr
	eb.Spec.CurrentHyper = newHyper
	eb.Spec.CurrentIPAddr = newIPAddr

	// Apply eip bind and unbind
	eb.Spec.Phase = extensionv1.PhaseProcessing
	if err := r.Client.Update(ctx, &eb); err != nil {
		log.Error(err, "update EipBinding phase failed")
		return ctrl.Result{}, err
	}

	if staleHyper != "" && staleIPAddr != "" {
		log.Info("Clean up staled hyper eip rules")
		// TODO(shawnlu): clean up staled eip rules on old hyper
	}

	log.Info("Apply eip rules on hyper")
	// TODO(shawnlu): apply eip rules on new hyper

	eb.Spec.Phase = extensionv1.PhaseReady
	if err := r.Client.Update(ctx, &eb); err != nil {
		log.Error(err, "update EipBinding phase failed")
		return ctrl.Result{}, err
	}

	log.Info("Eip bind succeed")
	return ctrl.Result{}, nil
}

// Get hyper ip address
func (r *EipBindingReconciler) getHyperIPAddr(name string) (string, error) {
	node := &corev1.Node{}
	err := r.Client.Get(context.TODO(), client.ObjectKey{
		Name: name,
	}, node)

	if err != nil {
		return "", err
	}

	return node.Status.Addresses[0].Address, nil
}

// Get EipBinding vmi hyper ip address infos
// When vmi not exist or not running return NotFound error
// and ignore it.
func (r *EipBindingReconciler) getVMIInfo(eb extensionv1.EipBinding) (string, string, error) {
	var vmi vmv1.VirtualMachineInstance
	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      eb.Spec.VmiName,
		Namespace: eb.Namespace,
	}, &vmi)
	if err != nil {
		return "", "", err
	}

	if vmi.Status.Phase != "Running" {
		return "", "", errors.NewNotFound(vmv1.Resource("VirtualMachineInstance"), "VirtualMachineInstance")
	}

	hyperAddr, err := r.getHyperIPAddr(vmi.Status.NodeName)
	if err != nil {
		return "", "", err
	}

	if len(vmi.Status.Interfaces) == 0 {
		return "", "", errors.NewBadRequest("vmi without interface, can't bind")
	}

	return hyperAddr, vmi.Status.Interfaces[0].IP, nil
}

func (r *EipBindingReconciler) findObjectsForVmi(ctx context.Context, vmi client.Object) []reconcile.Request {
	attachedEipBinding := &extensionv1.EipBindingList{}
	listOps := client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(".spec.vmiName", vmi.GetName()),
		Namespace:     vmi.GetNamespace(),
	}

	err := r.List(context.TODO(), attachedEipBinding, &listOps)
	if err != nil {
		return []reconcile.Request{}
	}

	requests := make([]reconcile.Request, len(attachedEipBinding.Items))
	for i, item := range attachedEipBinding.Items {
		requests[i] = reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      item.GetName(),
				Namespace: item.GetNamespace(),
			},
		}
	}

	return requests
}

// SetupWithManager sets up the controller with the Manager.
func (r *EipBindingReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &extensionv1.EipBinding{}, ".spec.vmiName", func(rawObj client.Object) []string {
		eb := rawObj.(*extensionv1.EipBinding)
		if eb.Spec.VmiName == "" {
			return nil
		}
		return []string{eb.Spec.VmiName}
	}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&extensionv1.EipBinding{}).
		Owns(&vmv1.VirtualMachineInstance{}).
		Watches(
			&vmv1.VirtualMachineInstance{},
			handler.EnqueueRequestsFromMapFunc(r.findObjectsForVmi),
		).
		Complete(r)
}

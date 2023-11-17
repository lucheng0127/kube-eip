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

	var eb extensionv1.EipBinding
	if err := r.Get(ctx, req.NamespacedName, &eb); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}

		log.Error(err, "unable to fetch EipBinding")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Handle delete
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
			return ctrl.Result{}, err
		}

		log.Info("Deleted EipBinding")
		return ctrl.Result{}, nil
	}

	// Work in process, skip
	if eb.Spec.Phase == extensionv1.PhaseProcessing {
		log.Info("EipBinding work in process")
		return ctrl.Result{}, nil
	}

	// Get vmi instance
	var vmi vmv1.VirtualMachineInstance
	r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      eb.Spec.VmiName,
		Namespace: eb.Namespace,
	}, &vmi)

	if vmi.Status.Phase != "Running" {
		log.Info("VMI not running")
		return ctrl.Result{}, nil
	}

	if len(vmi.Status.Interfaces) == 0 {
		log.Error(errors.NewBadRequest("vmi withou avaliable interface"), "eip bind failed")
		eb.Spec.Phase = extensionv1.PhaseError
		if err := r.Client.Update(ctx, &eb); err != nil {
			log.Error(err, "update EipBinding failed")
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	currentIPAddr := vmi.Status.Interfaces[0].IP
	currentHyper, err := r.getHyperIPAddr(vmi.Status.NodeName)
	if err != nil {
		log.Error(errors.NewInternalError(err), "get hyper IP address failed")
	}

	// Check vmi ip and hyper changed
	if eb.Spec.CurrentHyper == currentHyper && eb.Spec.CurrentIPAddr == currentIPAddr {
		// Nothing changed, skip
		return ctrl.Result{}, nil
	}

	// Handle create and update
	if eb.Spec.LastIPAddr == "" && eb.Spec.LastHyper == "" {
		// No old rule need be cleaned
		log.Info("Try to bind to vmi")
		eb.Spec.CurrentIPAddr = currentIPAddr
		eb.Spec.CurrentHyper = currentHyper
		eb.Spec.Phase = extensionv1.PhaseProcessing
		if err := r.Client.Update(ctx, &eb); err != nil {
			log.Error(err, "update EipBinding failed")
			return ctrl.Result{}, err
		}
		// TODO(shawnlu): apply eip rules
		eb.Spec.Phase = extensionv1.PhaseReady
		if err := r.Client.Update(ctx, &eb); err != nil {
			log.Error(err, "update EipBinding failed")
			return ctrl.Result{}, err
		}
		log.Info("Eip bind succeed")
		return ctrl.Result{}, nil
	}

	staleHyper := eb.Spec.LastHyper
	staleIPAddr := eb.Spec.LastIPAddr
	eb.Spec.LastHyper = eb.Spec.CurrentHyper
	eb.Spec.LastIPAddr = eb.Spec.CurrentIPAddr
	eb.Spec.CurrentHyper = currentHyper
	eb.Spec.CurrentIPAddr = currentIPAddr
	eb.Spec.Phase = extensionv1.PhaseProcessing
	if err := r.Client.Update(ctx, &eb); err != nil {
		log.Error(err, "update EipBinding failed")
		return ctrl.Result{}, err
	}

	// TODO(shawnlu)： Implement it
	// Clean up first then bind
	log.Info("Clean up last eip rules")
	log.Info(staleHyper)
	log.Info(staleIPAddr)

	// Apply eip rules
	eb.Spec.Phase = extensionv1.PhaseReady
	if err := r.Client.Update(ctx, &eb); err != nil {
		log.Error(err, "update EipBinding failed")
		return ctrl.Result{}, err
	}
	log.Info("Eip bind succeed")

	return ctrl.Result{}, nil
}

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

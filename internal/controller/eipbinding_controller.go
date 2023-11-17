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

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	extensionv1 "github.com/lucheng0127/kube-eip/api/v1"
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
		log.Error(err, "unable to fetch EipBinding")
		return ctrl.Result{}, err
	}
	// TODO(shawnlu): Implement it

	return ctrl.Result{}, nil
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

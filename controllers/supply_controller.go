/*

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
	"fmt"

	"github.com/go-logr/logr"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	fundv1beta1 "github.com/YaoZengzeng/variance/api/v1beta1"
	"github.com/YaoZengzeng/variance/types"
	"github.com/YaoZengzeng/variance/variance"
)

// SupplyReconciler reconciles a Supply object
type SupplyReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

/*
We generally want to ignore (not requeue) NotFound errors, since we'll get a
reconciliation request once the object exists, and requeuing in the meantime
won't help.
*/
func ignoreNotFound(err error) error {
	if apierrs.IsNotFound(err) {
		return nil
	}
	return err
}

// +kubebuilder:docs-gen:collapse=ignoreNotFound

// +kubebuilder:rbac:groups=fund.example.com,resources=supplies,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=fund.example.com,resources=supplies/status,verbs=get;update;patch

func (r *SupplyReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("supply", req.NamespacedName)

	var supply fundv1beta1.Supply
	if err := r.Get(ctx, req.NamespacedName, &supply); err != nil {
		log.Error(err, fmt.Sprintf("unable to fetch Supply: %v", req.NamespacedName))
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, ignoreNotFound(err)
	}

	if supply.Status.Phase != "" {
		// This supply has been processed, never process it again.
		return ctrl.Result{}, nil
	}

	// List all FundPool.
	var fundPools fundv1beta1.FundPoolList
	if err := r.List(ctx, &fundPools); err != nil {
		log.Error(err, "unable to list FundPool")
		return ctrl.Result{}, err
	}

	// If no FundPool, failed to satisfy supply.
	if len(fundPools.Items) == 0 {
		supply.Status.Phase = fundv1beta1.Failed
		if err := r.Update(ctx, &supply); err != nil {
			log.Error(err, fmt.Sprintf("failed to update Supply: %v", req.NamespacedName))
			return ctrl.Result{}, err
		}
	}

	m := map[string]*fundv1beta1.FundPool{}
	var pools types.Pools
	for i, _ := range fundPools.Items {
		pool := fundPools.Items[i].DeepCopy()
		name := pool.Name
		pools = append(pools, types.Pool{
			Name:  name,
			Value: *pool.Spec.Balance,
		})
		m[name] = pool
	}
	f := *supply.Spec.Request

	result, err := variance.MinVariance(pools, f)
	if err != nil {
		log.Error(err, "failed to get minimum variance")
		return ctrl.Result{}, err
	}

	// Update Fundpools.
	allocations := []fundv1beta1.Allocation{}
	for name, value := range result {
		pool := m[name]
		*pool.Spec.Balance -= value
		if err := r.Update(ctx, pool); err != nil {
			// TODO: the update of FundPools should be atomic operation.
			// If one update failed, all previous update should rollback.
			log.Error(err, fmt.Sprintf("unable to update FundPool: %v", name))
		}

		value := value
		allocations = append(allocations, fundv1beta1.Allocation{
			Pool:       name,
			Shortfalls: &value,
		})
	}

	// Update Supply.
	supply.Status.Phase = fundv1beta1.Succeeded
	supply.Status.Allocations = allocations
	if err := r.Update(ctx, &supply); err != nil {
		// TODO: if update failed, all previous operations should rollback.
		log.Error(err, fmt.Sprintf("failed to update Supply: %v", req.NamespacedName))
	}

	return ctrl.Result{}, nil
}

func (r *SupplyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&fundv1beta1.Supply{}).
		Complete(r)
}

/* Copyright 2025. McKinsey & Company */

package controller

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	arkv1alpha1 "mckinsey.com/ark/api/v1alpha1"
	"mckinsey.com/ark/internal/common"
)

// EvaluatorReconciler reconciles an Evaluator object
type EvaluatorReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
	resolver *common.ValueSourceResolver
}

// +kubebuilder:rbac:groups=ark.mckinsey.com,resources=evaluators,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=ark.mckinsey.com,resources=evaluators/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=ark.mckinsey.com,resources=evaluators/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch

func (r *EvaluatorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	var evaluator arkv1alpha1.Evaluator
	if err := r.Get(ctx, req.NamespacedName, &evaluator); err != nil {
		if errors.IsNotFound(err) {
			log.Info("Evaluator deleted", "evaluator", req.Name)
			return ctrl.Result{}, nil
		}
		log.Error(err, "unable to fetch Evaluator")
		return ctrl.Result{}, err
	}

	// State machine approach following Memory pattern
	switch evaluator.Status.Phase {
	case statusReady, statusError:
		// Terminal states - no further processing needed
		return ctrl.Result{}, nil
	case statusRunning:
		// Continue processing
		return r.processEvaluator(ctx, evaluator)
	default:
		// Initialize to running state
		if err := r.updateStatus(ctx, evaluator, statusRunning, "Resolving evaluator address"); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}
}

func (r *EvaluatorReconciler) getResolver() *common.ValueSourceResolver {
	if r.resolver == nil {
		r.resolver = common.NewValueSourceResolver(r.Client)
	}
	return r.resolver
}

func (r *EvaluatorReconciler) processEvaluator(ctx context.Context, evaluator arkv1alpha1.Evaluator) (ctrl.Result, error) {
	log := logf.FromContext(ctx)
	log.Info("Processing evaluator", "evaluator", evaluator.Name)

	resolver := r.getResolver()
	resolvedAddress, err := resolver.ResolveValueSource(ctx, evaluator.Spec.Address, evaluator.Namespace)
	if err != nil {
		log.Error(err, "failed to resolve Evaluator address", "evaluator", evaluator.Name)
		if err := r.updateStatus(ctx, evaluator, statusError, fmt.Sprintf("Failed to resolve address: %v", err)); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	// Update resolved address in status
	evaluator.Status.LastResolvedAddress = resolvedAddress

	// Mark as ready
	if err := r.updateStatus(ctx, evaluator, statusReady, "Evaluator address resolved successfully"); err != nil {
		return ctrl.Result{}, err
	}

	log.Info("Evaluator processed successfully", "evaluator", evaluator.Name, "resolvedAddress", resolvedAddress)
	return ctrl.Result{}, nil
}

func (r *EvaluatorReconciler) updateStatus(ctx context.Context, evaluator arkv1alpha1.Evaluator, phase, message string) error {
	log := logf.FromContext(ctx)

	evaluator.Status.Phase = phase
	evaluator.Status.Message = message

	if err := r.Status().Update(ctx, &evaluator); err != nil {
		log.Error(err, "failed to update Evaluator status", "evaluator", evaluator.Name)
		return err
	}

	log.Info("Updated Evaluator status", "evaluator", evaluator.Name, "phase", phase, "message", message)
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *EvaluatorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&arkv1alpha1.Evaluator{}).
		Named("evaluator").
		Complete(r)
}

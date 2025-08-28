/* Copyright 2025. McKinsey & Company */

package v1

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	arkv1alpha1 "mckinsey.com/ark/api/v1alpha1"
	"mckinsey.com/ark/internal/common"
)

var mcpserverlog = logf.Log.WithName("mcpserver-resource")

func SetupMCPServerWebhookWithManager(mgr ctrl.Manager) error {
	k8sClient := mgr.GetClient()
	return ctrl.NewWebhookManagedBy(mgr).
		For(&arkv1alpha1.MCPServer{}).
		WithValidator(&MCPServerValidator{
			Client:   k8sClient,
			Resolver: common.NewValueSourceResolver(k8sClient),
		}).
		Complete()
}

// +kubebuilder:webhook:path=/validate-ark-mckinsey-com-v1alpha1-mcpserver,mutating=false,failurePolicy=fail,sideEffects=None,groups=ark.mckinsey.com,resources=mcpserver,verbs=create;update,versions=v1alpha1,name=vmcpserver-v1.kb.io,admissionReviewVersions=v1

type MCPServerValidator struct {
	Client   client.Client
	Resolver *common.ValueSourceResolver
}

var _ webhook.CustomValidator = &MCPServerValidator{}

func (v *MCPServerValidator) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	mcpserver, ok := obj.(*arkv1alpha1.MCPServer)
	if !ok {
		return nil, fmt.Errorf("expected a MCPServer object but got %T", obj)
	}

	mcpserverlog.Info("Validating MCPServer", "name", mcpserver.GetName(), "namespace", mcpserver.GetNamespace())

	_, err := v.Resolver.ResolveValueSource(ctx, mcpserver.Spec.Address, mcpserver.GetNamespace())
	if err != nil {
		mcpserverlog.Error(err, "Failed to resolve Address", "mcpserver", mcpserver.GetName())
		return nil, fmt.Errorf("failed to resolve Address: %w", err)
	}

	for i, header := range mcpserver.Spec.Headers {
		if err := v.validateHeaderValue(ctx, header.Value, mcpserver.GetNamespace()); err != nil {
			mcpserverlog.Error(err, "Failed to validate header value", "mcpserver", mcpserver.GetName(), "header", header.Name)
			return nil, fmt.Errorf("failed to validate header %s (index %d): %w", header.Name, i, err)
		}
	}

	// Validate PollInterval
	if err := ValidatePollInterval(mcpserver.Spec.PollInterval.Duration); err != nil {
		mcpserverlog.Error(err, "Failed to validate pollInterval", "mcpserver", mcpserver.GetName())
		return nil, fmt.Errorf("failed to validate pollInterval: %w", err)
	}

	mcpserverlog.Info("MCPServer validation complete", "name", mcpserver.GetName())

	return nil, nil
}

func (v *MCPServerValidator) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	return v.ValidateCreate(ctx, newObj)
}

func (v *MCPServerValidator) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	return nil, nil
}

func (v *MCPServerValidator) validateHeaderValue(ctx context.Context, headerValue arkv1alpha1.HeaderValue, namespace string) error {
	if headerValue.Value != "" {
		return nil
	}

	if headerValue.ValueFrom == nil {
		return fmt.Errorf("header value must have either value or valueFrom specified")
	}

	if headerValue.ValueFrom.SecretKeyRef != nil {
		return v.validateSecretKeyRef(ctx, headerValue.ValueFrom.SecretKeyRef, namespace)
	}

	return fmt.Errorf("no valid valueFrom source specified for header")
}

func (v *MCPServerValidator) validateSecretKeyRef(ctx context.Context, secretRef *corev1.SecretKeySelector, namespace string) error {
	if secretRef.Name == "" {
		return fmt.Errorf("secret name is required")
	}

	secret := &corev1.Secret{}
	err := v.Client.Get(ctx, types.NamespacedName{
		Name:      secretRef.Name,
		Namespace: namespace,
	}, secret)
	if err != nil {
		return fmt.Errorf("failed to get secret %s/%s: %w", namespace, secretRef.Name, err)
	}

	if secretRef.Key != "" {
		if _, exists := secret.Data[secretRef.Key]; !exists {
			return fmt.Errorf("key %s not found in secret %s/%s", secretRef.Key, namespace, secretRef.Name)
		}
	}

	return nil
}

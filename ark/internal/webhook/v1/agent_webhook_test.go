/* Copyright 2025. McKinsey & Company */

package v1

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	arkv1alpha1 "mckinsey.com/ark/api/v1alpha1"
)

var _ = Describe("Agent Webhook", func() {
	var (
		ctx       context.Context
		agent     *arkv1alpha1.Agent
		validator *AgentCustomValidator
	)

	BeforeEach(func() {
		ctx = context.Background()

		// Setup scheme
		s := runtime.NewScheme()
		Expect(arkv1alpha1.AddToScheme(s)).To(Succeed())

		// Create fake client
		fakeClient := fake.NewClientBuilder().WithScheme(s).Build()

		// Create validator
		validator = &AgentCustomValidator{
			ResourceValidator: &ResourceValidator{Client: fakeClient},
		}

		// Create base agent
		agent = &arkv1alpha1.Agent{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-agent",
				Namespace: "default",
			},
			Spec: arkv1alpha1.AgentSpec{
				Description: "Test agent",
				Prompt:      "You are a test agent",
			},
		}
	})

	Context("When validating agent model requirements", func() {
		It("Should fail when no model is specified and no default model exists", func() {
			// Agent without modelRef and no default model in namespace
			warnings, err := validator.ValidateCreate(ctx, agent)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("no model specified for agent and no 'default' model found"))
			Expect(warnings).To(BeEmpty())
		})

		It("Should allow A2A agents without model validation", func() {
			// Set execution engine to A2A
			agent.Spec.ExecutionEngine = &arkv1alpha1.ExecutionEngineRef{
				Name: ExecutionEngineA2A,
			}

			warnings, err := validator.ValidateCreate(ctx, agent)
			Expect(err).NotTo(HaveOccurred())
			Expect(warnings).To(BeEmpty())
		})

		It("Should allow A2A agents to be updated without model validation", func() {
			// Set execution engine to A2A
			agent.Spec.ExecutionEngine = &arkv1alpha1.ExecutionEngineRef{
				Name: ExecutionEngineA2A,
			}

			oldAgent := agent.DeepCopy()
			agent.Spec.Description = "Updated A2A agent"

			warnings, err := validator.ValidateUpdate(ctx, oldAgent, agent)
			Expect(err).NotTo(HaveOccurred())
			Expect(warnings).To(BeEmpty())
		})

		It("Should still validate models for non-A2A agents", func() {
			// Set execution engine to something other than A2A
			agent.Spec.ExecutionEngine = &arkv1alpha1.ExecutionEngineRef{
				Name: "langchain",
			}

			warnings, err := validator.ValidateCreate(ctx, agent)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("no model specified for agent and no 'default' model found"))
			Expect(warnings).To(BeEmpty())
		})
	})
})

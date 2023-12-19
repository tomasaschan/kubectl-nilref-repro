package controller_test

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/client"

	addonsv1alpha1 "github.com/tomasaschan/kubectl-nilref-repro/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Reproduce nilref in kubectl", func() {
	ctx := context.Background()

	It("can apply cleanly", func() { // this test fails, with logs from the reconciler showing the nilref panic
		By("Creating the namespace")
		Expect(k8sClient.Create(ctx, &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: "testing",
			},
		})).Should(WithTransform(client.IgnoreAlreadyExists, Succeed()))

		By("Applying the resource")
		Expect(k8sClient.Create(ctx, &addonsv1alpha1.Repro{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "repro",
				Namespace: "testing",
			},
		})).Should(WithTransform(client.IgnoreAlreadyExists, Succeed()))

		Eventually(func() error {
			cm := corev1.ConfigMap{}
			key := client.ObjectKey{Name: "repro", Namespace: "default"}

			err := k8sClient.Get(ctx, key, &cm)
			if err != nil {
				return fmt.Errorf("get configmap: %w", err)
			}

			return nil
		}, "3s").Should(Succeed())
	})
})

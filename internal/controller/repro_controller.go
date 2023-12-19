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
	"fmt"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"sigs.k8s.io/kubebuilder-declarative-pattern/pkg/patterns/addon"
	"sigs.k8s.io/kubebuilder-declarative-pattern/pkg/patterns/addon/pkg/loaders"
	"sigs.k8s.io/kubebuilder-declarative-pattern/pkg/patterns/addon/pkg/status"
	"sigs.k8s.io/kubebuilder-declarative-pattern/pkg/patterns/declarative"

	addonsv1alpha1 "github.com/tomasaschan/kubectl-nilref-repro/api/v1alpha1"
)

var _ reconcile.Reconciler = &ReproReconciler{}

// ReproReconciler reconciles a Repro object
type ReproReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme

	declarative.Reconciler

	ChannelsPath string
}

//+kubebuilder:rbac:groups=addons,resources=reproes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=addons,resources=reproes/status,verbs=get;update;patch

// SetupWithManager sets up the controller with the Manager.
func (r *ReproReconciler) SetupWithManager(mgr ctrl.Manager) error {
	addon.Init()

	labels := map[string]string{
		"k8s-app": "repro",
	}

	watchLabels := declarative.SourceLabel(mgr.GetScheme())

	var loader declarative.ManifestLoaderFunc
	if r.ChannelsPath != "" {
		loader = func() (declarative.ManifestController, error) {
			return loaders.NewManifestLoader(r.ChannelsPath)
		}
	} else {
		loader = declarative.DefaultManifestLoader
	}
	mc, err := loader()
	if err != nil {
		return fmt.Errorf("creating manifest loader: %w", err)
	}

	if err := r.Reconciler.Init(mgr, &addonsv1alpha1.Repro{},
		declarative.WithManifestController(mc),
		declarative.WithObjectTransform(declarative.AddLabels(labels)),
		declarative.WithOwner(declarative.SourceAsOwner),
		declarative.WithLabels(watchLabels),
		declarative.WithStatus(status.NewBasic(mgr.GetClient())),
		declarative.WithObjectTransform(addon.ApplyPatches),
	); err != nil {
		return err
	}

	c, err := controller.New("repro-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to Repro
	err = c.Watch(source.Kind(mgr.GetCache(), &addonsv1alpha1.Repro{}), &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to deployed objects
	err = declarative.WatchAll(mgr.GetConfig(), c, r, watchLabels)
	if err != nil {
		return err
	}

	return nil
}

/*
Copyright 2018 The Kubernetes Authors.

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

package cloudproviders

import (
	"context"

	"github.com/carolynvs/cloudkinds/pkg/apis/cloudkinds/v1alpha1"
	"github.com/carolynvs/cloudkinds/pkg/controller/cloudkinds"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextensionclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// Add creates a new CloudProvider Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	crdClient, _ := apiextensionclient.NewForConfig(mgr.GetConfig())

	return &ReconcileCloudProvider{
		Client:    mgr.GetClient(),
		scheme:    mgr.GetScheme(),
		crdClient: crdClient,
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("CloudProviders-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &v1alpha1.CloudProvider{}}, &handler.EnqueueRequestForObject{})

	return nil
}

var _ reconcile.Reconciler = &ReconcileCloudProvider{}

// ReconcileCloudProvider reconciles a CloudResource object
type ReconcileCloudProvider struct {
	client.Client

	scheme    *runtime.Scheme
	crdClient apiextensionclient.Interface
}

// Reconcile creates CRDs for all kinds that our providers report that they work with
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=cloudkinds.k8s.io,resources=providers,verbs=get;list;watch
// +kubebuilder:rbac:groups=apiextensions.k8s.io,resources=customresourcedefinitions,verbs=get;list;watch;create
func (r *ReconcileCloudProvider) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	provider := &v1alpha1.CloudProvider{}
	err := r.Get(context.Background(), request.NamespacedName, provider)
	if err != nil {
		return reconcile.Result{Requeue: true}, err
	}

	// We should use the CRD informer cache instead
	registeredCRDs := make(map[string]bool)
	crds := &v1beta1.CustomResourceDefinitionList{}
	err = r.List(context.Background(), &client.ListOptions{}, crds)
	if err != nil {
		return reconcile.Result{Requeue: true}, err
	}

	for _, r := range crds.Items {
		if r.Spec.Group == v1alpha1.SchemeGroupVersion.Group {
			registeredCRDs[r.Spec.Names.Kind] = true
		}
	}

	for _, kind := range provider.Spec.Kinds {
		if _, ok := registeredCRDs[kind]; !ok {
			err = cloudkinds.RegisterCloudKind(r.crdClient, kind)
			if err != nil {
				return reconcile.Result{Requeue: true}, err
			}
		}
	}

	return reconcile.Result{}, nil
}

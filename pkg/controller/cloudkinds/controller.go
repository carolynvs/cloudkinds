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

package cloudkinds

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/carolynvs/cloudkinds/pkg/apis/cloudkinds/v1alpha1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// Add creates a new CloudKind Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileCloudKind{Client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("cloudkinds-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to any registered cloudkind
	kinds := []schema.GroupVersionKind{
		v1alpha1.SchemeGroupVersion.WithKind("CloudResource"),
	}
	for _, kind := range kinds {
		cloudKind := NewCloudKind(kind)
		err = c.Watch(&source.Kind{Type: cloudKind}, &handler.EnqueueRequestForObject{})
		if err != nil {
			return err
		}
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileCloudKind{}

// ReconcileCloudKind reconciles a CloudResource object
type ReconcileCloudKind struct {
	client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a CloudResource object and makes changes based on the state read
// and what is in the CloudResource.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  The scaffolding writes
// a Deployment as an example
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=cloudkinds.k8s.io,resources=cloudresources,verbs=get;list;watch;create;update;patch;delete
func (r *ReconcileCloudKind) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	fmt.Printf("farts are funny: %#v\n", request)

	obj := NewCloudKind(request.GroupVersionKind)
	err := r.Get(context.Background(), request.NamespacedName, obj)

	fmt.Printf("%#v\n", obj)
	return reconcile.Result{}, err
}

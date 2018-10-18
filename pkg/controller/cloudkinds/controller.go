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
// +kubebuilder:rbac:groups=cloudkinds.k8s.io,resources=providers,verbs=get;list;watch
// +kubebuilder:rbac:groups=cloudkinds.k8s.io,resources=cloudresources,verbs=get;list;watch;update;patch
func (r *ReconcileCloudKind) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	fmt.Println("farts are funny")
	// Resolve a provider for this kind
	providers := &v1alpha1.ProviderList{}
	err := r.List(context.Background(), &client.ListOptions{Namespace: request.NamespacedName.Namespace}, providers)

	var provider *v1alpha1.Provider
	for _, p := range providers.Items {
		for _, k := range p.Spec.Kinds {
			if k == request.Kind {
				provider = &p
			}
		}
	}

	if provider == nil {
		// We can't reconcile this _yet_ because there isn't a provider registered  - requeue the request with a bit of a buffer.
		err = fmt.Errorf("no provider registered for kind: %s", request.Kind)
		fmt.Println(err)
		return reconcile.Result{Requeue: true, RequeueAfter: 30 * time.Second}, err
	}

	obj := NewCloudKind(request.GroupVersionKind)
	err = r.Get(context.Background(), request.NamespacedName, obj)
	if err != nil {
		fmt.Println(err)

		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		// Error reading the object - requeue the request.
		return reconcile.Result{Requeue: true}, err
	}

	apiVersion, kind := request.ToAPIVersionAndKind()
	evt := v1alpha1.ResourceEvent{
		Action: v1alpha1.ResourceCreated, // TODO: Base the action on the status of the resource
		Resource: v1alpha1.ResourceReference{
			TypeMeta:  v1.TypeMeta{APIVersion: apiVersion, Kind: kind},
			Namespace: request.NamespacedName.Namespace,
			Name:      request.NamespacedName.Name,
		},
	}
	bodyJson, err := json.Marshal(evt)
	if err != nil {
		fmt.Println(err)
		return reconcile.Result{Requeue: true}, err
	}
	body := bytes.NewReader(bodyJson)

	response, err := http.DefaultClient.Post(provider.Spec.WebHook, "application/json", body)
	if err != nil {
		fmt.Println(err)
		return reconcile.Result{}, err
	}

	defer response.Body.Close()
	result, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		result = []byte("could not read response body")
	}

	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("provider failed: %v %v %s", response.StatusCode, provider.Spec.WebHook, string(result))
		fmt.Println(err)
		return reconcile.Result{}, err
	} else {
		fmt.Printf("%#v\n", obj)
		fmt.Println(string(result))
	}

	return reconcile.Result{}, nil
}

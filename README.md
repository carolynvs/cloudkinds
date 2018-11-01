# Cloud Kinds

This is a prototype for exploring a better user experience for using cloud services,
such as a mysql as a service from your cloud provider or a mega dev db of doom from
your IT, from inside kubernetes.

It uses a kind per service type, so that you can work with them without specifying
exactly which cloud should handle creating it:

```yaml
apiVersion: cloudkinds.k8s.io/v1alpha1
kind: MySQL
metadata:
  name: mydb
spec:
  requireSSL: true
```

The cluster operator is responsible for configuring backend providers and environment defaults
to handle creating services to support these new kinds. Examples of a backend provider are:
Service Catalog, AWS Operator, etcd operator.

A major goal is to support using the above resource and have it be portable across
clouds, so that someone could use the exact same definition and have it work with
different providers:

1. A developer using minikube and minibroker could work against an in-cluster mysql database.
1. When the app is deployed in the on-premise test cluster, it uses the shared test database server.
1. In production, it could use Amazon RDS, or Azure MySQL.

# Try it out

```console
# Run the cloud kinds service on the current cluster
$ TAG=alpha make deploy

# Tail the logs for a sample cloud kinds provider that just echos back webhook payloads
$ kubectl logs deploy/cloudkinds-sampleprovider -f&
I am a sample CloudKinds provider that does absolutely nothing useful! ☁️🌈
Listening on *:8080

# Create a new cloudkinds resource
$ kubectl apply -f hack/samples/sample.yaml
cloudresource.cloudkinds.k8s.io/sample created
Handled /
	{"action":"create","resource":{"kind":"CloudResource","apiVersion":"cloudkinds.k8s.io/v1alpha1","name":"sample","namespace":"cloudkinds"}}
```

This isn't terribly exciting (yet), but it shows that we can watch for arbitrary CRDs and call the provider's webhook.

# Implementation Notes

* Uses kubebuilder and takes advantage of the untyped client so that we can watch for any resource dynamically.
* The set of resources to watch for is configurable. When a new provider that supports more kinds is registered,
  the controller should be able to pick that up and start reacting to those new resources when they are created.
* Backend providers shouldn't require changes. Service Catalog or AWS operator continue to do their thing.
* Backend providers (or interested parties) can register an adapter that our controller can call out to and handle translating
  one of these resources into the appropriate request to the backend provider.
* Operators will be able to set default values that should be applied to certain kinds, so that they can say 
  "all mysql database should require ssl", or "whitelist these ips by default", so that each dev doesn't have to 
  re-specify this stuff. It also enables using per cluster/environment config without having to define it on the 
  resource itself. (This is from the [Service Plan Defaults proposal](https://github.com/carolynvs/service-catalog/blob/default-service-plan-proposal/docs/proposals/default-service-plans.md)
  that is being implemented in service catalog).

## Install Workflow
1. The service provider is installed. This could be a broker + service catalog, or could be an operator 
   like etcd operator or aws operator.
1. CloudKinds is installed.
1. Register adapters that convert from a custom kind, like "mysql", to a backend service provider.
   * An adapter for Service Catalog would convert resources into service instance as long as it finds a matching class/plan.
   * An adapter for aws operator would convert the resource into the resources understood by that operator.
1. CloudKinds manages registering the CRDs (such as `kind: MySQL`), with Kubernetes, not the adapters.

## Create Cloud Resource Workflow
1. User creates a cloud resource
    ```yaml
    apiVersion: cloudkinds.k8s.io/v1alpha1
    kind: MySQL
    metadata:
      name: carolyn-dev-db
   ```
1. We apply any relevant selectors and labels on the CRD to aid in picking a backend service provider.
1. We call the adapter's webhook and give it the CRD. It does it’s thing (which may just be to turn around and 
    make a service instance, or the more specific CRD for an operator).
1. The adapter's webhook handles interacting with the backend service provider and then, 
   setting owner references, extra labels, etc, and finally responding with OK.
   * Later we may use the response body to set more interesting information on the status sub-resource.
1. We set a finalizer on the resource so that we can manage its deletion later.
1. We also set on the status which provider handled the resource.

## Modify Cloud Resource Workflow
1. User edits a cloud resource.
1. We call the adapter and pass along a reference to the modified resource.
1. The adapter's webhook handles updating the resource with the backend.

## Delete Cloud Resource Workflow
1. User deletes the cloud resource
1. Our finalizer handles calling the adapter's webhook to pass along the delete request.

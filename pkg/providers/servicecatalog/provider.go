package servicecatalog

import (
	"encoding/json"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/pkg/errors"

	"github.com/carolynvs/cloudkinds/pkg/apis/cloudkinds/v1alpha1"
	svcatclient "github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	"github.com/kubernetes-incubator/service-catalog/pkg/util/kube"
	k8sclient "k8s.io/client-go/kubernetes"
)

func DealWithIt(payload []byte) ([]byte, error) {
	// parse the payload
	evt := &v1alpha1.ResourceEvent{}
	err := json.Unmarshal(payload, evt)
	if err != nil {
		return nil, err
	}

	// load up the current cluster config
	config := kube.GetConfig("", "")
	restConfig, err := config.ClientConfig()
	if err != nil {
		return nil, errors.New("could not get Kubernetes config")
	}
	namespace, _, err := config.Namespace()
	k8sClient, err := k8sclient.NewForConfig(restConfig)
	if err != nil {
		return nil, nil, "", err
	}
	svcatClient, err := svcatclient.NewForConfig(restConfig)

	// find the corresponding service instance
	svcatClient.ServicecatalogV1beta1().ServiceInstances(evt.Resource.Namespace).Get(evt.Resource.Name, metav1.GetOptions{})
	// create a service instance
	// mark the instance as owned by the crd
	return []byte(`{"msg": "farts are funny"}`), nil
}

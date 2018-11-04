package cloudkinds

import (
	"fmt"
	"strings"

	"github.com/carolynvs/cloudkinds/pkg/apis/cloudkinds/v1alpha1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextensionclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// NewCloudKind creates an unstructured k8s resource that can store anything,
// and does not need to be statically defined at compile time.
// Note: Due to checks for the actual struct, instead of the unstructured interface,
// We can't use an embedded field on a custom type.
func NewCloudKind(kind schema.GroupVersionKind) *unstructured.Unstructured {
	ret := &unstructured.Unstructured{}
	ret.SetGroupVersionKind(kind)
	return ret
}

func RegisterCloudKind(cl apiextensionclient.Interface, kind string) error {
	fmt.Printf("Discovered CloudKind: %s\n", kind)

	plural := strings.ToLower(kind) + "s"
	group := v1alpha1.SchemeGroupVersion.Group
	version := v1alpha1.SchemeGroupVersion.Version
	name := fmt.Sprintf("%s.%s", plural, group)
	crd := &apiextensionsv1beta1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{Name: name},
		Spec: apiextensionsv1beta1.CustomResourceDefinitionSpec{
			Group:   group,
			Version: version,
			Scope:   apiextensionsv1beta1.NamespaceScoped,
			Names: apiextensionsv1beta1.CustomResourceDefinitionNames{
				Plural: plural,
				Kind:   kind,
			},
		},
	}

	_, err := cl.ApiextensionsV1beta1().CustomResourceDefinitions().Create(crd)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		return err
	}
	fmt.Printf("Created CRD %s\n", name)
	return nil
}

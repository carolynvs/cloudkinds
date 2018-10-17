package cloudkinds

import (
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

package reconciler

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

// addFinalizer accepts a metav1 object and adds the provided finalizer if not present.
func addFinalizer(o metav1.Object, finalizer string) {
	f := o.GetFinalizers()
	for _, e := range f {
		if e == finalizer {
			return
		}
	}
	o.SetFinalizers(append(f, finalizer))
}

// removeFinalizer accepts a metav1 object and removes the provided finalizer if present.
func removeFinalizer(o metav1.Object, finalizer string) {
	f := o.GetFinalizers()
	for i, e := range f {
		if e == finalizer {
			f = append(f[:i], f[i+1:]...)
			o.SetFinalizers(f)
			return
		}
	}
}

// hasFinalizer accepts a metav1 object and returns true if the the object has the provided finalizer.
func hasFinalizer(o metav1.Object, finalizer string) bool {
	return containsString(o.GetFinalizers(), finalizer)
}

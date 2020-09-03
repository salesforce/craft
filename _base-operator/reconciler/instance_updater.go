// Copyright (c) 2020, salesforce.com, inc.
// All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// For full license text, see LICENSE.txt file in the repo root or https://opensource.org/licenses/BSD-3-Clause

package reconciler

import (
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	corev1 "k8s.io/api/core/v1"
	api "{{ .Repo }}/api/{{ .Version }}"
)

// modifies the runtime.Object in place
type statusUpdate = func(status *Status)
type metaUpdate = func(meta metav1.Object)

// instanceUpdater is a mechanism to enable updating the shared sections of the manifest
// Typically the Status section and the metadata.
type instanceUpdater struct {
	StatusUpdater
	metaUpdates   []metaUpdate
	statusUpdates []statusUpdate
}

func (updater *instanceUpdater) addFinalizer(name string) {
	updateFunc := func(meta metav1.Object) { addFinalizer(meta, name) }
	updater.metaUpdates = append(updater.metaUpdates, updateFunc)
}

func (updater *instanceUpdater) removeFinalizer(name string) {
	updateFunc := func(meta metav1.Object) { removeFinalizer(meta, name) }
	updater.metaUpdates = append(updater.metaUpdates, updateFunc)
}

func (updater *instanceUpdater) setStatusPayload(statusPayload interface{}) {
	updateFunc := func(s *Status) {
		s.StatusPayload = statusPayload
	}
	updater.statusUpdates = append(updater.statusUpdates, updateFunc)
}

func (updater *instanceUpdater) setTerminated(terminated *corev1.ContainerStateTerminated) {
	updateFunc := func(s *Status) {
		s.Terminated = terminated
	}
	updater.statusUpdates = append(updater.statusUpdates, updateFunc)
}

func (updater *instanceUpdater) setPodConfig(pod api.Pod) {
	updateFunc := func(s *Status) {
		s.Pod = pod
	}
	updater.statusUpdates = append(updater.statusUpdates, updateFunc)
}

func (updater *instanceUpdater) setReconcileState(state ReconcileState, message string) {
	updateFunc := func(s *Status) {
		s.State = state
		s.Message = message
	}
	updater.statusUpdates = append(updater.statusUpdates, updateFunc)
}

func (updater *instanceUpdater) setAnnotation(name string, value string) {
	updateFunc := func(meta metav1.Object) {
		annotations := meta.GetAnnotations()
		if annotations == nil {
			annotations = map[string]string{}
		}
		annotations[name] = value
		meta.SetAnnotations(annotations)
	}
	updater.metaUpdates = append(updater.metaUpdates, updateFunc)
}

func (updater *instanceUpdater) setOwnerReferences(owners []runtime.Object) {
	updateFunc := func(s metav1.Object) {
		references := make([]metav1.OwnerReference, len(owners))
		for i, o := range owners {
			controller := true
			meta, _ := apimeta.Accessor(o)
			references[i] = metav1.OwnerReference{
				APIVersion: "v1",
				Kind:       o.GetObjectKind().GroupVersionKind().Kind,
				Name:       meta.GetName(),
				UID:        meta.GetUID(),
				Controller: &controller,
			}
		}
		s.SetOwnerReferences(references)
	}
	updater.metaUpdates = append(updater.metaUpdates, updateFunc)
}

func (updater *instanceUpdater) applyUpdates(instance runtime.Object, status *Status) error {
	for _, f := range updater.statusUpdates {
		f(status)
	}
	err := updater.StatusUpdater(instance, status)
	m, _ := apimeta.Accessor(instance)
	for _, f := range updater.metaUpdates {
		f(m)
	}
	return err
}

func (updater *instanceUpdater) clear() {
	updater.metaUpdates = []metaUpdate{}
	updater.statusUpdates = []statusUpdate{}
}

func (updater *instanceUpdater) hasUpdates() bool {
	return len(updater.metaUpdates) > 0 || len(updater.statusUpdates) > 0
}

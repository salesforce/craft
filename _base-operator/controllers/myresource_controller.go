// Copyright (c) 2020, salesforce.com, inc.
// All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// For full license text, see LICENSE.txt file in the repo root or https://opensource.org/licenses/BSD-3-Clause

/*
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

package controllers

import (
	api "{{ .Repo }}/api/{{  .Version }}"
	"{{ .Repo }}/reconciler"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/go-logr/logr"
)

type ControllerFactory struct {
	ResourceManagerCreator func(logr.Logger, record.EventRecorder) ResourceManager
	Scheme                 *runtime.Scheme
}

// +kubebuilder:rbac:groups="",resources=pods;events,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=extensions;apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups={{  .Group }}.{{  .Domain }},resources={{  .LowerRes }}s,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups={{  .Group }}.{{  .Domain }},resources={{  .LowerRes }}s/status,verbs=get;update;patch

// +kubebuilder:printcolumn:name="State",type=string,JSONPath=`.status.state`
// +kubebuilder:printcolumn:name="Message",type=string,JSONPath=`.status.message`

const ResourceKind = "{{  .Resource }}"
const FinalizerName = "{{  .LowerRes }}s.{{  .Group }}.{{  .Domain }}"
const AnnotationBaseName = "{{  .Group }}.{{  .Domain }}"

func (factory *ControllerFactory) SetupWithManager(mgr ctrl.Manager, parameters reconciler.ReconcileParameters, log *logr.Logger) error {
	if log == nil {
		l := ctrl.Log.WithName("controllers")
		log = &l
	}
	gc, err := factory.createGenericController(mgr.GetClient(),
		(*log).WithName(ResourceKind),
		mgr.GetEventRecorderFor(ResourceKind+"-controller"), parameters)
	if err != nil {
		return err
	}
	// https://stuartleeks.com/posts/kubebuilder-event-filters-part-1-delete/
	// https://godoc.org/sigs.k8s.io/controller-runtime/pkg/event
	// https://github.com/kubernetes-sigs/kubebuilder/issues/618
	// https://godoc.org/sigs.k8s.io/controller-runtime/pkg/predicate#Predicate
	// https://book-v1.book.kubebuilder.io/beyond_basics/controller_watches.html
	return ctrl.NewControllerManagedBy(mgr).
		For(&api.{{  .Resource }}{}).
		WithEventFilter(predicate.Funcs{
			UpdateFunc: func(e event.UpdateEvent) bool {
				isSame := true
				for k, v := range e.ObjectOld.(*api.{{  .Resource }}).ObjectMeta.Annotations {
					isSame = isSame && (v == e.ObjectNew.(*api.{{  .Resource }}).ObjectMeta.Annotations[k])
				}
				if e.ObjectNew.(*api.{{  .Resource }}).ObjectMeta.DeletionTimestamp != nil {
					return true
				}
				return !isSame
			},
		}).
		Complete(gc)
}

func (factory *ControllerFactory) createGenericController(kubeClient client.Client, logger logr.Logger, recorder record.EventRecorder, parameters reconciler.ReconcileParameters) (*reconciler.GenericController, error) {
	resourceManagerClient := factory.ResourceManagerCreator(logger, recorder)

	return reconciler.CreateGenericController(parameters, ResourceKind, kubeClient, logger, recorder, factory.Scheme, &resourceManagerClient, &definitionManager{}, FinalizerName, AnnotationBaseName, nil)
}

func CreateResourceManager(logger logr.Logger, recorder record.EventRecorder) ResourceManager {
	return ResourceManager{
		Logger:     logger,
		Recorder:   recorder,
	}
}

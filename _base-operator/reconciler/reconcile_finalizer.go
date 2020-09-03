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

package reconciler

import (
	"context"
	"fmt"

	"github.com/prometheus/common/log"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *reconcileFinalizer) isDefined() bool {
	return hasFinalizer(r.objectMeta, r.FinalizerName)
}

func (r *reconcileFinalizer) add(ctx context.Context) (ctrl.Result, error) {
	updater := r.instanceUpdater

	updater.addFinalizer(r.FinalizerName)
	r.log.Info("Adding finalizer to resource")
	return r.applyTransition(ctx, "Finalizer", Pending, nil)
}
func (r *ReconcileRunner) terminationCheck(ctx context.Context) (string, string) {

	instance, err := convertInstance(r.instance)

	pod := instance.Status.Pod
	found := &corev1.Pod{}

	err = r.KubeClient.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, found)

	if err != nil && apierrors.IsNotFound(err) {
		return pod.Type, "Checking"
	}
	if len(found.Status.ContainerStatuses) == 0 {
		return pod.Type, "Checking"
	}
	terminated := found.Status.ContainerStatuses[0].State.Terminated

	if terminated != nil {
		// put the exitcode later
		r.instanceUpdater.setTerminated(terminated)
		r.KubeClient.Delete(ctx, found)
		return pod.Type, "PodDeleting"
		// return pod.Type, statusCodeMap[terminated.ExitCode]
	}
	return pod.Type, "Checking"
}

func (r *ReconcileRunner) deleteCheck(ctx context.Context) (string, string) {

	instance, err := convertInstance(r.instance)

	pod := instance.Status.Pod
	found := &corev1.Pod{}

	err = r.KubeClient.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, found)

	if err != nil && apierrors.IsNotFound(err) {
		terminated := instance.Status.Terminated
		return pod.Type, statusCodeMap[terminated.ExitCode]
	}
	return pod.Type, "PodDeleting"

}

func (r *reconcileFinalizer) handle() (ctrl.Result, error) {
	instance := r.instance
	updater := r.instanceUpdater
	ctx := context.Background()
	removeFinalizer := false
	requeue := false

	isTerminating := r.status.IsTerminating()
	r.log.Info(fmt.Sprintf("Finalizer: %s", r.status))
	if r.isDefined() {
		// Even before we cal ResourceManager.Delete, we verify the state of the resource
		// If it has not been created, we don't need to delete anything.
		var verifyResult VerifyResult
		var deleteResult DeleteResult
		var delete, checkState string

		if r.status.IsPodDeleting() {
			delete, checkState = r.deleteCheck(ctx)
		} else if r.status.IsChecking() {
			delete, checkState = r.terminationCheck(ctx)
		}

		if checkState == "PodDeleting" {
			return r.applyTransition(ctx, "deleteCheck", PodDeleting, nil)
		} else if checkState == "Checking" {
			return r.applyTransition(ctx, "terminationCheck", Checking, nil)
		} else if delete == "delete" {
			deleteResult := DeleteResult(checkState)
			r.log.Info(fmt.Sprintf("Finalizer deleteresult: %s", deleteResult))
		} else {
			// result is verify
			verifyResult = VerifyResult(checkState)
			r.log.Info(fmt.Sprintf("Finalizer verify result: %s", verifyResult))
		}

		if delete == "delete" {
			deleteResult = DeleteResult(checkState)
			if deleteResult.error() {
				log.Info("An error occurred attempting to delete managed object in finalizer. Cannot confirm that managed object has been deleted. Continuing deletion of kubernetes object anyway.")
				removeFinalizer = true
			} else if deleteResult.alreadyDeleted() || deleteResult.succeeded() {
				removeFinalizer = true
			} else if deleteResult.awaitingVerification() {
				requeue = true
			} else {
				// assert no more cases
				removeFinalizer = true
			}
		} else if verifyResult.missing() {
			removeFinalizer = true
		} else if verifyResult.deleting() {
			requeue = true
		} else if !isTerminating { // and one of verifyResult.ready() || verifyResult.recreateRequired() || verifyResult.updateRequired() || verifyResult.error()
			if verifyResult.error() {
				log.Info("An error occurred verifying state of managed object in finalizer. Cannot confirm that managed object can be deleted. Continuing deletion of kubernetes object anyway.")
				// TODO: maybe should rather retry a certain number of times before failing
			}
			permissions := r.getAccessPermissions()
			if !permissions.delete() {
				// if delete permission is turned off, just finalize, but don't delete
				r.log.Info("Resource is not managed by operator, bypassing delete of resource")
				removeFinalizer = true
			} else {
				// This block of code should only ever get called once.
				r.log.Info("Deleting resource")
				// delete state will be saved in the status
				podSpec, _ := r.ResourceManager.Delete(ctx)
				podValue, _ := r.ReconcileRunner.spawnPod("delete", podSpec)
				r.ReconcileRunner.instanceUpdater.setPodConfig(podValue)
				return r.applyTransition(ctx, "terminationCheck", Checking, nil)
			}
		} else {
			// this should never be called, as the first time r.ResourceManager.Delete is called isTerminating should be false
			// this implies that r.ResourceManager.Delete didn't throw an error, but didn't do anything either
			removeFinalizer = true
		}
	}

	if !isTerminating {
		updater.setReconcileState(Terminating, "")
	}
	if removeFinalizer {
		updater.removeFinalizer(r.FinalizerName)
	}

	requeueAfter := r.getRequeueAfter(Terminating)
	if removeFinalizer || !isTerminating {
		if err := r.updateInstance(ctx); err != nil {
			// if we can't update we have to requeue and hopefully it will remove the finalizer next time
			return ctrl.Result{Requeue: true, RequeueAfter: requeueAfter}, fmt.Errorf("Error removing finalizer: %v", err)
		}
		if !isTerminating {
			r.Recorder.Event(instance, corev1.EventTypeNormal, "Finalizer", "Setting state to terminating for "+r.Name)
		}
		if removeFinalizer {
			r.Recorder.Event(instance, corev1.EventTypeNormal, "Finalizer", "Removing finalizer for "+r.Name)
		}
	}

	if requeue {
		return ctrl.Result{Requeue: true, RequeueAfter: requeueAfter}, nil
	} else {
		r.Recorder.Event(instance, corev1.EventTypeNormal, "Finalizer", r.Name+" finalizer complete")
		return ctrl.Result{}, nil
	}
}

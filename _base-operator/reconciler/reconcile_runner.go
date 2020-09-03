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
	"encoding/json"
	"io/ioutil"
	"os"

	// "os/exec"
	"fmt"
	"strings"
	"time"

	"errors"
	{{ .Resource }}{{ .Version }} "{{ .Repo }}/api/{{ .Version }}"

	"github.com/hashicorp/vault/api"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/go-logr/logr"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	ctrl "sigs.k8s.io/controller-runtime"
)

var vault *api.Logical
var token = os.Getenv("VAULT_TOKEN")
var vault_addr = os.Getenv("VAULT_ADDR")

const LastAppliedAnnotation = "/last-applied-spec"

// Contains all the state involved in running a single reconcile event in the reconcile loo[
type ReconcileRunner struct {
	*GenericController
	*ResourceDefinition
	*DependencyDefinitions
	types.NamespacedName
	instance        runtime.Object
	objectMeta      metav1.Object
	status          *Status
	req             ctrl.Request
	log             logr.Logger
	instanceUpdater *instanceUpdater
	owner           runtime.Object
	dependencies    map[types.NamespacedName]runtime.Object
}

type reconcileFinalizer struct {
	ReconcileRunner
}

//runs a single reconcile on the
func (r *ReconcileRunner) run(ctx context.Context) (ctrl.Result, error) {

	// Verify that all dependencies are present in the cluster, and they are
	owner := r.Owner
	var allDeps []*Dependency
	if owner != nil {
		allDeps = append([]*Dependency{owner}, r.Dependencies...)
	} else {
		allDeps = r.Dependencies
	}
	status := r.status
	r.dependencies = map[types.NamespacedName]runtime.Object{}

	// jump out and requeue if any of the dependencies are missing
	for i, dep := range allDeps {
		instance := dep.InitialInstance
		err := r.KubeClient.Get(ctx, dep.NamespacedName, instance)
		log := r.log.WithValues("Dependency", dep.NamespacedName)

		// if any of the dependencies are not found, we jump out.
		if err != nil { // note that dependencies should be an empty array
			if apierrors.IsNotFound(err) {
				log.Info("Dependency not found for " + dep.NamespacedName.Name + ". Requeuing request.")
			} else {
				log.Info(fmt.Sprintf("Unable to retrieve dependency for %s: %v", dep.NamespacedName.Name, err.Error()))
			}
			return r.applyTransition(ctx, "Dependency", Pending, client.IgnoreNotFound(err))
		}

		// set the owner reference if owner is present and references have not been set
		// currently we only have single object ownership, but it is poosible to have multiple owners
		if owner != nil && i == 0 {
			if len(r.objectMeta.GetOwnerReferences()) == 0 {
				return r.setOwner(ctx, instance)
			}
			r.owner = instance
		}
		r.dependencies[dep.NamespacedName] = instance

		succeeded, err := dep.SucceededAccessor(instance)
		if err != nil {
			log.Info(fmt.Sprintf("Cannot get success state for %s. terminal failure.", dep.NamespacedName.Name))
			// Fail if cannot get Status accessor for dependency
			return r.applyTransition(ctx, "Dependency", Failed, err)
		}

		if !succeeded {
			log.Info("One of the dependencies is not in 'Succeeded' state, requeuing")
			return r.applyTransition(ctx, "Dependency", Pending, nil)
		}
	}
	// status = &Status{State: Checking}
	r.log.Info(fmt.Sprintf("ReconcileState: %s", status))
	// ****  checking for termination of dockerfile
	if status.IsChecking() {
		return r.check(ctx)
	}

	// ****  podPreviousPod for termination of dockerfile
	if status.IsPodDeleting() {
		return r.podDelete(ctx)
	}

	if status.IsCompleted(){
		return r.applyTransition(ctx, "from completed", Pending, nil)
	}

	if status.IsFailed(){
		return r.applyTransition(ctx, "from failed", Pending, nil)
	}

	// Verify the resource state
	if status.IsVerifying() || status.IsPending() || status.IsSucceeded() || status.IsRecreating() {
		podSpec, _ := r.ResourceManager.Verify(ctx)
		podValue, err := r.spawnPod("verify", podSpec)
		fmt.Printf("%+v\n", err)
		r.instanceUpdater.setPodConfig(podValue)
		return r.applyTransition(ctx, "Check", Checking, nil)
	}

	// dependencies are now satisfied, can now reconcile the manifest and create or update the resource
	if status.IsCreating() {
		podSpec, _ := r.ResourceManager.Create(ctx)
		podValue, err := r.spawnPod("create", podSpec)
		fmt.Printf("%+v\n", err)
		r.instanceUpdater.setPodConfig(podValue)
		return r.applyTransition(ctx, "Check", Checking, nil)
	}

	// **** Updating
	if status.IsUpdating() {
		podSpec, _ := r.ResourceManager.Update(ctx)
		podValue, err := r.spawnPod("update", podSpec)
		fmt.Printf("%+v\n", err)
		r.instanceUpdater.setPodConfig(podValue)
		return r.applyTransition(ctx, "Check", Checking, nil)
	}

	// **** Completing
	// has created or updated, running completion step
	if status.IsCompleting() {
		return r.runCompletion(ctx)
	}

	// **** Terminating
	if status.IsTerminating() {
		r.log.Info("unexpected condition. Terminating state should be handled in finalizer")
		return ctrl.Result{}, nil
	}

	// if has no Status, set to pending
	return r.applyTransition(ctx, "run", Pending, nil)
}

func (r *ReconcileRunner) podDelete(ctx context.Context) (ctrl.Result, error) {
	instance, err := convertInstance(r.instance)

	pod := instance.Status.Pod
	found := &v1.Pod{}

	err = r.KubeClient.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, found)
	terminated := instance.Status.Terminated
	permissions := r.getAccessPermissions()
	r.log.Info(fmt.Sprintf("PodDeleting status: %v", err))
	if err != nil && apierrors.IsNotFound(err) {
		switch pod.Type {
		case "recreate":
			state, err := r.recreate(DeleteResult("Succeeded"))
			return r.applyTransition(ctx, "check", state, err)
		case "create":
			if !permissions.create() {
				// this should never be the case - this is more of an assertion (as the state Verify or Create should never have been set in the first place)
				return r.applyTransition(ctx, "check", Failed, fmt.Errorf(rejectCreateManagedResource))
			}
			return r.apply(ctx, ApplyResponse{
				Result: ApplyResult(statusCodeMap[terminated.ExitCode]),
				Status: terminated.Message,
			})
		case "update":
			if !permissions.update() {
				// this should never be the case - this is more of an assertion (as the state Verify or Create should never have been set in the first place)
				return r.applyTransition(ctx, "check", Failed, fmt.Errorf(rejectCreateManagedResource))
			}
			return r.apply(ctx, ApplyResponse{
				Result: ApplyResult(statusCodeMap[terminated.ExitCode]),
				Status: terminated.Message,
			})
		case "verify":
			return r.verify(ctx, VerifyResponse{
				Result: VerifyResult(statusCodeMap[terminated.ExitCode]),
				Status: terminated.Message,
			})
		}
	}
	return r.applyTransition(ctx, "podDelete", PodDeleting, nil)
}

func (r *ReconcileRunner) check(ctx context.Context) (ctrl.Result, error) {

	instance, err := convertInstance(r.instance)

	pod := instance.Status.Pod
	found := &v1.Pod{}

	err = r.KubeClient.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, found)
	r.log.Info(fmt.Sprintf("PodChecking status: %v", err))
	if err != nil && apierrors.IsNotFound(err) {
		return r.applyTransition(ctx, "check", Pending, err)
	}
	if len(found.Status.ContainerStatuses) == 0 {
		return r.applyTransition(ctx, "check", Checking, errors.New("ContainerStatuses array is of 0 length"))
	}
	terminated := found.Status.ContainerStatuses[0].State.Terminated
	// fmt.Printf("%+v\n", terminated)

	if terminated != nil {
		r.instanceUpdater.setTerminated(terminated)
		r.KubeClient.Delete(ctx, found)
		return r.applyTransition(ctx, "check", PodDeleting, nil)
	}
	return r.applyTransition(ctx, "check", Checking, nil)
}

func (r *ReconcileRunner) setOwner(ctx context.Context, owner runtime.Object) (ctrl.Result, error) {
	//set owner reference if it exists
	r.instanceUpdater.setOwnerReferences([]runtime.Object{owner})
	if err := r.updateAndLog(ctx, v1.EventTypeNormal, "OwnerReferences", "setting OwnerReferences for "+r.Name); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func (r *ReconcileRunner) verify(ctx context.Context, verifyResp VerifyResponse) (ctrl.Result, error) {
	nextState, ensureErr := r.verifyExecute(ctx, verifyResp)
	return r.applyTransition(ctx, "Verify", nextState, ensureErr)
}

const rejectCreateManagedResource = "permission to create resource is not set. Annotation '*/access-permissions' is present, but the flag 'C' is not set"
const rejectUpdateManagedResource = "permission to update resource is not set. Annotation '*/access-permissions' is present, but the flag 'U' is not set"
const rejectDeleteManagedResource = "permission to delete or recreate resource is not set. Annotation '*/access-permissions' is present, but the flag 'D' is not set"

func (r *ReconcileRunner) verifyExecute(ctx context.Context, verifyResp VerifyResponse) (ReconcileState, error) {
	status := r.status
	currentState := status.State

	r.log.Info("Verifying state of resource")
	// verifyResponse, err := r.ResourceManager.Verify(ctx, r.resourceSpec())
	verifyResult := verifyResp.Result
	permissions := r.getAccessPermissions()
	// **** Error
	// either an error was returned from the SDK, or the
	// if err != nil {
	// 	return Failed, err
	// }

	// **** Deleting
	// if the resource is any state where it exists, but the K8s resource is recreating
	// we assume that the resource is being deleted asynchronously, but there is no way to distinguish
	// from the SDK that is deleting - it is either present or not. so we requeue the loop and wait for it to become `missing`
	if verifyResult.deleting() || verifyResult.exists() && status.IsRecreating() {
		r.log.Info("Retrying verification: resource awaiting deletion before recreation can begin, requeuing reconcile loop")
		return currentState, nil
	}

	// **** Ready
	// The resource is finished creating or updating, completion step can take place if necessary
	if verifyResult.ready() {
		// set the Status payload if there is any
		r.instanceUpdater.setStatusPayload(verifyResp.Status)
		nextState := Completed
		// nextState := r.succeedOrComplete()
		// if ( "{{.RunOnce}}" == "1"  ) {
		// 	nextState = Completed
		// }
		return nextState, nil
	}

	// **** Missing
	// We can now create the resource
	if verifyResult.missing() {
		// if not requeing failure, leave as failed
		if status.IsFailed() && r.Parameters.RequeueAfterFailure == 0 {
			return Failed, nil
		}
		if !permissions.create() {
			// fail if permission to create is not present
			return Failed, fmt.Errorf(rejectCreateManagedResource)
		}
		return Creating, nil
	}

	// **** InProgress
	// if still is in progress with create or update, requeue the reconcile loop
	if verifyResult.inProgress() {
		r.log.Info("Retrying verification: create or update in progress, requeuing reconcile loop")
		return currentState, nil
	}

	// **** UpdateRequired
	// The resource exists and is invalid but updateable, so doesn't need to be recreated
	if verifyResult.updateRequired() {
		if !permissions.update() {
			// fail if permission to update is not present
			return Failed, fmt.Errorf(rejectUpdateManagedResource)
		}
		return Updating, nil
	}

	// **** RecreateRequired
	// The resource exists and is invalid and needs to be created
	if verifyResult.recreateRequired() {
		if !permissions.delete() {
			// fail if permission to delete is not present
			return Failed, fmt.Errorf(rejectDeleteManagedResource)
		}
		// deleteResult, err := r.ResourceManager.Delete(ctx, r)
		podSpec, _ := r.ResourceManager.Delete(ctx)
		podValue, err := r.spawnPod("recreate", podSpec)
		fmt.Printf("%+v\n", err)
		r.instanceUpdater.setPodConfig(podValue)
		return Checking, nil
		// if err != nil || deleteResult == DeleteError {
		// 	return Failed, err
		// }

		// // set it back to pending and let it go through the whole process again
		// if deleteResult.awaitingVerification() {
		// 	return Recreating, err
		// }

		// if deleteResult.alreadyDeleted() || deleteResult.succeeded() {
		// 	return Creating, err
		// }

		// return Failed, fmt.Errorf("invalid DeleteResult for %s %s in Verify", r.ResourceKind, r.Name)
	}

	// **** Error
	return Failed, fmt.Errorf("invalid VerifyResult for %s %s in Verify, and no error was specified", r.ResourceKind, r.Name)
}
func (r *ReconcileRunner) recreate(deleteResult DeleteResult) (ReconcileState, error) {
	if deleteResult == DeleteError {
		return Failed, errors.New("Undetermined error")
	}

	// set it back to pending and let it go through the whole process again
	if deleteResult.awaitingVerification() {
		return Recreating, errors.New("Undetermined error")
	}

	if deleteResult.alreadyDeleted() || deleteResult.succeeded() {
		return Creating, errors.New("Undetermined error")
	}

	return Failed, fmt.Errorf("invalid DeleteResult for %s %s in Verify", r.ResourceKind, r.Name)
}
func (r *ReconcileRunner) apply(ctx context.Context, applyResp ApplyResponse) (ctrl.Result, error) {
	r.log.Info("Ready to create or update resource")
	nextState, ensureErr := r.applyExecute(ctx, applyResp)
	return r.applyTransition(ctx, "Ensure", nextState, ensureErr)
}

func (r *ReconcileRunner) applyExecute(ctx context.Context, applyResp ApplyResponse) (ReconcileState, error) {

	resourceName := r.Name
	instance := r.instance
	lastAppliedAnnotation := r.AnnotationBaseName + LastAppliedAnnotation

	// apply that the resource is created or updated (though it won't necessarily be ready, it still needs to be verified)
	// var err error
	// if status.IsCreating() {
	// 	if !permissions.create() {
	// 		// this should never be the case - this is more of an assertion (as the state Verify or Create should never have been set in the first place)
	// 		return Failed, fmt.Errorf(rejectCreateManagedResource)
	// 	}
	// 	applyResponse, err = r.ResourceManager.Create(ctx, r.resourceSpec())
	// } else {
	// 	if !permissions.update() {
	// 		// this should never be the case - this is more of an assertion (as the state Verify or Create should never have been set in the first place)
	// 		return Failed, fmt.Errorf(rejectCreateManagedResource)
	// 	}
	// 	applyResponse, err = r.ResourceManager.Update(ctx, r.resourceSpec())
	// }
	applyResult := applyResp.Result
	if applyResult == "" || applyResult.failed() {
		// clear last update annotation
		r.instanceUpdater.setAnnotation(lastAppliedAnnotation, "")
		errMsg := "Undetermined error"
		// if err != nil {
		// 	errMsg = err.Error()
		// }
		r.Recorder.Event(instance, v1.EventTypeWarning, "Failed", fmt.Sprintf("Couldn't create or update resource: %v", errMsg))
		return Failed, errors.New(errMsg)
	}

	// if successful
	// save the last updated spec as a metadata annotation
	r.instanceUpdater.setAnnotation(lastAppliedAnnotation, r.getJsonSpec())

	// set it to succeeded, completing (if there is a CompletionHandler), or await verification
	if applyResult.awaitingVerification() {
		r.instanceUpdater.setStatusPayload(applyResp.Status)
		return Verifying, nil
	} else if applyResult.succeeded() {
		r.instanceUpdater.setStatusPayload(applyResp.Status)
		return r.succeedOrComplete(), nil
	} else {
		return Failed, fmt.Errorf("invalid response from Create for resource '%s'", resourceName)
	}
}

func (r *ReconcileRunner) succeedOrComplete() ReconcileState {
	if r.CompletionFactory == nil || r.status.IsSucceeded() {
		return Succeeded
	} else {
		return Completing
	}
}

func (r *ReconcileRunner) runCompletion(ctx context.Context) (ctrl.Result, error) {
	var ppError error = nil
	if r.CompletionFactory != nil {
		if handler := r.CompletionFactory(r.GenericController); handler != nil {
			ppError = handler.Run(ctx, r.instance)
		}
	}
	if ppError != nil {
		return r.applyTransition(ctx, "Completion", Failed, ppError)
	} else {
		return r.applyTransition(ctx, "Completion", Succeeded, nil)
	}
}

func (r *ReconcileRunner) updateInstance(ctx context.Context) error {
	if !r.instanceUpdater.hasUpdates() {
		return nil
	}
	return r.tryUpdateInstance(ctx, 5)
}

// this is to get rid of the pesky errors
// "Operation cannot be fulfilled on xxx:the object has been modified; please apply your changes to the latest version and try again"
func (r *ReconcileRunner) tryUpdateInstance(ctx context.Context, count int) error {
	// refetch the instance and apply the updates to it
	baseInstance := r.instance
	instance := baseInstance.DeepCopyObject()
	err := r.KubeClient.Get(ctx, r.NamespacedName, baseInstance)
	if err != nil {
		if apierrors.IsNotFound(err) {
			r.log.Info("Unable to update deleted resource. it may have already been finalized. this error is ignorable. Resource: " + r.Name)
			return nil
		} else {
			r.log.Info("Unable to retrieve resource. falling back to prior instance: " + r.Name + ": err " + err.Error())
		}
	}
	status, _ := r.StatusAccessor(instance)
	err = r.instanceUpdater.applyUpdates(instance, status)
	if err != nil {
		r.log.Info("Unable to convert Object to resource")
		r.instanceUpdater.clear()
		return err
	}
	err = r.KubeClient.Update(ctx, instance)
	if err != nil {
		if count == 0 {
			r.Recorder.Event(instance, v1.EventTypeWarning, "Update", fmt.Sprintf("failed to update  %s instance %s on K8s cluster.", r.ResourceKind, r.Name))
			r.instanceUpdater.clear()
			return err
		}
		r.log.Info(fmt.Sprintf("Failed to update CRD instance on K8s cluster. retries left=%d", count))
		time.Sleep(2 * time.Second)
		return r.tryUpdateInstance(ctx, count-1)
	} else {
		r.instanceUpdater.clear()
		return nil
	}
}

func (r *ReconcileRunner) updateAndLog(ctx context.Context, eventType string, reason string, message string) error {
	instance := r.instance
	if !r.instanceUpdater.hasUpdates() {
		return nil
	}
	if err := r.updateInstance(ctx); err != nil {
		r.log.Info(fmt.Sprintf("K8s update failure: %v", err))
		r.Recorder.Event(instance, v1.EventTypeWarning, reason, fmt.Sprintf("failed to update instance of %s %s in kubernetes cluster", r.ResourceKind, r.Name))
		return err
	}
	r.Recorder.Event(instance, eventType, reason, message)
	return nil
}

func (r *ReconcileRunner) getTransitionDetails(nextState ReconcileState) (ctrl.Result, string) {
	requeueAfter := r.getRequeueAfter(nextState)
	requeueResult := ctrl.Result{Requeue: requeueAfter > 0, RequeueAfter: requeueAfter}
	message := ""
	switch nextState {
	case Pending:
		message = fmt.Sprintf("%s %s in pending state.", r.ResourceKind, r.Name)
	case Creating:
		message = fmt.Sprintf("%s %s ready for creation.", r.ResourceKind, r.Name)
	case Updating:
		message = fmt.Sprintf("%s %s ready to be updated.", r.ResourceKind, r.Name)
	case Verifying:
		message = fmt.Sprintf("%s %s verification in progress.", r.ResourceKind, r.Name)
	case Completing:
		message = fmt.Sprintf("%s %s create or update succeeded and ready for completion step", r.ResourceKind, r.Name)
	case Succeeded:
		message = fmt.Sprintf("%s %s successfully applied and ready for use.", r.ResourceKind, r.Name)
	case Recreating:
		message = fmt.Sprintf("%s %s deleting and recreating in progress.", r.ResourceKind, r.Name)
	case Failed:
		message = fmt.Sprintf("%s %s failed.", r.ResourceKind, r.Name)
	case Checking:
		message = fmt.Sprint("Checking.")
	case PodDeleting:
		message = fmt.Sprint("PodDeleting.")
	case Terminating:
		message = fmt.Sprintf("%s %s termination in progress.", r.ResourceKind, r.Name)
	case Completed:
		message = fmt.Sprintf("completed")
	default:
		message = fmt.Sprintf("%s %s set to state %s", r.ResourceKind, r.Name, nextState)
	}
	return requeueResult, message
}

func (r *ReconcileRunner) applyTransition(ctx context.Context, reason string, nextState ReconcileState, transitionErr error) (ctrl.Result, error) {
	eventType := v1.EventTypeNormal
	if nextState == Failed {
		eventType = v1.EventTypeWarning
	}
	errorMsg := ""
	if transitionErr != nil {
		errorMsg = transitionErr.Error()
	}
	if nextState != r.status.State {
		r.instanceUpdater.setReconcileState(nextState, errorMsg)
	}
	result, transitionMsg := r.getTransitionDetails(nextState)
	updateErr := r.updateAndLog(ctx, eventType, reason, transitionMsg)
	if transitionErr != nil {
		if updateErr != nil {
			// TODO: is the transition error is more important?
			// we don't requeue if there is an update error
			return ctrl.Result{}, transitionErr
		} else {
			return result, nil
		}
	}
	if updateErr != nil {
		return ctrl.Result{}, updateErr
	}
	return result, nil
}

func (r *ReconcileRunner) getRequeueAfter(transitionState ReconcileState) time.Duration {
	parameters := r.Parameters
	requeueAfterDuration := func(requeueSeconds int) time.Duration {
		requeueAfter := time.Duration(requeueSeconds) * time.Millisecond
		return requeueAfter
	}
	
	if transitionState == Completed {
		if ("{{.RunOnce}}" == "1") {
			r.log.Info("Suspended reconciliation as per the configuration")
			return time.Duration(0)
		} else{
			requeueAfter := requeueAfterDuration(parameters.RequeueAfter)
			r.log.Info(fmt.Sprintf("Reconcile runner freq set to - %s", requeueAfter))
			return requeueAfter
		}		
	}

	if transitionState == Verifying ||
		transitionState == PodDeleting ||
		transitionState == Checking ||
		transitionState == Creating ||
		transitionState == Updating ||
		transitionState == Pending ||
		transitionState == Succeeded ||
		transitionState == Recreating {
		// must by default have a non zero requeue for these states
		return requeueAfterDuration(parameters.RequeueAfterSuccess)
	} else if transitionState == Failed {
		r.log.Info("Terminated reconciliation as pod failed")
		return time.Duration(0)
	}
	return 0
}

func (r *ReconcileRunner) getJsonSpec() string {
	fetch := func() (string, error) {
		b, err := json.Marshal(r.instance)
		if err != nil {
			return "", err
		}
		var asMap map[string]interface{}
		err = json.Unmarshal(b, &asMap)
		if err != nil {
			return "", err
		}
		spec := asMap["spec"]
		b, err = json.Marshal(spec)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}

	jsonSpec, err := fetch()
	if err != nil {
		r.log.Info("Error fetching Json for instance spec")
		return ""
	}
	return jsonSpec
}

func (r *ReconcileRunner) getAccessPermissions() AccessPermissions {
	annotations := r.objectMeta.GetAnnotations()
	return AccessPermissions(annotations[r.AnnotationBaseName+AccessPermissionAnnotation])
}

type PodSpec struct {
	Name             string
	Namespace        string
	LogPath          string
	Arguments        string
	Image            string
	ImagePullPolicy  string
	ImagePullSecrets string
}

func (r *ReconcileRunner) spawnPod(event string, podInputSpec PodSpec) ({{ .Resource }}{{ .Version }}.Pod, error) {
	instance, err := getSpec(r.instance)
	var args []string
	for _, v := range strings.Split(podInputSpec.Arguments, " ") {
		vToAppend := v
		switch v {
		case "{spec}":
			tmp, _ := instance.Spec.MarshalJSON()
			vToAppend = string(tmp)
			// vToAppend = instance.Spec
		case "{type}":
			vToAppend = event
		case "{status}":
			tmp, _ := json.Marshal(instance.Status)
			vToAppend = string(tmp)
		case "{logPath}":
			vToAppend = podInputSpec.LogPath
		}
		args = append(args, vToAppend)
	}
	podName := fmt.Sprintf("%s-%s", r.NamespacedName.Name, podInputSpec.Name)
	podNamespace := podInputSpec.Namespace

	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: podNamespace,
			Labels: map[string]string{
				"app": podInputSpec.Name,
			},
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:                   podInputSpec.Name,
					Image:                  podInputSpec.Image,
					Env:                    create_envs(),
					Args:                   args,
					TerminationMessagePath: podInputSpec.LogPath,
					ImagePullPolicy:        v1.PullPolicy(podInputSpec.ImagePullPolicy),
					Resources: 				getResourceRequirements(getResourceList("", ""), getResourceList("{{ .CpuLimit }}", "{{ .MemoryLimit }}")),
				},
			},

			ImagePullSecrets: []v1.LocalObjectReference{{"{{"}}
				Name: podInputSpec.ImagePullSecrets{{"}}"}},
			RestartPolicy:    v1.RestartPolicyNever,
		},
	}

	if err := controllerutil.SetControllerReference(instance, pod, r.GenericController.Scheme); err != nil {
		// requeue with error
		return {{ .Resource }}{{ .Version }}.Pod{}, err
	}

	found := &v1.Pod{}
	err = r.KubeClient.Get(context.TODO(), types.NamespacedName{Name: podName, Namespace: podNamespace}, found)
	// fmt.Println(event, err)
	if err != nil && apierrors.IsNotFound(err) {
		er2 := r.KubeClient.Create(context.TODO(), pod)
		r.log.Info(fmt.Sprintf("SpawnPod status: %v", err))
		if er2 != nil {
			return {{ .Resource }}{{ .Version }}.Pod{}, er2
		}
	}
	return {{ .Resource }}{{ .Version }}.Pod{
		Name:      podName,
		Namespace: podNamespace,
		Type:      event,
	}, nil
}

func getResourceRequirements(requests, limits v1.ResourceList) v1.ResourceRequirements {
	res := v1.ResourceRequirements{}
	res.Requests = requests
	res.Limits = limits
	return res
}

func getResourceList(cpu, memory string) v1.ResourceList {
	res := v1.ResourceList{}
	if cpu != "" {
		res[v1.ResourceCPU] = resource.MustParse(cpu)
	}
	if memory != "" {
		res[v1.ResourceMemory] = resource.MustParse(memory)
	}
	return res
}

func create_envs() []v1.EnvVar {
	filePath := "./controller.json"
	//fmt.Printf("// reading file %s\n", filePath)
	file, err1 := ioutil.ReadFile(filePath)
	if err1 != nil {
		fmt.Printf("// error while reading file %s\n", filePath)
		fmt.Printf("File error: %v\n", err1)
		os.Exit(1)
	}

	var apiconfigs map[string]string

	err2 := json.Unmarshal(file, &apiconfigs)
	if err2 != nil {
		fmt.Println("error:", err2)
	}

	var envs []v1.EnvVar
	for k := range apiconfigs {
		envs = append(envs, v1.EnvVar{Name: k, Value: apiconfigs[k]})
	}

	return envs

}
func (r *ReconcileRunner) get_docker_creds(secret_engine string) (string, string, string) {
	defer func() {
		if err := recover(); err != nil {
			r.log.Info(fmt.Sprintf("Exception occured while secret reveal from vault - %v", err))
			return
		}
	}()

	secret, err := vault.Read(secret_engine)

	if err != nil {
		r.log.Info(fmt.Sprintf("Panic occured while reading secret, error - %s", err.Error()))
		return "", "", ""
	}

	username := secret.Data["username"].(string)
	password := secret.Data["password"].(string)
	server := secret.Data["server"].(string)

	return username, password, server
}

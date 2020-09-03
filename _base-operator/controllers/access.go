// Copyright (c) 2020, salesforce.com, inc.
// All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// For full license text, see LICENSE.txt file in the repo root or https://opensource.org/licenses/BSD-3-Clause

package controllers

import (
	"fmt"

	api "{{ .Repo }}/api/{{  .Version }}"
	"{{ .Repo }}/reconciler"

	"k8s.io/apimachinery/pkg/runtime"
)

func GetStatus(instance runtime.Object) (*reconciler.Status, error) {
	x, err := convertInstance(instance)
	if err != nil {
		return nil, err
	}
	status := x.Status

	return &reconciler.Status{
		State:   reconciler.ReconcileState(status.State),
		Message: status.Message,
	}, nil
}

func updateStatus(instance runtime.Object, status *reconciler.Status) error {
	x, err := convertInstance(instance)
	if err != nil {
		return err
	}
	x.Status.State = string(status.State)
	x.Status.Message = status.Message
	if status.Pod != (api.Pod{}) {
		x.Status.Pod = status.Pod
	}
	x.Status.Terminated = status.Terminated

	switch status.StatusPayload.(type) {
	case string:
		x.Status.StatusPayload = status.StatusPayload.(string)
	}
	return nil
}

func convertInstance(obj runtime.Object) (*api.{{  .Resource }}, error) {
	local, ok := obj.(*api.{{  .Resource }})
	if !ok {
		return nil, fmt.Errorf("failed type assertion on kind: A")
	}
	return local, nil
}

func getSpec(object runtime.Object) (*api.{{  .Resource }}, error) {
	instance, err := convertInstance(object)
	if err != nil {
		return nil, err
	}
	return instance, nil
}

func GetSuccess(object runtime.Object) (bool, error) {
	instance, err := GetStatus(object)
	if err != nil {
		return false, err
	}
	return instance.IsSucceeded(), nil
}

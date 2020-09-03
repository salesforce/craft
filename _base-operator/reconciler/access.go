// Copyright (c) 2020, salesforce.com, inc.
// All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// For full license text, see LICENSE.txt file in the repo root or https://opensource.org/licenses/BSD-3-Clause

package reconciler

import (
	"fmt"

	api "{{ .Repo }}/api/{{  .Version }}"

	"k8s.io/apimachinery/pkg/runtime"
)

var statusCodeMap = map[int32]string{
	201: "Succeeded", // create or update
	202: "AwaitingVerification", // create or update
	203: "Error", // create or update
	211: "Ready", // verify
	212: "InProgress", // verify
	213: "Error", // verify
	214: "Missing", // verify
	215: "UpdateRequired", // verify
	216: "RecreateRequired", // verify
	217: "Deleting", // verify
	221: "Succeeded", // delete
	222: "InProgress", // delete
	223: "Error", // delete
	224: "Missing", // delete
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
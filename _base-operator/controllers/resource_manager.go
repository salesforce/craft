// Copyright (c) 2020, salesforce.com, inc.
// All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// For full license text, see LICENSE.txt file in the repo root or https://opensource.org/licenses/BSD-3-Clause

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"{{ .Repo }}/reconciler"
	"k8s.io/client-go/tools/record"

)


type ResourceManager struct {
	Logger     logr.Logger
	Recorder   record.EventRecorder
}

func (resManager *ResourceManager) Create(ctx context.Context) (reconciler.PodSpec, error) {
	return reconciler.PodSpec{
		Image: "{{ .Image }}",
		Arguments: "--type create --spec {spec}",
		Name: "create",
		Namespace: "{{ .Namespace }}",
		ImagePullPolicy: "{{ .ImagePullPolicy }}",
		ImagePullSecrets: "{{ .ImagePullSecrets }}",
	},nil
}
func (resManager *ResourceManager) Update(ctx context.Context) (reconciler.PodSpec, error) {
	return reconciler.PodSpec{
		Image: "{{ .Image }}",
		Arguments: "--type update --spec {spec}",
		Name: "update",
		Namespace: "{{ .Namespace }}",
		ImagePullPolicy: "{{ .ImagePullPolicy }}",
		ImagePullSecrets: "{{ .ImagePullSecrets }}",
	}, nil
}

func (resManager *ResourceManager) Verify(ctx context.Context)  (reconciler.PodSpec, error) {
	return reconciler.PodSpec{
		Image: "{{ .Image }}",
		Arguments: "--type verify --spec {spec}",
		Name: "verify",
		Namespace: "{{ .Namespace }}",
		ImagePullPolicy: "{{ .ImagePullPolicy }}",
		ImagePullSecrets: "{{ .ImagePullSecrets }}",
	},nil
}

func (resManager *ResourceManager) Delete(ctx context.Context) (reconciler.PodSpec, error) {
	return reconciler.PodSpec{
		Image: "{{ .Image }}",
		Arguments: "--type delete --spec {spec}",
		Name: "delete",
		Namespace: "{{ .Namespace }}",
		ImagePullPolicy: "{{ .ImagePullPolicy }}",
		ImagePullSecrets: "{{ .ImagePullSecrets }}",
	}, nil
}

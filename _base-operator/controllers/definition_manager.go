// Copyright (c) 2020, salesforce.com, inc.
// All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// For full license text, see LICENSE.txt file in the repo root or https://opensource.org/licenses/BSD-3-Clause

package controllers

import (
	"context"

	"{{ .Repo }}/api/{{  .Version }}"
	"{{ .Repo }}/reconciler"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
)

type definitionManager struct{}

func (dm *definitionManager) GetDefinition(ctx context.Context, namespacedName types.NamespacedName) *reconciler.ResourceDefinition {
	return &reconciler.ResourceDefinition{
		InitialInstance: &{{  .Version }}.{{  .Resource }}{},
		StatusAccessor:  GetStatus,
		StatusUpdater:   updateStatus,
	}
}

func (dm *definitionManager) GetDependencies(ctx context.Context, thisInstance runtime.Object) (*reconciler.DependencyDefinitions, error) {
	return &reconciler.NoDependencies, nil
}

// Copyright (c) 2020, salesforce.com, inc.
// All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// For full license text, see LICENSE.txt file in the repo root or https://opensource.org/licenses/BSD-3-Clause

package utils

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	log "github.com/sirupsen/logrus"

	"gopkg.in/yaml.v2"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/validation"
	"k8s.io/apimachinery/pkg/conversion"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func hasAnyStatusEnabled(crd *apiextensions.CustomResourceDefinitionSpec) bool {
	if hasStatusEnabled(crd.Subresources) {
		return true
	}
	for _, v := range crd.Versions {
		if hasStatusEnabled(v.Subresources) {
			return true
		}
	}
	return false
}

// hasStatusEnabled returns true if given CRD Subresources has non-nil Status set.
func hasStatusEnabled(subresources *apiextensions.CustomResourceSubresources) bool {
	if subresources != nil && subresources.Status != nil {
		return true
	}
	return false
}

func Validate(crdPath string) {
	jsonFile, err := ioutil.ReadFile(crdPath)
	if err != nil {
		log.Fatal(err)
	}
	sepYamlfiles := strings.Split(string(jsonFile), "---")
	for _, f := range sepYamlfiles {
		if strings.Contains(f, "kind: CustomResourceDefinition") {
			var obj v1beta1.CustomResourceDefinition
			var apiObj apiextensions.CustomResourceDefinition
			var s conversion.Scope
			var body interface{}
			yaml.Unmarshal([]byte(f), &body)
			body = convert(body)
			if b, err := json.Marshal(body); err != nil {
				log.Fatal(err)
			} else {
				json.Unmarshal(b, &obj)
				// log.Debugf("%+v\n", obj)
				v1beta1.Convert_v1beta1_CustomResourceDefinition_To_apiextensions_CustomResourceDefinition(&obj, &apiObj, s)
				// log.Debugf("%+v\n", apiObj.Spec.Validation.OpenAPIV3Schema)
				// version v0.18.2
				requestGV := schema.GroupVersion{
					Group:   apiObj.Spec.Group,
					Version: apiObj.Spec.Version,
				}
				errList := validation.ValidateCustomResourceDefinition(&apiObj, requestGV)
				for _, e := range errList {
					if !strings.Contains(e.Error(), "status.storedVersion") {
						log.Warnf("Error: %s\n", e)
					}
				}
			}
		}

	}
}

func convert(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = convert(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = convert(v)
		}
	}
	return i
}

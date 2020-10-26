// Copyright (c) 2020, salesforce.com, inc.
// All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// For full license text, see LICENSE.txt file in the repo root or https://opensource.org/licenses/BSD-3-Clause

package cmd

import (
	"crypto/rand"
	"encoding/hex"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

var (
	apiFile    string
	apiFileObj ApiFileStruct
)

type ApiFileStruct struct {
	Group            string `json:"group" yaml:"group"`
	Repo             string `json:"repo" yaml:"repo"`
	Resource         string `json:"resource" yaml:"resource"`
	Domain           string `json:"domain" yaml:"domain"`
	Version          string `json:"version" yaml:"version"`
	Namespace        string `json:"namespace" yaml:"namespace"`
	Image            string `json:"image" yaml:"image"`
	CpuLimit         string `json:"cpu_limit" yaml:"cpu_limit"`
	MemoryLimit      string `json:"memory_limit" yaml:"memory_limit"`
	ImagePullPolicy  string `json:"imagePullPolicy" yaml:"imagePullPolicy"`
	ImagePullSecrets string `json:"imagePullSecrets" yaml:"imagePullSecrets"`
	OperatorImage    string `json:"operator_image" yaml:"operator_image"`
	ReconcileFreq    string `json:"reconcileFreq" yaml:"reconcileFreq"`
	RunOnce          string `json:"runOnce" yaml:"runOnce"`
	LeaderElectionID string
	LowerRes         string
}

func (api *ApiFileStruct) loadApi(path string) {
	stream, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	_ = yaml.Unmarshal([]byte(stream), &api)

	api.LeaderElectionID = randomHex()
}

func randomHex() string {
	bytes := make([]byte, 20)
	if _, err := rand.Read(bytes); err != nil {
		return "123.salesforce.com"
	}

	return hex.EncodeToString(bytes)
}

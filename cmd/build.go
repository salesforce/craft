// Copyright (c) 2020, salesforce.com, inc.
// All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// For full license text, see LICENSE.txt file in the repo root or https://opensource.org/licenses/BSD-3-Clause

package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	pathLib "path"

	"craft/utils"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	apps "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/client-go/kubernetes/scheme"
)

func buildCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "build",
		Aliases: []string{"b"},
		Short:   "for building (code | deploy| images)",
		Long:    `for building (code |deploy| images)`,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	cmd.AddCommand(
		imageBuildCmd(),
		deployBuildCmd(),
		codeBuildCmd(),
	)
	cmd.PersistentFlags().StringVarP(&apiFile, "controllerFile", "c", "", "controller file with group, resource and other info")
	cmd.MarkPersistentFlagRequired("controllerFile")

	if err := viper.BindPFlag("controllerFile", cmd.Flags().Lookup("controllerFile")); err != nil {
		log.Fatal(err)
	}
	return cmd
}

var (
	environ string
)

func RWYaml(deployFile string) {
	sch := runtime.NewScheme()
	_ = scheme.AddToScheme(sch)
	_ = apiextv1beta1.AddToScheme(sch)
	decode := serializer.NewCodecFactory(sch).UniversalDeserializer().Decode
	stream, err := ioutil.ReadFile(deployFile)
	if err != nil {
		log.Fatal(err)
	}
	objList := strings.Split(fmt.Sprintf("%s", stream), "---\n")
	newFile, err := os.Create(deployFile)
	if err != nil {
		log.Fatal(err)
	}
	defer newFile.Close()
	y := printers.YAMLPrinter{}
	for _, f := range objList {
		obj, gKV, err := decode([]byte(f), nil, nil)
		if err != nil {
			log.Println(fmt.Sprintf("Error while decoding YAML object. Err was: %s", err))
			continue
		}
		switch gKV.Kind {
		case "Namespace":
			n := obj.(*corev1.Namespace)
			n.ObjectMeta.Name = apiFileObj.Namespace
			y.PrintObj(n, newFile)
		case "Deployment":
			n := obj.(*apps.Deployment)
			n.Spec.Template.Spec.ImagePullSecrets = []corev1.LocalObjectReference{
				{Name: apiFileObj.ImagePullSecrets},
			}
			y.PrintObj(n, newFile)
		case "CustomResourceDefinition":
			n := obj.(*apiextv1beta1.CustomResourceDefinition)
			var m apiextv1beta1.JSONSchemaProps
			s, err := ioutil.ReadFile(resourceFile)
			if err != nil {
				log.Fatal(err)
			}
			json.Unmarshal(s, &m)
			n.Spec.Validation.OpenAPIV3Schema.Properties["spec"] = m
			y.PrintObj(n, newFile)
		default:
			y.PrintObj(obj, newFile)
		}
	}
}

func deployBuildCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "deploy",
		Aliases: []string{"d"},
		Short:   "build deploy",
		Long:    `build deploy`,
		Run: func(cmd *cobra.Command, args []string) {
			absPath()
			apiFileObj.loadApi(apiFile)
			baseDir := filepath.Dir(apiFile)
			newOperatorPath := pathLib.Join(goSrc, apiFileObj.Repo)

			deployPath := pathLib.Join(baseDir, "deploy")
			deployFile := pathLib.Join(deployPath, "operator.yaml")

			log.Debugf("Mkdir deploy %s", deployPath)
			os.MkdirAll(deployPath, os.ModePerm)

			utils.EnvCmdExec("go build -a -o bin/manager main.go",
				newOperatorPath,
				[]string{"CGO_ENABLED=0", "GOOS=linux", "GOARCH=amd64", "GO111MODULE=on"})

			cmdString := fmt.Sprintf("make operator IMG=%s NAMESPACE=%s FILE=%s",
				apiFileObj.OperatorImage,
				apiFileObj.Namespace,
				deployFile)
			utils.CmdExec(cmdString, newOperatorPath)

			RWYaml(deployFile)

			utils.Validate(deployFile)

		},
	}
	cmd.PersistentFlags().StringVarP(&resourceFile, "resourceFile", "r", "", "resourcefile with properties of resource")
	cmd.PersistentFlags().StringVarP(&environ, "environ", "e", "", "which environment to use for envyaml")

	cmd.MarkPersistentFlagRequired("resourceFile")

	if err := viper.BindPFlag("resourceFile", cmd.Flags().Lookup("resourceFile")); err != nil {
		log.Fatal(err)
	}
	if err := viper.BindPFlag("environ", cmd.Flags().Lookup("environ")); err != nil {
		log.Fatal(err)
	}
	return cmd
}

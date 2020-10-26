// Copyright (c) 2020, salesforce.com, inc.
// All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// For full license text, see LICENSE.txt file in the repo root or https://opensource.org/licenses/BSD-3-Clause

package cmd

import (
	"craft/utils"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	pathLib "path"
	"path/filepath"
	"strings"
	"text/template"
)

var (
	resourceFile string
	baseOperator string
)

func renderTemplate(operatorPath string) {
	resourceDef := fmt.Sprintf("api/%s/%s_types.go",
		apiFileObj.Version,
		strings.ToLower(apiFileObj.Resource))

	utils.MinCmdExec("svn checkout https://github.com/salesforce/craft/trunk/_base-operator", operatorPath)
	baseOperator = pathLib.Join(operatorPath, "_base-operator")

	dirs := []string{"controllers", "reconciler", "main.go", "Dockerfile", "v1/resource.go", "Makefile"}
	for _, dir := range dirs {
		path := pathLib.Join(baseOperator, dir)

		err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if utils.FileExists(path) {
				tpl, err := template.ParseFiles(path)
				if err != nil {
					log.Fatal(err)
				}

				newPath := strings.Replace(path, baseOperator, operatorPath, 1)
				if strings.HasSuffix(path, "v1/resource.go") {
					newPath = pathLib.Join(operatorPath, resourceDef)
				}
				log.Debugf("Rendering file: %s", newPath)
				fi, err := os.Create(newPath)
				if err != nil {
					log.Fatal(err)
				}

				err = tpl.Execute(fi, apiFileObj)

				if err != nil {
					log.Fatal(err)
				}
			}
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
	}

	utils.CmdExec("rm -rf _base-operator", operatorPath)
}

func cpFile(operatorPath string) {
	dirs := []string{"controllers", "reconciler"}
	for _, dir := range dirs {
		pth := pathLib.Join(operatorPath, dir)
		os.MkdirAll(pth, os.ModePerm)
	}
}

func cpAPIFile(apiFile string, operatorPath string) {
	input, err := ioutil.ReadFile(apiFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	dstFile := filepath.Join(operatorPath, "controller.json")
	err = ioutil.WriteFile(dstFile, input, 0644)
	if err != nil {
		fmt.Println("Error creating", dstFile)
		fmt.Println(err)
		return
	}
}

func codeBuildCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "code",
		Aliases: []string{"c"},
		Short:   "create operator template in $GOPATH/src",
		Long:    `create operator template in $GOPATH/src`,
		Run: func(cmd *cobra.Command, args []string) {
			absPath()
			apiFileObj.loadApi(apiFile)

			apiFileObj.LowerRes = strings.ToLower(apiFileObj.Resource)
			var kubeCmdString string
			newOperatorPath := pathLib.Join(goSrc, apiFileObj.Repo)

			os.RemoveAll(newOperatorPath)
			os.MkdirAll(newOperatorPath, os.ModePerm)

			kubeCmdString = fmt.Sprintf("kubebuilder init --domain %s --repo %s", apiFileObj.Domain, apiFileObj.Repo)
			utils.CmdExec(kubeCmdString, newOperatorPath)

			kubeCmdString = fmt.Sprintf("kubebuilder create api --group %s --version %s --kind %s --resource=true --controller=true",
				apiFileObj.Group,
				apiFileObj.Version,
				apiFileObj.Resource,
			)
			utils.CmdExec(kubeCmdString, newOperatorPath)

			utils.CmdExec("rm -rf controllers", newOperatorPath)

			cpFile(newOperatorPath)
			cpAPIFile(apiFile, newOperatorPath)

			kubeCmdString = fmt.Sprintf("rm -rf api/%s/%s_types.go",
				apiFileObj.Version,
				apiFileObj.LowerRes)
			utils.CmdExec(kubeCmdString, newOperatorPath)

			kubeCmdString = fmt.Sprintf("schema-generate -p %s -o api/%s/spec_type.go %s",
				apiFileObj.Version,
				apiFileObj.Version,
				resourceFile)
			utils.CmdExec(kubeCmdString, newOperatorPath)

			renderTemplate(newOperatorPath)
			utils.CmdExec("make generate", newOperatorPath)
		},
	}

	cmd.PersistentFlags().StringVarP(&apiFile, "controllerFile", "c", "", "controller file with group, resource and other info")
	cmd.PersistentFlags().StringVarP(&resourceFile, "resourceFile", "r", "", "resourcefile with properties of resource")
	cmd.MarkPersistentFlagRequired("controllerFile")
	cmd.MarkPersistentFlagRequired("resourceFile")

	if err := viper.BindPFlag("controllerFile", cmd.Flags().Lookup("controllerFile")); err != nil {
		log.Fatal(err)
	}
	if err := viper.BindPFlag("resourceFile", cmd.Flags().Lookup("resourceFile")); err != nil {
		log.Fatal(err)
	}

	return cmd
}

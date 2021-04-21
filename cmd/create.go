// Copyright (c) 2020, salesforce.com, inc.
// All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// For full license text, see LICENSE.txt file in the repo root or https://opensource.org/licenses/BSD-3-Clause

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"craft/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	dockerPush bool
)

func absAPIPath() {
	var err error
	apiFile, err = filepath.Abs(apiFile)
	if err != nil {
		log.Fatal(err)
	}
	log.Debug("apiFile: ", apiFile)
}
func absResourcePath() {
	var err error
	resourceFile, err = filepath.Abs(resourceFile)
	if err != nil {
		log.Fatal(err)
	}
	log.Debug("resourceFile: ", resourceFile)

}
func absPodPath() {
	var err error
	podDockerFile, err = filepath.Abs(podDockerFile)
	if err != nil {
		log.Fatal(err)
	}
	log.Debug("podDockerFile: ", podDockerFile)
}
func absPath() {
	absAPIPath()
	absResourcePath()
	absPodPath()
}

func createCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Aliases: []string{"c"},
		Short:   "create operator in $GOPATH/src",
		Long:    `create operator in $GOPATH/src`,
		Run: func(cmd *cobra.Command, args []string) {
			absPath()
			apiFileObj.loadApi(apiFile)
			pwd, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}
			total := time.Now()
			fmt.Println("Creating operator in $GOPATH/src")
			start := time.Now()
			cmdString := fmt.Sprintf("craft build code -c %s -r %s", apiFile, resourceFile)
			utils.CmdExec(cmdString, pwd)
			elapsed:= time.Since(start)
			fmt.Println("Created operator in $GOPATH/src in", elapsed.Round(time.Second).String())

			fmt.Println("Building operator.yaml for deployment")
			start = time.Now()
			cmdString = fmt.Sprintf("craft build deploy -c %s -r %s", apiFile, resourceFile)
			utils.CmdExec(cmdString, pwd)
			elapsed = time.Since(start)
			fmt.Println("Built operator.yaml for deployment in", elapsed.Round(time.Second).String())

			fmt.Println("Building operator and resource docker images")
			start = time.Now()
			cmdString = fmt.Sprintf("craft build image -b -c %s --podDockerFile %s", apiFile, podDockerFile)
			utils.CmdExec(cmdString, pwd)
			elapsed = time.Since(start)
			fmt.Println("Built operator and resource docker images in", elapsed.Round(time.Second).String())

			if dockerPush {
				fmt.Println("Pushing operator image to docker")
				start = time.Now()
				cmdString = fmt.Sprintf("docker push %s", apiFileObj.OperatorImage)
				utils.CmdExec(cmdString, pwd)
				elapsed = time.Since(start)
				fmt.Println("Pushed operator image to docker in", elapsed.Round(time.Second).String())

				fmt.Println("Pushing resource image to docker")
				start = time.Now()
				cmdString = fmt.Sprintf("docker push %s", apiFileObj.Image)
				utils.CmdExec(cmdString, pwd)
				elapsed = time.Since(start)
				fmt.Println("Pushed resource image to docker in", elapsed.Round(time.Second).String())
			}
			totalTime := time.Since(total)
			fmt.Println("Total time required for completion is", totalTime.Round(time.Second).String())
		},
	}

	cmd.PersistentFlags().StringVarP(&apiFile, "controllerFile", "c", "", "controller file with group, resource and other info")
	cmd.PersistentFlags().StringVarP(&resourceFile, "resourceFile", "r", "", "resourcefile with properties of resource")
	cmd.PersistentFlags().StringVarP(&podDockerFile, "podDockerFile", "P", "", "pod Dockerfile")
	cmd.PersistentFlags().BoolVarP(&dockerPush, "push", "p", false, "If set to true, pushes images to docker")
	cmd.MarkPersistentFlagRequired("controllerFile")
	cmd.MarkPersistentFlagRequired("resourceFile")
	cmd.MarkPersistentFlagRequired("podDockerFile")
	// cmd.MarkPersistentFlagRequired("environ")

	if err := viper.BindPFlag("controllerFile", cmd.Flags().Lookup("controllerFile")); err != nil {
		log.Fatal(err)
	}
	if err := viper.BindPFlag("resourceFile", cmd.Flags().Lookup("resourceFile")); err != nil {
		log.Fatal(err)
	}
	if err := viper.BindPFlag("podDockerFile", cmd.Flags().Lookup("podDockerFile")); err != nil {
		log.Fatal(err)
	}
	return cmd
}
// Copyright (c) 2020, salesforce.com, inc.
// All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// For full license text, see LICENSE.txt file in the repo root or https://opensource.org/licenses/BSD-3-Clause

package cmd

import (
	"fmt"
	"os"
	pathLib "path"
	"path/filepath"

	"craft/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	podDockerFile      string
	operatorDockerFile string
	buildImage         bool
)

func absOperatorPath() {
	var err error
	operatorDockerFile, err = filepath.Abs(operatorDockerFile)
	if err != nil {
		log.Fatal(err)
	}
	log.Debug("operatorDockerFile: ", operatorDockerFile)
}
func logDockerBuild(cmd, dir, imageName string) {
	log.Infof("cd %s", dir)
	log.Infof("Run command for building %s : %s", imageName, cmd)
}
func imageBuildCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "image",
		Aliases: []string{"i"},
		Short:   "build images",
		Long:    `build images`,
		Run: func(cmd *cobra.Command, args []string) {
			absAPIPath()
			apiFileObj.loadApi(apiFile)
			newOperatorPath := pathLib.Join(goSrc, apiFileObj.Repo)

			var podDockerDir, operatorDockerDir string

			if podDockerFile == "" {
				podDockerDir = filepath.Dir(apiFile)
				podDockerFile = pathLib.Join(podDockerDir, "Dockerfile")
			} else {
				absPodPath()
				podDockerDir = filepath.Dir(podDockerFile)
			}
			if operatorDockerFile == "" {
				operatorDockerDir = newOperatorPath
				operatorDockerFile = pathLib.Join(operatorDockerDir, "Dockerfile")
			} else {
				absOperatorPath()
				operatorDockerDir = filepath.Dir(operatorDockerFile)
			}

			if !utils.FileExists(operatorDockerFile) {
				log.Fatal("specify -o for operator docker file.")
			}
			if !utils.FileExists(podDockerFile) {
				log.Fatal("specify -p for pod docker file.")
			}

			if !buildImage {
				log.Info(`


				Pass -b for craft to build image (but since building image takes long time you won't see any output during that)
				At end we specified instructions for building images yourself
				
				
				`)
			}
			imageCmd := fmt.Sprintf("docker build -t %s -f %s .",
				apiFileObj.OperatorImage,
				operatorDockerFile)
			if buildImage {
				utils.CmdExec(imageCmd, operatorDockerDir)
			} else {
				logDockerBuild(imageCmd, operatorDockerDir, "operator")
			}

			imageCmd = fmt.Sprintf("docker build --build-arg vault_token=%s -t %s -f %s .",
				os.Getenv("VAULT_TOKEN"),
				apiFileObj.Image,
				podDockerFile)
			if buildImage {
				utils.CmdExec(imageCmd, podDockerDir)
			} else {
				logDockerBuild(imageCmd, podDockerDir, "pod")
			}
		},
	}
	cmd.PersistentFlags().StringVarP(&podDockerFile, "podDockerFile", "p", "", "pod Dockerfile")
	cmd.PersistentFlags().StringVarP(&operatorDockerFile, "operatorDockerFile", "o", "", "pod Dockerfile")
	cmd.PersistentFlags().BoolVarP(&buildImage, "build", "b", false, "pod Dockerfile")

	if err := viper.BindPFlag("podDockerFile", cmd.Flags().Lookup("podDockerFile")); err != nil {
		log.Fatal(err)
	}
	if err := viper.BindPFlag("operatorDockerFile", cmd.Flags().Lookup("operatorDockerFile")); err != nil {
		log.Fatal(err)
	}
	return cmd
}

// Copyright (c) 2020, salesforce.com, inc.
// All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// For full license text, see LICENSE.txt file in the repo root or https://opensource.org/licenses/BSD-3-Clause

package cmd

import (
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"

	"craft/utils"
	"github.com/spf13/cobra"
)

func initCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "init",
		Aliases: []string{"i"},
		Short:   "init folder with craft sample declaration",
		Long:    `init folder with craft sample declaration`,
		Run: func(cmd *cobra.Command, args []string) {
			folderName, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}
			folder, err := filepath.Abs(folderName)
			if err != nil {
				log.Fatal(err)
			}

			cmdString := "svn checkout https://github.com/salesforce/craft/trunk/init"
			utils.MinCmdExec(cmdString, folder)
			log.Infof("Created sample controller and resource files")
		},
	}

	return cmd
}

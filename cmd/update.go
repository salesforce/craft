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
	"os"
	"runtime"
)

func updateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update",
		Aliases: []string{"u"},
		Short:   "update existing version of craft to latest version",
		Long:    `update existing version of craft to latest version`,
		Run: func(cmd *cobra.Command, args []string) {
			path, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}
			goos := runtime.GOOS
			cmdStr := fmt.Sprintf("wget https://github.com/salesforce/craft/releases/latest/download/craft_%s.tar.gz", goos)
			utils.MinCmdExec(cmdStr, path)
			cmdStr = fmt.Sprintf("sudo tar -xzf craft_%s.tar.gz -C /usr/local/craft", goos)
			utils.CmdExec(cmdStr, path)
			cmdStr = fmt.Sprintf("sudo rm -rf craft_%s.tar.gz", goos)
			utils.CmdExec(cmdStr, path)
		},
	}

	return cmd
}
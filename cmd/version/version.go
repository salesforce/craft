// Copyright (c) 2020, salesforce.com, inc.
// All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// For full license text, see LICENSE.txt file in the repo root or https://opensource.org/licenses/BSD-3-Clause

package version

import (
	"fmt"
	"os"

	"craft/cmd/base"
	"github.com/spf13/cobra"
)

func VersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Displays the version of the current build of craft",
		Long:  `Displays the version of the current build of craft`,
		Run: func(cmd *cobra.Command, args []string) {
				fmt.Printf(base.VersionStr,
					base.Info["version"],
					base.Info["revision"],
					base.Info["branch"],
					base.Info["buildUser"],
					base.Info["buildDate"],
					base.Info["goVersion"])
				os.Exit(0)
		},
	}

	return cmd
}
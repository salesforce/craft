// Copyright (c) 2020, salesforce.com, inc.
// All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// For full license text, see LICENSE.txt file in the repo root or https://opensource.org/licenses/BSD-3-Clause

package cmd

import (
	"log"
	"path/filepath"

	"craft/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var crdPath string

func validateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "validate",
		Aliases: []string{"v"},
		Short:   "validate crd",
		Long:    `validate crd`,
		Run: func(cmd *cobra.Command, args []string) {
			crdPath, err := filepath.Abs(crdPath)
			if err != nil {
				log.Fatal(err)
			}
			utils.Validate(crdPath)
		},
	}
	cmd.PersistentFlags().StringVarP(&crdPath, "crdPath", "v", "", "path to crd definition")
	if err := viper.BindPFlag("crdPath", rootCmd.Flags().Lookup("crdPath")); err != nil {
		log.Fatal(err)
	}
	cmd.MarkPersistentFlagRequired("crdPath")

	return cmd
}

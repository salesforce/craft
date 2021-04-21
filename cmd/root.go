// Copyright (c) 2020, salesforce.com, inc.
// All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// For full license text, see LICENSE.txt file in the repo root or https://opensource.org/licenses/BSD-3-Clause

package cmd

import (
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	"craft/cmd/version"
	"craft/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	craftInstallPath = "/usr/local/craft"
	goSrc            = os.ExpandEnv("$GOPATH/src")
	craftDir         string
	debug            bool
)

func setCraftDir() {
	var err error
	craftDir, err = filepath.Abs(craftDir)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("CraftDir: ", craftDir)
}
func initLoad() {
	setCraftDir()
	utils.CheckGoPath()
	log.SetOutput(os.Stdout)
	if debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}
func init() {
	cobra.OnInitialize(initLoad)
	cobra.EnableCommandSorting = false

	rootCmd.SilenceUsage = true

	// Register subcommands
	rootCmd.AddCommand(
		version.VersionCmd(),
		createCmd(),
		initCmd(),
		validateCmd(),
		buildCmd(),
		updateCmd(),
	)
	rootCmd.PersistentFlags().StringVarP(&craftDir, "craftDir", "C", craftInstallPath, "craft dir")
	rootCmd.PersistentFlags().MarkHidden("craftDir")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "debug")

	if err := viper.BindPFlag("craftDir", rootCmd.Flags().Lookup("craftDir")); err != nil {
		log.Fatal(err)
	}
	if err := viper.BindPFlag("debug", rootCmd.Flags().Lookup("debug")); err != nil {
		log.Fatal(err)
	}
}

var rootCmd = &cobra.Command{
	Use:   "craft",
	Short: "Craft is tool for creating generic operator",
	Long:  strings.TrimSpace(``),
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

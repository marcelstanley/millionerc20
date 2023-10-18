// Copyright (c) Marcel Moura
// SPDX-License-Identifier: MIT (see LICENSE)

package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

const CLI_NAME = "millionerc20"
const CLI_LOG_PREFIX = "[millionerc20] "

var (
	rootCmd = &cobra.Command{
		Use:   CLI_NAME,
		Short: "A CLI for the Million ERC-20 DApp",
		Long:  "A CLI for the Million ERC-20 DApp",
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// Copyright (c) Marcel Moura
// SPDX-License-Identifier: MIT (see LICENSE)

package main

import (
	"log"

	"github.com/marcelstanley/millionerc20/cli/cmd"
)

func init() {
	// TODO Make sure this is set globally
	log.SetPrefix(cmd.CLI_LOG_PREFIX)
}

func main() {
	cmd.Execute()
}

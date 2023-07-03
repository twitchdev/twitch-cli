// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package cmd

import (
	"fmt"

	"github.com/twitchdev/twitch-cli/internal/util"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Returns the current version of the CLI.",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("twitch-cli/" + util.GetVersion())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

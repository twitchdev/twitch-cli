// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/twitchdev/twitch-cli/internal/drops"
)

var (
	gameID   string
	userID   string
	filename string
)

// dropsCmd represents the drops command
var dropsCmd = &cobra.Command{
	Use:   "drops",
	Short: "Used to interface with Drops services.",
}

var exportDropsCmd = &cobra.Command{
	Use:   "export",
	Short: "Exports a CSV with a list of entitlements from the stored Client ID",
	Run:   runDropsCmd,
}

func init() {
	rootCmd.AddCommand(dropsCmd)
	dropsCmd.AddCommand(exportDropsCmd)

	exportDropsCmd.Flags().StringVarP(&filename, "filename", "f", "", "Filename to write the output to.")
	exportDropsCmd.Flags().StringVarP(&gameID, "game-id", "g", "", "ID of the game to get entitlements for.")
	exportDropsCmd.Flags().StringVarP(&userID, "user-id", "u", "", "ID of the user to get entitlements for.")
	exportDropsCmd.MarkFlagRequired("filename")

}

func runDropsCmd(cmd *cobra.Command, args []string) {
	drops.ExportEntitlements(filename, gameID, userID)
}

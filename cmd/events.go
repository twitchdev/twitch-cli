// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/twitchdev/twitch-cli/cmd/events"
)

var noConfig bool

var eventCmd = &cobra.Command{
	Use:   "event",
	Short: "Used to interface with EventSub topics.",
}

func init() {
	rootCmd.AddCommand(eventCmd)

	eventCmd.AddCommand(
		events.TriggerCommand(),
		events.RetriggerCommand(),
		events.VerifySubscriptionCommand(),
		events.WebsocketCommand(),
		events.StartWebsocketServerCommand(),
		events.ConfigureCommand(),
	)

	eventCmd.Flags().BoolVarP(&noConfig, "no-config", "D", false, "Disables the use of the configuration, if it exists.")
}

// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package cmd

import (
	"fmt"
	"net/url"

	trigger "github.com/twitchdev/twitch-cli/internal/events"

	"github.com/spf13/cobra"
)

var (
	isAnonymous    bool
	forwardAddress string
	event          string
	transport      string
	fromUser       string
	toUser         string
	giftUser       string
	eventID        string
	secret         string
	count          int
)

var eventCmd = &cobra.Command{
	Use:   "event",
	Short: "Used to interface with Event services, such as Eventsub and Websub.",
}

var triggerCmd = &cobra.Command{
	Use:   "trigger [event]",
	Short: "Creates mock events that can be forwarded to a local webserver for event testing.",
	Long: fmt.Sprintf(`Creates mock events that can be forwarded to a local webserver for event testing.
	Supported:
	%s`, trigger.ValidTriggers()),
	Args:      cobra.MaximumNArgs(1),
	ValidArgs: trigger.ValidTriggers(),
	Run:       triggerCmdRun,
	Example:   `twitch trigger subscribe`,
	Aliases: []string{
		"fire", "emit",
	},
}

var retriggerCmd = &cobra.Command{
	Use:     "retrigger",
	Short:   "Refires events based on the event ID. Can be forwarded to the local webserver for event testing.",
	Run:     retriggerCmdRun,
	Example: `twitch trigger subscribe`,
}

func init() {
	rootCmd.AddCommand(eventCmd)
	eventCmd.AddCommand(triggerCmd, retriggerCmd)

	triggerCmd.Flags().StringVarP(&forwardAddress, "forward-address", "F", "", "Forward address for mock event.")
	triggerCmd.Flags().StringVarP(&transport, "transport", "T", "eventsub", fmt.Sprintf("Preferred transport method for event. Defaults to Webhooks 2/EventSub.\nSupported values: %s", trigger.ValidTransports()))
	triggerCmd.Flags().StringVarP(&toUser, "to-user", "t", "", "User ID of the receiver of the event. For example, the user that receives a follow. In most contexts, this is the broadcaster.")
	triggerCmd.Flags().StringVarP(&fromUser, "from-user", "f", "", "User ID of the user sending the event, for example the user following another user.")
	triggerCmd.Flags().StringVarP(&giftUser, "gift-user", "g", "", "Used only for \"gift\" events. Denotes the User ID of the gifting user.")
	triggerCmd.Flags().BoolVarP(&isAnonymous, "anonymous", "a", false, "Denotes if the event is anonymous. Only applies to Gift and Sub events.")
	triggerCmd.Flags().IntVarP(&count, "count", "c", 1, "Count of events to trigger. This will simulate a sub gift, or large number of cheers.")
	triggerCmd.Flags().StringVarP(&secret, "secret", "s", "", "Webhook secret. If defined, signs all forwarded events with the SHA256 HMAC.")

	retriggerCmd.Flags().StringVarP(&forwardAddress, "forward-address", "F", "", "Forward address for mock event.")
	retriggerCmd.Flags().StringVarP(&eventID, "id", "i", "", "ID of the event to be refired.")
	retriggerCmd.Flags().StringVarP(&secret, "secret", "s", "", "Webhook secret. If defined, signs all forwarded events with the SHA256 HMAC.")

	retriggerCmd.MarkFlagRequired("id")
}

func triggerCmdRun(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		cmd.Help()
		return
	}

	// Validate that the forward address is actually a URL
	if len(forwardAddress) > 0 {
		_, err := url.ParseRequestURI(forwardAddress)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	for i := 0; i < count; i++ {
		res, err := trigger.Fire(trigger.TriggerParameters{
			Event:          args[0],
			Transport:      transport,
			ForwardAddress: forwardAddress,
			FromUser:       fromUser,
			ToUser:         toUser,
			GiftUser:       giftUser,
			Secret:         secret,
			IsAnonymous:    isAnonymous,
		})

		if err != nil {
			println(err.Error())
			return
		}

		fmt.Println(res)
	}
}

func retriggerCmdRun(cmd *cobra.Command, args []string) {
	res, err := trigger.RefireEvent(eventID, trigger.TriggerParameters{
		ForwardAddress: forwardAddress,
		Secret:         secret,
	})
	if err != nil {
		fmt.Printf("Error refiring event: %s", err)
		return
	}

	fmt.Println(res)
}

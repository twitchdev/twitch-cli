// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package cmd

import (
	"fmt"
	"log"
	"net/url"

	"github.com/spf13/cobra"
	"github.com/twitchdev/twitch-cli/internal/events"
	"github.com/twitchdev/twitch-cli/internal/events/mock_wss_server"
	"github.com/twitchdev/twitch-cli/internal/events/trigger"
	"github.com/twitchdev/twitch-cli/internal/events/verify"
)

const websubDeprecationNotice = "Halt! It appears you are trying to use WebSub, which has been deprecated. For more information, see: https://discuss.dev.twitch.tv/t/deprecation-of-websub-based-webhooks/32152"

var (
	isAnonymous        bool
	forwardAddress     string
	event              string
	transport          string
	fromUser           string
	toUser             string
	giftUser           string
	eventID            string
	secret             string
	eventStatus        string
	subscriptionStatus string
	itemID             string
	itemName           string
	cost               int64
	count              int
	description        string
	gameID             string
	timestamp          string
	debug              bool
	wssReconnectTest   int
	sslEnabled         bool
)

var eventCmd = &cobra.Command{
	Use:   "event",
	Short: "Used to interface with EventSub topics.",
}

var triggerCmd = &cobra.Command{
	Use:   "trigger [event]",
	Short: "Creates mock events that can be forwarded to a local webserver for event testing.",
	Long: fmt.Sprintf(`Creates mock events that can be forwarded to a local webserver for event testing.
	Supported:
	%s`, events.ValidTriggers()),
	Args:      cobra.MaximumNArgs(1),
	ValidArgs: events.ValidTriggers(),
	Run:       triggerCmdRun,
	Example:   `twitch event trigger subscribe`,
	Aliases: []string{
		"fire", "emit",
	},
}

var verifyCmd = &cobra.Command{
	Use:   "verify-subscription [event]",
	Short: "Mocks the subscription verification event. Can be forwarded to a local webserver for testing.",
	Long: fmt.Sprintf(`Mocks the subscription verification event that can be forwarded to a local webserver for testing.
	Supported:
	%s`, events.ValidTriggers()),
	Args:      cobra.MaximumNArgs(1),
	ValidArgs: events.ValidTriggers(),
	Run:       verifyCmdRun,
	Example:   `twitch event verify-subscription subscribe`,
	Aliases: []string{
		"verify",
	},
}

var retriggerCmd = &cobra.Command{
	Use:     "retrigger",
	Short:   "Refires events based on the event ID. Can be forwarded to the local webserver for event testing.",
	Run:     retriggerCmdRun,
	Example: `twitch event retrigger subscribe`,
}

var startWebsocketServerCmd = &cobra.Command{
	Use:     "start-websocket-server",
	Short:   `Starts a local websocket server at "ws://localhost:8080/eventsub" or at another preferred port.`,
	Run:     startWebsocketServerCmdRun,
	Example: `twitch event start-websocket-server`,
	Aliases: []string{
		"ws",
		"wss",
	},
}

func init() {
	rootCmd.AddCommand(eventCmd)
	eventCmd.AddCommand(triggerCmd, retriggerCmd, verifyCmd, startWebsocketServerCmd)

	// trigger flags
	//// flags for forwarding functionality/changing payloads
	triggerCmd.Flags().StringVarP(&forwardAddress, "forward-address", "F", "", "Forward address for mock event.")
	triggerCmd.Flags().StringVarP(&transport, "transport", "T", "eventsub", fmt.Sprintf("Preferred transport method for event. Defaults to /EventSub.\nSupported values: %s", events.ValidTransports()))
	triggerCmd.Flags().StringVarP(&secret, "secret", "s", "", "Webhook secret. If defined, signs all forwarded events with the SHA256 HMAC and must be 10-100 characters in length.")

	// trigger flags
	//// per-topic flags
	triggerCmd.Flags().StringVarP(&toUser, "to-user", "t", "", "User ID of the receiver of the event. For example, the user that receives a follow. In most contexts, this is the broadcaster.")
	triggerCmd.Flags().StringVarP(&fromUser, "from-user", "f", "", "User ID of the user sending the event, for example the user following another user.")
	triggerCmd.Flags().StringVarP(&giftUser, "gift-user", "g", "", "Used only for \"gift\" events. Denotes the User ID of the gifting user.")
	triggerCmd.Flags().BoolVarP(&isAnonymous, "anonymous", "a", false, "Denotes if the event is anonymous. Only applies to Gift and Sub events.")
	triggerCmd.Flags().IntVarP(&count, "count", "c", 1, "Count of events to events. This will simulate a sub gift, or large number of cheers.")
	triggerCmd.Flags().StringVarP(&eventStatus, "event-status", "S", "", "Status of the Event object (.event.status in JSON); currently applies to channel points redemptions.")
	triggerCmd.Flags().StringVarP(&subscriptionStatus, "subscription-status", "r", "enabled", "Status of the Subscription object (.subscription.status in JSON). Defaults to \"enabled\".")
	triggerCmd.Flags().StringVarP(&itemID, "item-id", "i", "", "Manually set the ID of the event payload item (for example the reward ID in redemption events). For stream events, this is the game ID.")
	triggerCmd.Flags().StringVarP(&itemName, "item-name", "n", "", "Manually set the name of the event payload item (for example the reward ID in redemption events). For stream events, this is the game title.")
	triggerCmd.Flags().Int64VarP(&cost, "cost", "C", 0, "Amount of bits or channel points redeemed/used in the event.")
	triggerCmd.Flags().StringVarP(&description, "description", "d", "", "Title the stream should be updated with.")
	triggerCmd.Flags().StringVarP(&gameID, "game-id", "G", "", "Sets the game/category ID for applicable events.")
	triggerCmd.Flags().StringVarP(&eventID, "subscription-id", "u", "", "Manually set the subscription/event ID of the event itself.") // TODO: This description will need to change with https://github.com/twitchdev/twitch-cli/issues/184
	triggerCmd.Flags().StringVarP(&timestamp, "timestamp", "Z", "", "Sets the timestamp to be used in payloads and headers. Must be in RFC3339Nano format.")

	// retrigger flags
	retriggerCmd.Flags().StringVarP(&forwardAddress, "forward-address", "F", "", "Forward address for mock event.")
	retriggerCmd.Flags().StringVarP(&eventID, "id", "i", "", "ID of the event to be refired.")
	retriggerCmd.Flags().StringVarP(&secret, "secret", "s", "", "Webhook secret. If defined, signs all forwarded events with the SHA256 HMAC and must be 10-100 characters in length.")
	retriggerCmd.MarkFlagRequired("id")

	// verify-subscription flags
	verifyCmd.Flags().StringVarP(&forwardAddress, "forward-address", "F", "", "Forward address for mock event.")
	verifyCmd.Flags().StringVarP(&transport, "transport", "T", "eventsub", fmt.Sprintf("Preferred transport method for event. Defaults to EventSub.\nSupported values: %s", events.ValidTransports()))
	verifyCmd.Flags().StringVarP(&secret, "secret", "s", "", "Webhook secret. If defined, signs all forwarded events with the SHA256 HMAC and must be 10-100 characters in length.")
	verifyCmd.Flags().StringVarP(&eventID, "subscription-id", "u", "", "Manually set the subscription/event ID of the event itself.") // TODO: This description will need to change with https://github.com/twitchdev/twitch-cli/issues/184
	verifyCmd.MarkFlagRequired("forward-address")

	// start-websocket-server flags
	startWebsocketServerCmd.Flags().IntVarP(&port, "port", "p", 8080, "Defines the port that the mock EventSub websocket server will run on.")
	startWebsocketServerCmd.Flags().BoolVar(&debug, "debug", false, "Set on/off for debug messages for the EventSub WebSocket server.")
	startWebsocketServerCmd.Flags().BoolVar(&sslEnabled, "ssl", false, "Sets on/off for SSL. Recommended to keep 'false', as most testing does not require this.")
	startWebsocketServerCmd.Flags().IntVarP(&wssReconnectTest, "reconnect", "r", -1, "Used to test WebSocket Reconnect message. Sets delay (in seconds) from first client connection until the reconnect occurs.")
}

func triggerCmdRun(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		cmd.Help()
		return
	}

	if transport == "websub" {
		fmt.Println(websubDeprecationNotice)
		return
	}

	if secret != "" && (len(secret) < 10 || len(secret) > 100) {
		fmt.Println("Invalid secret provided. Secrets must be between 10-100 characters")
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
			Event:              args[0],
			EventID:            eventID,
			Transport:          transport,
			ForwardAddress:     forwardAddress,
			FromUser:           fromUser,
			ToUser:             toUser,
			GiftUser:           giftUser,
			Secret:             secret,
			IsAnonymous:        isAnonymous,
			EventStatus:        eventStatus,
			ItemID:             itemID,
			Cost:               cost,
			Description:        description,
			ItemName:           itemName,
			GameID:             gameID,
			SubscriptionStatus: subscriptionStatus,
			Timestamp:          timestamp,
		})

		if err != nil {
			println(err.Error())
			return
		}

		fmt.Println(res)
	}
}

func retriggerCmdRun(cmd *cobra.Command, args []string) {
	if transport == "websub" {
		fmt.Println(websubDeprecationNotice)
		return
	}

	if secret != "" && (len(secret) < 10 || len(secret) > 100) {
		fmt.Println("Invalid secret provided. Secrets must be between 10-100 characters")
		return
	}

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

func verifyCmdRun(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		cmd.Help()
		return
	}

	if transport == "websub" {
		fmt.Println(websubDeprecationNotice)
		return
	}

	if secret != "" && (len(secret) < 10 || len(secret) > 100) {
		fmt.Println("Invalid secret provided. Secrets must be between 10-100 characters")
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

	_, err := verify.VerifyWebhookSubscription(verify.VerifyParameters{
		Event:          args[0],
		Transport:      transport,
		ForwardAddress: forwardAddress,
		Secret:         secret,
		EventID:        eventID,
	})

	if err != nil {
		println(err.Error())
		return
	}
}

func startWebsocketServerCmdRun(cmd *cobra.Command, args []string) {
	wsStr := "ws"
	if sslEnabled {
		wsStr = "wss"
	}

	log.Printf("Starting mock EventSub WebSocket servers on %v://localhost:%v/eventsub (alternate on port %v)", wsStr, port, port+1)
	mock_wss_server.StartServer(port, debug, wssReconnectTest, sslEnabled)
}

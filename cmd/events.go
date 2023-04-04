// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package cmd

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/spf13/cobra"
	"github.com/twitchdev/twitch-cli/internal/events"
	"github.com/twitchdev/twitch-cli/internal/events/trigger"
	"github.com/twitchdev/twitch-cli/internal/events/types"
	"github.com/twitchdev/twitch-cli/internal/events/verify"
	"github.com/twitchdev/twitch-cli/internal/events/websocket"
	"github.com/twitchdev/twitch-cli/internal/events/websocket/mock_server"
	"github.com/twitchdev/twitch-cli/internal/util"
)

const websubDeprecationNotice = "Halt! It appears you are trying to use WebSub, which has been deprecated. For more information, see: https://discuss.dev.twitch.tv/t/deprecation-of-websub-based-webhooks/32152"

var (
	isAnonymous         bool
	forwardAddress      string
	event               string
	transport           string
	fromUser            string
	toUser              string
	giftUser            string
	eventID             string
	secret              string
	eventStatus         string
	subscriptionStatus  string
	itemID              string
	itemName            string
	cost                int64
	count               int
	description         string
	gameID              string
	tier                string
	timestamp           string
	charityCurrentValue int
	charityTargetValue  int
	clientId            string
	version             string
	websocketClient     string
	websocketServerIP   string
	websocketServerPort int
)

// websocketCmd-specific flags
var (
	wsDebug        bool
	wsStrict       bool
	wsClient       string
	wsSubscription string
	wsStatus       string
	wsReason       string
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
	%s`, types.AllWebhookTopics()),
	Args:      cobra.MaximumNArgs(1),
	ValidArgs: types.AllWebhookTopics(),
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
	%s`, types.AllWebhookTopics()),
	Args:      cobra.MaximumNArgs(1),
	ValidArgs: types.AllWebhookTopics(),
	Run:       verifyCmdRun,
	Example:   `twitch event verify-subscription subscribe`,
	Aliases: []string{
		"verify",
	},
}

var websocketCmd = &cobra.Command{
	Use:   "websocket [action]",
	Short: `Executes actions regarding the mock EventSub WebSocket server. See "twitch event websocket --help" for usage info.`,
	Long:  fmt.Sprintf(`Executes actions regarding the mock EventSub WebSocket server.`),
	Args:  cobra.MaximumNArgs(1),
	Run:   websocketCmdRun,
	Example: fmt.Sprintf(`  twitch event websocket start-server
  twitch event websocket reconnect
  twitch event websocket close --session=e411cc1e_a2613d4e --reason=4006
  twitch event websocket subscription --status=user_removed --subscription=82a855-fae8-93bff0`,
	),
	Aliases: []string{
		"websockets",
		"ws",
		"wss",
	},
}

var retriggerCmd = &cobra.Command{
	Use:     "retrigger",
	Short:   "Refires events based on the event ID. Can be forwarded to the local webserver for event testing.",
	Run:     retriggerCmdRun,
	Example: `twitch event retrigger subscribe`,
}

var startWebsocketServerCmd = &cobra.Command{
	Use:        "start-websocket-server",
	Deprecated: `use "twitch event websocket start-server" instead.`,
}

func init() {
	rootCmd.AddCommand(eventCmd)

	eventCmd.AddCommand(triggerCmd, retriggerCmd, verifyCmd, websocketCmd, startWebsocketServerCmd)

	// trigger flags
	//// flags for forwarding functionality/changing payloads
	triggerCmd.Flags().StringVarP(&forwardAddress, "forward-address", "F", "", "Forward address for mock event (webhook only).")
	triggerCmd.Flags().StringVarP(&transport, "transport", "T", "webhook", fmt.Sprintf("Preferred transport method for event. Defaults to /EventSub.\nSupported values: %s", events.ValidTransports()))
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
	triggerCmd.Flags().StringVarP(&tier, "tier", "", "", "Sets the subscription tier. Valid values are 1000, 2000, and 3000.")
	triggerCmd.Flags().StringVarP(&eventID, "subscription-id", "u", "", "Manually set the subscription/event ID of the event itself.") // TODO: This description will need to change with https://github.com/twitchdev/twitch-cli/issues/184
	triggerCmd.Flags().StringVar(&timestamp, "timestamp", "", "Sets the timestamp to be used in payloads and headers. Must be in RFC3339Nano format.")
	triggerCmd.Flags().IntVar(&charityCurrentValue, "charity-current-value", 0, "Only used for \"charity-*\" events. Manually set the current dollar value for charity events.")
	triggerCmd.Flags().IntVar(&charityTargetValue, "charity-target-value", 1500000, "Only used for \"charity-*\" events. Manually set the target dollar value for charity events.")
	triggerCmd.Flags().StringVar(&clientId, "client-id", "", "Manually set the Client ID used in revoke, grant, and bits transaction events.")
	triggerCmd.Flags().StringVarP(&version, "version", "v", "", "Chooses the EventSub version used for a specific event. Not required for most events.")
	triggerCmd.Flags().StringVar(&websocketClient, "session", "", "Defines a specific websocket client/session to forward an event to. Used only with \"websocket\" transport.")

	// retrigger flags
	retriggerCmd.Flags().StringVarP(&forwardAddress, "forward-address", "F", "", "Forward address for mock event (webhook only).")
	retriggerCmd.Flags().StringVarP(&eventID, "id", "i", "", "ID of the event to be refired.")
	retriggerCmd.Flags().StringVarP(&secret, "secret", "s", "", "Webhook secret. If defined, signs all forwarded events with the SHA256 HMAC and must be 10-100 characters in length.")
	retriggerCmd.MarkFlagRequired("id")

	// verify-subscription flags
	verifyCmd.Flags().StringVarP(&forwardAddress, "forward-address", "F", "", "Forward address for mock event (webhook only).")
	verifyCmd.Flags().StringVarP(&transport, "transport", "T", "webhook", fmt.Sprintf("Preferred transport method for event. Defaults to EventSub.\nSupported values: %s", events.ValidTransports()))
	verifyCmd.Flags().StringVarP(&secret, "secret", "s", "", "Webhook secret. If defined, signs all forwarded events with the SHA256 HMAC and must be 10-100 characters in length.")
	verifyCmd.Flags().StringVar(&timestamp, "timestamp", "", "Sets the timestamp to be used in payloads and headers. Must be in RFC3339Nano format.")
	verifyCmd.Flags().StringVarP(&eventID, "subscription-id", "u", "", "Manually set the subscription/event ID of the event itself.") // TODO: This description will need to change with https://github.com/twitchdev/twitch-cli/issues/184
	verifyCmd.Flags().StringVarP(&version, "version", "v", "", "Chooses the EventSub version used for a specific event. Not required for most events.")
	verifyCmd.MarkFlagRequired("forward-address")

	// websocket flags
	/// flags for start-server
	websocketCmd.Flags().StringVar(&websocketServerIP, "ip", "127.0.0.1", "Defines the ip that the mock EventSub websocket server will bind to.")
	websocketCmd.Flags().IntVarP(&websocketServerPort, "port", "p", 8080, "Defines the port that the mock EventSub websocket server will run on.")
	websocketCmd.Flags().BoolVar(&wsDebug, "debug", false, "Set on/off for debug messages for the EventSub WebSocket server.")
	websocketCmd.Flags().BoolVarP(&wsStrict, "require-subscription", "S", false, "Requires subscriptions for all events, and activates 10 second subscription requirement.")

	// websocket flags
	/// flags for everything else
	websocketCmd.Flags().StringVarP(&wsClient, "session", "s", "", "WebSocket client/session to target with your server command. Used in multiple commands.")
	websocketCmd.Flags().StringVar(&wsSubscription, "subscription", "", `Subscription to target with your server command. Used with "websocket subscription".`)
	websocketCmd.Flags().StringVar(&wsStatus, "status", "", `Changes the status of an existing subscription. Used with "websocket subscription".`)
	websocketCmd.Flags().StringVar(&wsReason, "reason", "", `Sets the close reason when sending a Close message to the client. Used with "websocket close".`)
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
			Event:               args[0],
			EventID:             eventID,
			Transport:           transport,
			ForwardAddress:      forwardAddress,
			FromUser:            fromUser,
			ToUser:              toUser,
			GiftUser:            giftUser,
			Secret:              secret,
			IsAnonymous:         isAnonymous,
			EventStatus:         eventStatus,
			ItemID:              itemID,
			Cost:                cost,
			Description:         description,
			ItemName:            itemName,
			GameID:              gameID,
			Tier:                tier,
			SubscriptionStatus:  subscriptionStatus,
			Timestamp:           timestamp,
			CharityCurrentValue: charityCurrentValue,
			CharityTargetValue:  charityTargetValue,
			ClientID:            clientId,
			Version:             version,
			WebSocketClient:     websocketClient,
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
		Timestamp:      util.GetTimestamp().Format(time.RFC3339Nano),
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

	if timestamp == "" {
		timestamp = util.GetTimestamp().Format(time.RFC3339Nano)
	} else {
		// Verify custom timestamp
		_, err := time.Parse(time.RFC3339Nano, timestamp)
		if err != nil {
			fmt.Println(
				`Discarding verify: Invalid timestamp provided.
Please follow RFC3339Nano, which is used by Twitch as seen here:
https://dev.twitch.tv/docs/eventsub/handling-webhook-events#processing-an-event`)
			return
		}
	}

	_, err := verify.VerifyWebhookSubscription(verify.VerifyParameters{
		Event:          args[0],
		Transport:      transport,
		ForwardAddress: forwardAddress,
		Secret:         secret,
		Timestamp:      timestamp,
		EventID:        eventID,
	})

	if err != nil {
		println(err.Error())
		return
	}
}

func websocketCmdRun(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		cmd.Help()
		return
	}

	if args[0] == "start-server" || args[0] == "start" {
		log.Printf("`Ctrl + C` to exit mock WebSocket servers.")
		mock_server.StartWebsocketServer(wsDebug, websocketServerIP, websocketServerPort, wsStrict)
	} else {
		// Forward all other commands via RPC
		websocket.ForwardWebsocketCommand(args[0], websocket.WebsocketCommandParameters{
			Client:             wsClient,
			Subscription:       wsSubscription,
			SubscriptionStatus: wsStatus,
			CloseReason:        wsReason,
		})
	}
}

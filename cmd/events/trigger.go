package events

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
	"github.com/twitchdev/twitch-cli/internal/events"
	configure_event "github.com/twitchdev/twitch-cli/internal/events/configure"
	"github.com/twitchdev/twitch-cli/internal/events/trigger"
	"github.com/twitchdev/twitch-cli/internal/events/types"
)

func TriggerCommand() (command *cobra.Command) {
	command = &cobra.Command{
		Use:   "trigger [event]",
		Short: "Creates mock events that can be forwarded to a local webserver for event testing.",
		Long: fmt.Sprintf(`Creates mock events that can be forwarded to a local webserver for event testing.
		Supported:
		%s`, types.AllWebhookTopics()),
		Args:      cobra.MaximumNArgs(1),
		ValidArgs: types.AllWebhookTopics(),
		RunE:      triggerCmdRun,
		Example:   `twitch event trigger subscribe`,
		Aliases: []string{
			"fire", "emit",
		},
	}

	// flags for forwarding functionality/changing payloads
	command.Flags().StringVarP(&forwardAddress, "forward-address", "F", "", "Forward address for mock event (webhook only).")
	command.Flags().StringVarP(&transport, "transport", "T", "webhook", fmt.Sprintf("Preferred transport method for event. Defaults to /EventSub.\nSupported values: %s", events.ValidTransports()))
	command.Flags().StringVarP(&secret, "secret", "s", "", "Webhook secret. If defined, signs all forwarded events with the SHA256 HMAC and must be 10-100 characters in length.")
	command.Flags().BoolVarP(&noConfig, "no-config", "D", false, "Disables the use of the configuration, if it exists.")

	// per-topic flags
	command.Flags().StringVarP(&toUser, "to-user", "t", "", "User ID of the receiver of the event. For example, the user that receives a follow. In most contexts, this is the broadcaster.")
	command.Flags().StringVarP(&fromUser, "from-user", "f", "", "User ID of the user sending the event, for example the user following another user.")
	command.Flags().StringVarP(&giftUser, "gift-user", "g", "", "Used only for \"gift\" events. Denotes the User ID of the gifting user.")
	command.Flags().BoolVarP(&isAnonymous, "anonymous", "a", false, "Denotes if the event is anonymous. Only applies to Gift and Sub events.")
	command.Flags().IntVarP(&count, "count", "c", 1, "Number of times to run an event. This can be used to simulate rapid events, such as multiple sub gift, or large number of cheers.")
	command.Flags().StringVarP(&eventStatus, "event-status", "S", "", "Status of the Event object (.event.status in JSON); currently applies to channel points redemptions.")
	command.Flags().StringVarP(&subscriptionStatus, "subscription-status", "r", "enabled", "Status of the Subscription object (.subscription.status in JSON). Defaults to \"enabled\".")
	command.Flags().StringVarP(&itemID, "item-id", "i", "", "Manually set the ID of the event payload item (for example the reward ID in redemption events). For stream events, this is the game ID.")
	command.Flags().StringVarP(&itemName, "item-name", "n", "", "Manually set the name of the event payload item (for example the reward ID in redemption events). For stream events, this is the game title.")
	command.Flags().Int64VarP(&cost, "cost", "C", 0, "Amount of drops, subscriptions, bits, or channel points redeemed/used in the event.")
	command.Flags().StringVarP(&description, "description", "d", "", "Title the stream should be updated with.")
	command.Flags().StringVarP(&gameID, "game-id", "G", "", "Sets the game/category ID for applicable events.")
	command.Flags().StringVarP(&tier, "tier", "", "", "Sets the subscription tier. Valid values are 1000, 2000, and 3000.")
	command.Flags().StringVarP(&subscriptionID, "subscription-id", "u", "", "Manually set the subscription/event ID of the event itself.")
	command.Flags().StringVarP(&eventMessageID, "event-id", "I", "", "Manually set the Twitch-Eventsub-Message-Id header value for the event.")
	command.Flags().StringVar(&timestamp, "timestamp", "", "Sets the timestamp to be used in payloads and headers. Must be in RFC3339Nano format.")
	command.Flags().IntVar(&charityCurrentValue, "charity-current-value", 0, "Only used for \"charity-*\" events. Manually set the current dollar value for charity events.")
	command.Flags().IntVar(&charityTargetValue, "charity-target-value", 1500000, "Only used for \"charity-*\" events. Manually set the target dollar value for charity events.")
	command.Flags().StringVar(&clientId, "client-id", "", "Manually set the Client ID used in revoke, grant, and bits transaction events.")
	command.Flags().StringVarP(&version, "version", "v", "", "Chooses the EventSub version used for a specific event. Not required for most events.")
	command.Flags().StringVar(&websocketClient, "session", "", "Defines a specific websocket client/session to forward an event to. Used only with \"websocket\" transport.")
	command.Flags().StringVar(&banStart, "ban-start", "", "Sets the timestamp a ban started at.")
	command.Flags().StringVar(&banEnd, "ban-end", "", "Sets the timestamp a ban is intended to end at. If not set, the ban event will appear as permanent. This flag can take a timestamp or relative time (600, 600s, 10d4h12m55s)")

	return
}

func triggerCmdRun(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		cmd.Help()
		return fmt.Errorf("")
	}

	if transport == "websub" {
		return fmt.Errorf(websubDeprecationNotice)
	}

	defaults := configure_event.GetEventConfiguration(noConfig)

	if secret != "" {
		if len(secret) < 10 || len(secret) > 100 {
			return fmt.Errorf("Invalid secret provided. Secrets must be between 10-100 characters")
		}
	} else {
		secret = defaults.Secret
	}

	// Validate that the forward address is actually a URL
	if len(forwardAddress) > 0 {
		_, err := url.ParseRequestURI(forwardAddress)
		if err != nil {
			return err
		}
	} else {
		forwardAddress = defaults.ForwardAddress
	}

	for i := 0; i < count; i++ {
		res, err := trigger.Fire(trigger.TriggerParameters{
			Event:               args[0],
			SubscriptionID:      subscriptionID,
			EventMessageID:      eventMessageID,
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
			BanStartTimestamp:   banStart,
			BanEndTimestamp:     banEnd,
		})

		if err != nil {
			return err
		}

		fmt.Println(res)
	}

	return nil
}

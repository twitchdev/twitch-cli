package events

import (
	"fmt"
	"net/url"
	"time"

	"github.com/spf13/cobra"
	"github.com/twitchdev/twitch-cli/internal/events"
	configure_event "github.com/twitchdev/twitch-cli/internal/events/configure"
	"github.com/twitchdev/twitch-cli/internal/events/types"
	"github.com/twitchdev/twitch-cli/internal/events/verify"
	"github.com/twitchdev/twitch-cli/internal/util"
)

func VerifySubscriptionCommand() (command *cobra.Command) {
	command = &cobra.Command{
		Use:   "verify-subscription [event]",
		Short: "Mocks the subscription verification event. Can be forwarded to a local webserver for testing.",
		Long: fmt.Sprintf(`Mocks the subscription verification event that can be forwarded to a local webserver for testing.
		Supported:
		%s`, types.AllWebhookTopics()),
		Args:      cobra.MaximumNArgs(1),
		ValidArgs: types.AllWebhookTopics(),
		RunE:      verifyCmdRun,
		Example:   `twitch event verify-subscription subscribe`,
		Aliases: []string{
			"verify",
		},
	}

	command.Flags().StringVarP(&forwardAddress, "forward-address", "F", "", "Forward address for mock event (webhook only).")
	command.Flags().StringVarP(&transport, "transport", "T", "webhook", fmt.Sprintf("Preferred transport method for event. Defaults to EventSub.\nSupported values: %s", events.ValidTransports()))
	command.Flags().StringVarP(&secret, "secret", "s", "", "Webhook secret. If defined, signs all forwarded events with the SHA256 HMAC and must be 10-100 characters in length.")
	command.Flags().StringVar(&timestamp, "timestamp", "", "Sets the timestamp to be used in payloads and headers. Must be in RFC3339Nano format.")
	command.Flags().StringVarP(&eventID, "subscription-id", "u", "", "Manually set the subscription/event ID of the event itself.") // TODO: This description will need to change with https://github.com/twitchdev/twitch-cli/issues/184
	command.Flags().StringVarP(&version, "version", "v", "", "Chooses the EventSub version used for a specific event. Not required for most events.")
	command.Flags().BoolVarP(&noConfig, "no-config", "D", false, "Disables the use of the configuration, if it exists.")
	command.Flags().StringVarP(&toUser, "broadcaster", "b", "", "User ID of the broadcaster for the verification event.")

	return
}

func verifyCmdRun(cmd *cobra.Command, args []string) error {
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

	if timestamp == "" {
		timestamp = util.GetTimestamp().Format(time.RFC3339Nano)
	} else {
		// Verify custom timestamp
		_, err := time.Parse(time.RFC3339Nano, timestamp)
		if err != nil {
			return fmt.Errorf(
				`Discarding verify: Invalid timestamp provided.
Please follow RFC3339Nano, which is used by Twitch as seen here:
https://dev.twitch.tv/docs/eventsub/handling-webhook-events#processing-an-event`)
		}
	}

	_, err := verify.VerifyWebhookSubscription(verify.VerifyParameters{
		Event:             args[0],
		Transport:         transport,
		ForwardAddress:    forwardAddress,
		Secret:            secret,
		Timestamp:         timestamp,
		EventID:           eventID,
		BroadcasterUserID: toUser,
		Version:           version,
	})

	if err != nil {
		return err
	}

	return nil
}

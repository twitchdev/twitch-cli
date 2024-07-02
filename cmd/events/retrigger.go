package events

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	configure_event "github.com/twitchdev/twitch-cli/internal/events/configure"
	"github.com/twitchdev/twitch-cli/internal/events/trigger"
	"github.com/twitchdev/twitch-cli/internal/util"
)

func RetriggerCommand() (command *cobra.Command) {
	command = &cobra.Command{
		Use:     "retrigger",
		Short:   "Refires events based on the event ID. Can be forwarded to the local webserver for event testing.",
		RunE:    retriggerCmdRun,
		Example: `twitch event retrigger subscribe`,
	}

	command.Flags().StringVarP(&forwardAddress, "forward-address", "F", "", "Forward address for mock event (webhook only).")
	command.Flags().StringVarP(&eventMessageID, "id", "i", "", "ID of the event to be refired.")
	command.Flags().StringVarP(&secret, "secret", "s", "", "Webhook secret. If defined, signs all forwarded events with the SHA256 HMAC and must be 10-100 characters in length.")
	command.Flags().BoolVarP(&noConfig, "no-config", "D", false, "Disables the use of the configuration, if it exists.")
	command.MarkFlagRequired("id")

	return
}

func retriggerCmdRun(cmd *cobra.Command, args []string) error {
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

	if forwardAddress == "" {
		if defaults.ForwardAddress == "" {
			return fmt.Errorf("if a default configuration is not set, forward-address must be provided")
		}
		forwardAddress = defaults.ForwardAddress
	}

	//color.New().Add(color.FgGreen).Println(fmt.Sprintf(`Refire %v`, eventMessageID));
	res, err := trigger.RefireEvent(eventMessageID, trigger.TriggerParameters{
		ForwardAddress: forwardAddress,
		Secret:         secret,
		Timestamp:      util.GetTimestamp().Format(time.RFC3339Nano),
	})
	if err != nil {
		return fmt.Errorf("Error refiring event: %s", err)
	}

	fmt.Println(res)
	return nil
}

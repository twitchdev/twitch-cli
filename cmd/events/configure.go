package events

import (
	"github.com/spf13/cobra"
	configure_event "github.com/twitchdev/twitch-cli/internal/events/configure"
)

func ConfigureCommand() (command *cobra.Command) {
	command = &cobra.Command{
		Use:     "configure",
		Short:   "Allows users to configure defaults for the twitch event subcommands.",
		RunE:    configureEventRun,
		Example: `twitch event configure`,
	}

	command.Flags().StringVarP(&forwardAddress, "forward-address", "F", "", "Forward address for mock event (webhook only).")
	command.Flags().StringVarP(&secret, "secret", "s", "", "Webhook secret. If defined, signs all forwarded events with the SHA256 HMAC and must be 10-100 characters in length.")

	return
}

func configureEventRun(cmd *cobra.Command, args []string) error {
	return configure_event.ConfigureEvents(configure_event.EventConfigurationParams{
		ForwardAddress: forwardAddress,
		Secret:         secret,
	})
}

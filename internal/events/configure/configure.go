package configure_event

import (
	"fmt"
	"net/url"

	"github.com/spf13/viper"
	"github.com/twitchdev/twitch-cli/internal/util"
)

type EventConfigurationParams struct {
	Secret         string
	ForwardAddress string
}

func ConfigureEvents(p EventConfigurationParams) error {
	var err error
	if p.ForwardAddress == "" && p.Secret == "" {
		return fmt.Errorf("you must provide at least one of --secret or --forward-address")
	}

	// Validate that the forward address is actually a URL
	if len(p.ForwardAddress) > 0 {
		_, err := url.ParseRequestURI(p.ForwardAddress)
		if err != nil {
			return err
		}
		viper.Set("forwardAddress", p.ForwardAddress)
	}
	if p.Secret != "" {
		if len(p.Secret) < 10 || len(p.Secret) > 100 {
			return fmt.Errorf("invalid secret provided. Secrets must be between 10-100 characters")
		}
		viper.Set("eventSecret", p.Secret)
	}

	configPath, err := util.GetConfigPath()
	if err != nil {
		return err
	}

	if err := viper.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("failed to write configuration: %v", err.Error())
	}

	fmt.Println("Updated configuration.")
	return nil
}

func GetEventConfiguration(noConfig bool) EventConfigurationParams {
	if noConfig {
		return EventConfigurationParams{}
	}
	return EventConfigurationParams{
		ForwardAddress: viper.GetString("forwardAddress"),
		Secret:         viper.GetString("eventSecret"),
	}
}

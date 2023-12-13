package configure_event_test

import (
	"testing"

	"github.com/spf13/viper"
	configure_event "github.com/twitchdev/twitch-cli/internal/events/configure"
	"github.com/twitchdev/twitch-cli/test_setup"
)

func TestWriteEventConfig(t *testing.T) {
	a := test_setup.SetupTestEnv(t)
	defaultForwardAddress := "http://localhost:3000/"
	defaultSecret := "12345678910"
	test_config := configure_event.EventConfigurationParams{
		ForwardAddress: defaultForwardAddress,
		Secret:         defaultSecret,
	}

	// test a good config writes correctly
	a.NoError(configure_event.ConfigureEvents(test_config))

	a.Equal(defaultForwardAddress, viper.Get("forwardAddress"))
	a.Equal(defaultSecret, viper.Get("eventSecret"))

	// test for secret length validation
	test_config.Secret = "1"
	a.Error(configure_event.ConfigureEvents(test_config))
	a.NotEqual("1", viper.Get("eventSecret"))
	test_config.Secret = defaultSecret

	// test for forward address validation
	test_config.ForwardAddress = "not a url"
	a.Error(configure_event.ConfigureEvents(test_config))
	a.NotEqual("not a url", viper.Get("forwardAddress"))
}

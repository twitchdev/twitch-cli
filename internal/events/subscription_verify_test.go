package events

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
)

func TestSubscriptionVerify(t *testing.T) {
	a := util.SetupTestEnv(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var challenge string
		w.WriteHeader(http.StatusAccepted)

		body, err := ioutil.ReadAll(r.Body)
		a.Nil(err)

		if r.Method == http.MethodPost {
			var verification models.EventsubSubscriptionVerification
			err := json.Unmarshal(body, &verification)
			a.Nil(err)

			challenge = verification.Challenge
		} else if r.Method == http.MethodGet {
			q := r.URL.Query()

			challenge = q.Get("hub.challenge")
		}

		w.Header().Add("Content-Type", "text/plain")
		w.Write([]byte(challenge))
	}))
	defer ts.Close()

	p := VerifyParameters{
		Transport:      TransportEventSub,
		Event:          "subscribe",
		ForwardAddress: ts.URL,
		Secret:         "potatoes",
	}
	res, err := VerifyWebhookSubscription(p)
	a.Nil(err)
	a.Equal(res.IsValid, true)

	p.Transport = TransportWebSub
	res, err = VerifyWebhookSubscription(p)
	a.Nil(err)
	a.Equal(res.IsValid, true)

	p.Event = "cheer"
	_, err = VerifyWebhookSubscription(p)
	a.NotNil(err)

	p.Transport = "notarealtransport"
	_, err = VerifyWebhookSubscription(p)
	a.NotNil(err)

}

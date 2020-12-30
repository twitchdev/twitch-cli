package request

import (
	"strings"
	"testing"
)

func TestNewRequest(t *testing.T) {
	r, err := NewRequest("GET", "https://api.twitch.tv/helix/users", nil)
	if err != nil {
		t.Errorf("Received error %v for valid request", err)
	}

	if !strings.Contains(r.Header.Get("User-Agent"), "twitch-cli/") {
		t.Error("User agent not properly set")
	}
}

package request

import (
	"testing"
)

func TestNewRequest(t *testing.T) {
	_, err := NewRequest("GET", "https://api.twitch.tv/helix/users", nil)
	if err != nil {
		t.Errorf("Received error %v for valid request", err)
	}
}

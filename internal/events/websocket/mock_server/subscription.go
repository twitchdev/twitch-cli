package mock_server

import "github.com/twitchdev/twitch-cli/internal/models"

type Subscription struct {
	SubscriptionID    string // Random GUID for the subscription
	ClientID          string // Client ID included in headers
	Type              string // EventSub topic
	Version           string // EventSub topic version
	CreatedAt         string // Timestamp of when the subscription was created
	Status            string // Status of the subscription
	SessionClientName string // Client name of the session this is associated with.

	ClientConnectedAt    string // Time client connected
	ClientDisconnectedAt string // Time client disconnected

	Conditions models.EventsubCondition // Values of the subscription's condition object
}

// Request - POST /eventsub/subscriptions
type SubscriptionPostRequest struct {
	Type    string `json:"type"`
	Version string `json:"version"`

	Condition models.EventsubCondition         `json:"condition"`
	Transport SubscriptionPostRequestTransport `json:"transport"`
}

// Request - POST /eventsub/subscriptions
type SubscriptionPostRequestTransport struct {
	Method    string `json:"method"`
	SessionID string `json:"session_id"`
}

// Response (Success) - POST /eventsub/subscriptions
type SubscriptionPostSuccessResponse struct {
	Data []SubscriptionPostSuccessResponseBody `json:"data"`

	Total        int `json:"total"`
	MaxTotalCost int `json:"max_total_cost"`
	TotalCost    int `json:"total_cost"`
}

// Response (Success) - POST /eventsub/subscriptions
// Response (Success) - GET /eventsub/subscriptions
type SubscriptionPostSuccessResponseBody struct {
	ID        string `json:"id"`
	Status    string `json:"status"`
	Type      string `json:"type"`
	Version   string `json:"version"`
	CreatedAt string `json:"created_at"`
	Cost      int    `json:"cost"`

	Condition models.EventsubCondition `json:"condition"`
	Transport SubscriptionTransport    `json:"transport"`
}

// Response (Error) - POST /eventsub/subscriptions
type SubscriptionPostErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// Response (Success) - GET /eventsub/subscriptions
type SubscriptionGetSuccessResponse struct {
	Total        int         `json:"total"`
	TotalCost    int         `json:"total_cost"`
	MaxTotalCost int         `json:"max_total_cost"`
	Pagination   EmptyStruct `json:"pagination"`

	Data []SubscriptionPostSuccessResponseBody `json:"data"`
}

// Cross-usage
type SubscriptionTransport struct {
	Method         string `json:"method"`
	SessionID      string `json:"session_id"`
	ConnectedAt    string `json:"connected_at,omitempty"`
	DisconnectedAt string `json:"disconnected_at,omitempty"`
}

// Cross-usage
type EmptyStruct struct {
}

// Subscription Statuses
// Only includes status values that apply to WebSocket connections
// https://dev.twitch.tv/docs/api/reference/#get-eventsub-subscriptions
const (
	STATUS_ENABLED                            = "enabled"
	STATUS_AUTHORIZATION_REVOKED              = "revoked"
	STATUS_MODERATOR_REMOVED                  = "moderator_removed"
	STATUS_USER_REMOVED                       = "user_removed"
	STATUS_VERSION_REMOVED                    = "version_removed"
	STATUS_WEBSOCKET_DISCONNECTED             = "websocket_disconnected"
	STATUS_WEBSOCKET_FAILED_PING_PONG         = "websocket_failed_ping_pong"
	STATUS_WEBSOCKET_RECEIVED_INBOUND_TRAFFIC = "websocket_received_inbound_traffic"
	STATUS_WEBSOCKET_CONNECTION_UNUSED        = "websocket_connection_unused"
	STATUS_INTERNAL_ERROR                     = "websocket_internal_error"
	STATUS_NETWORK_TIMEOUT                    = "network_timeout"
	STATUS_NETWORK_ERROR                      = "websocket_network_error"
)

func IsValidSubscriptionStatus(status string) bool {
	switch status {
	case STATUS_ENABLED, STATUS_AUTHORIZATION_REVOKED,
		STATUS_MODERATOR_REMOVED, STATUS_USER_REMOVED,
		STATUS_VERSION_REMOVED, STATUS_WEBSOCKET_DISCONNECTED,
		STATUS_WEBSOCKET_FAILED_PING_PONG, STATUS_WEBSOCKET_RECEIVED_INBOUND_TRAFFIC,
		STATUS_WEBSOCKET_CONNECTION_UNUSED, STATUS_INTERNAL_ERROR,
		STATUS_NETWORK_TIMEOUT, STATUS_NETWORK_ERROR:
		return true
	default:
		return false
	}
}

func getStatusFromCloseMessage(reason *CloseMessage) string {
	code := reason.code

	switch code {
	case 1000:
		return STATUS_WEBSOCKET_DISCONNECTED
	case 4000:
		return STATUS_INTERNAL_ERROR
	case 4001:
		return STATUS_WEBSOCKET_RECEIVED_INBOUND_TRAFFIC
	case 4002:
		return STATUS_WEBSOCKET_FAILED_PING_PONG
	case 4003:
		return STATUS_WEBSOCKET_CONNECTION_UNUSED
	case 4004: // grace time expired. Subscriptions stay open
		return STATUS_ENABLED
	case 4005:
		return STATUS_NETWORK_TIMEOUT
	case 4006:
		return STATUS_NETWORK_ERROR
	default:
		return STATUS_WEBSOCKET_DISCONNECTED
	}
}

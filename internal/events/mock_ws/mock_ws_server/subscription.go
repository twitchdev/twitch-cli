package mock_ws_server

type Subscription struct {
	SubscriptionID string // Random GUID for the subscription
	ClientID       string // Client ID included in headers
	Type           string // EventSub topic
	Version        string // EventSub topic version
	CreatedAt      string // Timestamp of when the subscription was created
}

// Request - POST /eventsub/subscriptions
type SubscriptionPostRequest struct {
	Type      string      `json:"type"`
	Version   string      `json:"version"`
	Condition interface{} `json:"condition"`

	Transport SubscriptionPostRequestTransport `json:"transport"`
}

// Request - POST /eventsub/subscriptions
type SubscriptionPostRequestTransport struct {
	Method    string `json:"method"`
	SessionID string `json:"session_id"`
}

// Response (Success) - POST /eventsub/subscriptions
type SubscriptionPostSuccessResponse struct {
	Body SubscriptionPostSuccessResponseBody `json:"body"`

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

	Condition EmptyStruct           `json:"condition"`
	Transport SubscriptionTransport `json:"transport"`
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
	Method      string `json:"method"`
	SessionID   string `json:"session_id"`
	ConnectedAt string `json:"connected_at"`
}

// Cross-usage
type EmptyStruct struct {
}

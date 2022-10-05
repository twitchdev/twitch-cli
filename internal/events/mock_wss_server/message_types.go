package mock_wss_server

// Generic response message Metadata; Always the same

type MessageMetadata struct {
	MessageID        string `json:"message_id"`
	MessageType      string `json:"message_type"`
	MessageTimestamp string `json:"message_timestamp"`
}

// Welcome message

type WelcomeMessagePayload struct {
	Websocket WelcomeMessagePayloadWebsocket `json:"websocket"`
}

type WelcomeMessagePayloadWebsocket struct {
	ID                             string `json:"id"`
	Status                         string `json:"status"`
	MinimumMessageFrequencySeconds int    `json:"minimum_message_frequency_seconds"`
	ConnectedAt                    string `json:"connected_at"`
}

type WelcomeMessage struct {
	Metadata MessageMetadata       `json:"metadata"`
	Payload  WelcomeMessagePayload `json:"payload"`
}

// Reconnect message

type ReconnectMessagePayload struct {
	Websocket ReconnectMessagePayloadWebsocket `json:"websocket"`
}

type ReconnectMessagePayloadWebsocket struct {
	ID                             string `json:"id"`
	Status                         string `json:"status"`
	MinimumMessageFrequencySeconds int    `json:"minimum_message_frequency_seconds"`
	Url                            string `json:"url"`
	ConnectedAt                    string `json:"connected_at"`
	ReconnectingAt                 string `json:"reconnecting_at"`
}

type ReconnectMessage struct {
	Metadata MessageMetadata         `json:"metadata"`
	Payload  ReconnectMessagePayload `json:"payload"`
}

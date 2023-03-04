package mock_ws_server

type CloseMessage struct {
	code    int
	message string
}

var (
	closeInternalServerError = &CloseMessage{
		code:    4000,
		message: "internal server error",
	}

	closeClientSentInboundTraffic = &CloseMessage{
		code:    4001,
		message: "sent inbound traffic",
	}

	closeClientFailedPingPong = &CloseMessage{
		code:    4002,
		message: "failed ping pong",
	}

	closeConnectionUnused = &CloseMessage{
		code:    4003,
		message: "connection unused",
	}

	closeReconnectGraceTimeExpired = &CloseMessage{
		code:    4004,
		message: "client reconnect grace time expired",
	}

	closeNetworkTimeout = &CloseMessage{
		code:    4005,
		message: "network timeout",
	}

	closeNetworkError = &CloseMessage{
		code:    4006,
		message: "network error",
	}

	closeInvalidReconnect = &CloseMessage{
		code:    4007,
		message: "invalid reconnect attempt",
	}
)

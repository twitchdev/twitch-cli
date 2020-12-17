version = "0.1.0"

build:
	xgo -out build/twitch --targets "darwin/amd64,windows/amd64,linux/amd64" --ldflags "-X main.buildVersion=$(version)" ./


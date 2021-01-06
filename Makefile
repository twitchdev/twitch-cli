version = "0.2.0"

release:
	xgo -out build/twitch --targets "darwin/amd64,windows/amd64,linux/amd64" --ldflags "-X main.buildVersion=$(version)" ./

build:
	go build --ldflags "-X main.buildVersion=${version}"
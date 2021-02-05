version = "0.3.1"

release:
	docker build https://github.com/mailchain/goreleaser-xcgo.git --build-arg GORELEASER_VERSION=0.155.0 --build-arg GORELEASER_SHA=2a33aa15933cfd5bd2b714860c4876fa76f1fab8f46a7c6d29a8e32c7f9445f2 -t twitch-cli:latest
	docker run --rm --privileged \
		-v $$PWD:/go/src/github.com/twitchdev/twitch-cli \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-w /go/src/github.com/twitchdev/twitch-cli \
		-e GITHUB_TOKEN=${GITHUB_TOKEN} \
		twitch-cli:latest --rm-dist 

build:
	go build --ldflags "-X main.buildVersion=${version}"

build_all:
	xgo -out build/twitch --targets "darwin/amd64,windows/amd64,linux/amd64" --ldflags "-X main.buildVersion=$(version)" ./
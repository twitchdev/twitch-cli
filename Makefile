version = "0.2.2"

release:
	docker run --rm --privileged \
		-v $$PWD:/go/src/github.com/twitchdev/twitch-cli \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-w /go/src/github.com/twitchdev/twitch-cli \
		mailchain/goreleaser-xcgo --rm-dist 

build:
	go build --ldflags "-X main.buildVersion=${version}"

build_all:
	xgo -out build/twitch --targets "darwin/amd64,windows/amd64,linux/amd64" --ldflags "-X main.buildVersion=$(version)" ./
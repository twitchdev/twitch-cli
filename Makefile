version = "0.3.2"

release:
	docker build . -t twitch-cli:latest
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
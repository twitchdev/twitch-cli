release:
	docker build . -t twitch-cli:latest
	docker run --rm --privileged \
		-v $$PWD:/go/src/github.com/twitchdev/twitch-cli \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-w /go/src/github.com/twitchdev/twitch-cli \
		-e GITHUB_TOKEN=${GITHUB_TOKEN} \
		twitch-cli:latest --rm-dist

test-release:
	docker build . -t twitch-cli:latest
	docker run --rm --privileged \
		-v $$PWD:/go/src/github.com/twitchdev/twitch-cli \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-w /go/src/github.com/twitchdev/twitch-cli \
		-e GITHUB_TOKEN=${GITHUB_TOKEN} \
		twitch-cli:latest --rm-dist --skip-publish --snapshot
	
build:
	go build --ldflags "-s -w -X main.buildVersion=source"

build_all:
	xgo -out build/twitch --targets "darwin/amd64,windows/amd64,linux/amd64" --ldflags "-s -w -X main.buildVersion=source" ./

clean: 
	rm -rf ~/.twitch-cli/eventCache.db
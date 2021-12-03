FROM karalabe/xgo-latest
RUN apt-get update && apt-get install curl -y
RUN curl -sfL https://install.goreleaser.com/github.com/goreleaser/goreleaser.sh | sh
ENV GOPROXY=direct
ENTRYPOINT ["goreleaser"]
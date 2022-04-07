FROM techknowlogick/xgo:latest
RUN echo 'deb [trusted=yes] https://repo.goreleaser.com/apt/ /' | tee /etc/apt/sources.list.d/goreleaser.list
RUN apt-get update && apt-get install goreleaser -y

ENTRYPOINT ["goreleaser"]
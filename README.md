# Twitch CLI (Open Beta Version)

- [Twitch CLI (Open Beta Version)](#twitch-cli-open-beta-version)
  - [Download](#download)
    - [Homebrew](#homebrew)
    - [Scoop](#scoop)
    - [Manual Download](#manual-download)
  - [Usage](#usage)
  - [Commands](#commands)
  - [Contributing](#contributing)
  - [License](#license)

## Download

There are two options to download/install the Twitch CLI for each platform. 

### Homebrew

If you are using MacOS or Linux, we recommend using [Homebrew](https://brew.sh/) for installing the CLI as it will also manage the versioning for you. 

To install via Homebrew, run `brew install twitchdev/twitch/twitch-cli` and it'll be callable via `twitch`. 

### Scoop

If you are using Windows, we recommend using [Scoop](https://scoop.sh/) for installing the CLI, as it'll also manage versioning. 

To install via Scoop, run: 

```sh
scoop bucket add twitch https://github.com/twitchdev/scoop-bucket.git
scoop install twitch-cli
```

This will install it into your path, and it'll be callable via `twitch`. 

### Manual Download

To download, go to the [Releases tab of GitHub](https://github.com/twitchdev/twitch-cli/releases). The examples in the documentation assume you have put this into your PATH and renamed to `twitch` (or symlinked as such).

**Note**: If using MacOS and downloading manually, you may need to adjust the permissions of the file to allow for execution.

To do so, please run: `chmod 755 <filename>` where the filename is the name of the downloaded binary. 

## Usage

The CLI largely follows a standard format: 

```sh
twitch <product> <action>
```

The commands are described below, and any accompanying args/flags will be in the accompanying subsections.

## Commands

The CLI currently supports the following products: 

- [api](./docs/api.md)
- [configure](./docs/configure.md)
- [event](docs/event.md)
- [token](docs/token.md)
- [version](docs/version.md)

## Contributing

Check out [CONTRIBUTING.md](./CONTRIBUTING.md) for notes on making contributions.

## License 

This library is licensed under the Apache 2.0 License.

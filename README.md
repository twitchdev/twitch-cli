# Twitch CLI (Early Preview Version)

- [Twitch CLI (Early Preview Version)](#twitch-cli-early-preview-version)
  - [Download](#download)
  - [Usage](#usage)
  - [Commands](#commands)
  - [Contributing](#contributing)
  - [License](#license)

## Download

To download, go to the [Releases tab of GitHub](https://github.com/twitchdev/twitch-cli/releases). The examples in the documentation assume you have put this into your PATH and renamed to `twitch` (or symlinked as such).

**Note**: if using Mac OS, you may need to adjust the permissions of the file to allow for execution. This is only temporary while we work to get this into Homebrew in the future. 

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
# Twitch CLI 

- [Twitch CLI](#twitch-cli)
  - [Download](#download)
    - [Homebrew](#homebrew)
    - [Scoop](#scoop)
    - [WinGet](#winget)
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

### WinGet

Alternatively on Windows you can use [WinGet](https://learn.microsoft.com/en-us/windows/package-manager/winget/) for installing the CLI

To install via Winget, run:

```sh
winget install Twitch.TwitchCLI
```

### Manual Download

To download, go to the [Releases tab of GitHub](https://github.com/twitchdev/twitch-cli/releases). The examples in the documentation assume you have put this into your PATH and renamed to `twitch` (or symlinked as such).

**Note**: If using MacOS and downloading manually, you may need to adjust the permissions of the file to allow for execution.

To do so, please run: `chmod 755 <filename>` where the filename is the name of the downloaded binary. 

## Updating

To update the Twitch CLI, run the command relevant to your installation method.

**NOTE:** Once a day the program will make an HTTP call to Github to check if the application is of the latest version. For information on disabling this, see *Disabling release version checks and notices* below.

### Homebrew

To update using Homebrew, run:

```sh
brew upgrade twitchdev/twitch/twitch-cli
```

### Scoop

To update using Scoop, run:

```sh
scoop update twitch-cli
```

### WinGet

To update using WinGet, run:

```sh
winget update Twitch.TwitchCLI
```

### Manual Download

To download, go to the [Releases tab of GitHub](https://github.com/twitchdev/twitch-cli/releases). The examples in the documentation assume you have put this into your PATH and renamed to `twitch` (or symlinked as such).

**Note**: If using MacOS and downloading manually, you may need to adjust the permissions of the file to allow for execution.

To do so, please run: `chmod 755 <filename>` where the filename is the name of the downloaded binary.

## Disabling release version checks and notices

When the Twitch CLI exits successfully, the application will automatically check the Twitch CLI's Github releases at the following URL:

```
https://api.github.com/repos/twitchdev/twitch-cli/releases/latest
```

If the version of the Twitch CLI you are running is older than the latest released version, a notice will be printed to the console.

To prevent this from happening, make one of the following changes:

- Set the environment variable `CI` to `true`
- Set the environment variable `TWITCH_DISABLE_UPDATE_CHECKS` to `true`
- Add `DISABLE_UPDATE_CHECKS=true` to your **.twitch-cli.env** configuration file
- SET `LAST_UPDATE_CHECK` to `3000-01-01` in your **.twitch-cli.env** configuration file, which will prevent it from running until the year 3000

If you're running the Twitch CLI in a CI/CD environment, most environments will have already set the `CI` environment variable to `true`.

## Usage

The CLI largely follows a standard format: 

```sh
twitch <product> <action>
```

The commands are described below, and any accompanying args/flags will be in the accompanying subsections.

## Commands

The CLI currently supports the following products: 

- [api](./docs/api.md)
- [completion](./docs/completion.md)
- [configure](./docs/configure.md)
- [event](docs/event.md)
- [mock-api](docs/mock-api.md)
- [token](docs/token.md)
- [version](docs/version.md)

## Contributing

Check out [CONTRIBUTING.md](./CONTRIBUTING.md) for notes on making contributions.

## License 

This library is licensed under the Apache 2.0 License.

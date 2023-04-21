# Contributing

Thanks for helping make the Twitch CLI better! 
- [Contributing](#contributing)
  - [Design Principles](#design-principles)
  - [Report an Issue](#report-an-issue)
  - [Contributing Code with Pull Requests](#contributing-code-with-pull-requests)
    - [Requirements](#requirements)
    - [Profiling](#profiling)
  - [Code of Conduct](#code-of-conduct)
  - [Licensing](#licensing)

## Design Principles

Contributions to the Twitch CLI should align with the project’s design principles:

 * Maintain backwards compatibility whenever possible. This tool has an opportunity to be used in CI/CD pipelines and breaking changes should be avoided at all costs
 * Use only publicly documented endpoints. We will not accept PRs for functionality that leverages on undocumented endpoints
 * Limit dependencies where possible, so that they are easier to integrate and upgrade


Examples of contributions that should be addressed with high priority:

 * Security updates.
 * Performance improvements.
 * Supporting new versions of key dependencies such as Go or Cobra
 * Documentation

## Report an Issue

If you have run into a bug or want to discuss a new feature, please [file an issue](https://github.com/twitchdev/twitch-cli/issues).

## Contributing Code with Pull Requests

The Twitch CLI uses [Github pull requests](https://github.com/twitchdev/twitch-cli/pulls). Fork, hack away at your changes and submit. Most pull requests will go through a few iterations before they get merged. Different contributors will sometimes have different opinions, and often patches will need to be revised before they can get merged.

### Requirements

 *  The Twitch CLI officially supports Mac, Windows, and Linux Intel-based systems
 *  All commands and functionality should be documented appropriately
 *  All new functionality/features should have appropriate unit testing

To confirm it will build with these systems, feel free to run `make build_all`. 

The Twitch CLI strives to have a consistent set of documentation that matches the command structure and any new functionality must have accompanying documentation in the PR.

As noted in the [README](./README.md), all commands follow the following structure: `twitch <product> <action>`. Each product should live within it's own file in the `cmd` directory, with the applicable actions within it. The logic is then split into the `internal` directory. 

Some commands may not be part of a designated product (for example, the `token` and `version` commands) - if you are building functionality that is not tied to a Twitch product, please open the PR to discuss further. 

### Profiling

The Twitch CLI makes use of [pprof](https://github.com/google/pprof) for CPU profiling. This can be enabled on any system by setting the environment variable `TWITCH_CLI_ENABLE_CPU_PROFILER` to `true`. 
By default, the CPU profile will be written to your system as `cpu.prof` when the program exits. This filename can be modified with the environment variable `TWITCH_CLI_CPU_PROFILER_FILE`. 

## Code of Conduct

This project has adopted the [Amazon Open Source Code of Conduct](https://aws.github.io/code-of-conduct).
For more information see the [Code of Conduct FAQ](https://aws.github.io/code-of-conduct-faq) or contact
opensource-codeofconduct@amazon.com with any additional questions or comments.

## Licensing

See the [LICENSE](https://github.com/twitchdev/twitch-cli//blob/master/LICENSE) file for our project's licensing. We will ask you to confirm the licensing of your contribution.

We may ask you to sign a [Contributor License Agreement (CLA)](http://en.wikipedia.org/wiki/Contributor_License_Agreement) for larger changes.


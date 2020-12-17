# Configure

Allows the user to set basic information required for CLI usage via interactive prompt. At the moment, this is just the Client ID and Secret. After this is run, the ability to run both [`token`](token.md) and [`api`](api.md) commands are enabled, as they both require OAuth tokens which need both a Client ID and Secret.

This should be the first step if the intent is to use either of those functionalities.

**Args**

None.

**Flags**

None

**Examples**

```sh
twitch configure // configures the CLI tool
```
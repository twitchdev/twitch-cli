# Configure

Allows the user to set basic information required for CLI usage via interactive prompt. At the moment, this is just the Client ID and Secret. After this is run, the ability to run both [`token`](token.md) and [`api`](api.md) commands are enabled, as they both require OAuth tokens which need both a Client ID and Secret.

This should be the first step if the intent is to use either of those functionalities.

If you'd prefer to not use the interactive shell, you can pass the settings via the below flags.

**Args**

None.

**Flags**

| Flag              | Shorthand | Description                       | Example                    | Required? (Y/N) |
|-------------------|-----------|-----------------------------------|----------------------------|-----------------|
| `--client-id`     | `-i`      | Client ID to use for the CLI.     | `configure -i test_client` | N               |
| `--client-secret` | `-s`      | Client Secret to use for the CLI. | `configure -s test_secret` | N               |


**Examples**

```sh
twitch configure // configures the CLI tool
```
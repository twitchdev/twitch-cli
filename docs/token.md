# Token

Used to fetch an access token. This can fetch either an app access token or user token (with scopes) and stores for further use. It is important to note that this requires the running of `twitchdev configure` before it'll fetch tokens. If you don't run [`configure`](configure.md) prior to running `token`, you will be prompted to configure the application.

**IMPORTANT** 

To use user tokens, you will need to set up your Client ID provided during [`configure`](configure.md) with a redirect URI of `http://localhost:3000`. App access tokens will work without this step. You can configure that [on the Twitch Developer console](https://dev.twitch.tv/console).

**Args**

None.


**Flags**

| Flag           | Shorthand | Description                                                  | Example                          | Required? (Y/N) |
|----------------|-----------|--------------------------------------------------------------|----------------------------------|-----------------|
| `--user-token` | `-u`      | Whether to fetch a user token or not. Default is false.      | `token -u`                       | N               |
| `--scopes`     | `-s`      | The space separated scopes to use when getting a user token. | `-s "user:read:email user_read"` | N               |

**Examples**

```sh
twitch token -u -s "user:read:email" // gets a user token with the user:read:email scope
twitch token // fetches an app access token
```
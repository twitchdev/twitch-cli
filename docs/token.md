# Token

## Overview

The `token` product is used to fetch access tokens for an application

## Token Types

The Twitch CLI provides access to two types of tokens: App Access Tokens and User Access Tokens. 

**App Access Tokens** 

[App Access Tokens](https://dev.twitch.tv/docs/authentication/#app-access-tokens) can access APIs that don't require the user's permission. For example, getting a list of videos. 

**User Access Tokens** 

[User Access Tokens](https://dev.twitch.tv/docs/authentication/#user-access-tokens) provide access to APIs a user must grant permission to. For example, starting or stopping a poll. The specific APIs and functionality granted to a token is defined via [scopes](https://dev.twitch.tv/docs/authentication/scopes/).

## Client IDs and Client Secrets

Getting Access Tokens requires use of a Client ID and Client Secret which are associated with a Twitch Developer's Application. Creating an application is done by registration. Details on that process [are here](https://dev.twitch.tv/docs/authentication/register-app/). Client IDs are generated automatically when an application is registered. Client Secrets must be generated explicitly. This can be done by visiting [the Developer Applications Console](https://dev.twitch.tv/console/apps), choosing "Manage" for the app, then using the "New Secret" button at the bottom of the page. 

Adding the Client ID and Client Secret to the CLI tool is done with:

```
twitch configure
```

Running that starts prompts asking for the credentials. 

## Fetching App Access Tokens

App Access Tokens can be fetched once the Client ID and Client Secret have been entered. No other configuration is necessary. The command is:

```
twitch token
```

Running that returns a result with the token like:

```
2023/08/23 13:19:08 App Access Token: 01234abcdetc...
```


## Fetching User Access Tokens

Fetching User Access Tokens requires setting an _OAuth Redirect URL_. Those URLs are defined on the _Manage_ page for each app in the [Developer's Application Console](https://dev.twitch.tv/console/apps). The twitch CLI uses `http://localhost:3000`. Two important notes when adding that to the OAuth Redirect URLs section:

1. Do not add a `/` to the end of the URL (i.e. use `http://localhost:3000` and not `http://localhost:3000/`)
2. The OAuth Redirect URLs section allows up to ten URLs. The `http://localhost:3000` must be the first one

**The User Flag**

The `-u` flag is what sets the `token` product to fetch a User Access Token instead of an App Access Token. 


**Scopes**

User Access Tokens use scopes to determine which APIs and features they have access to. The requested scopes are defined via a space separated list following an `-s` flag with the `token` product. 

The full list of available scopes [here in the Twitch Documentation](https://dev.twitch.tv/docs/authentication/scopes/)

**Example**

A full example fetching a User Access Token with the ability to do shoutouts and set shield modes looks like this:

```
twitch token -u -s "moderator:manage:shoutouts moderator:manage:shield_mode"
```

Running that produce some initial output in the terminal and opens a browser to a Twitch authorization page. If you're not already signed in, you'll be asked to do so. When signed-in, the page displays the authorization request including the requested scopes. Clicking the "Authorize" button at the bottom redirects the browser back to the `http://localhost:3000` address where the `twitch` CLI picks it up and complete the process by parsing the data returned in the URL. 

The browser will display a message like:

```
Feel free to close this browser window.
```

The terminal outputs the tokens in a response like this:

```
Opening browser. Press Ctrl+C to cancel...
2023/08/23 13:50:00 Waiting for authorization response ...
2023/08/23 13:50:03 Closing local server ...
2023/08/23 13:50:10 User Access Token: 012345asdfetc...
Refresh Token: 012345asdfetc...
Expires At: 2023-08-23 22:06:47.036137 +0000 UTC
Scopes: [moderator:manage:shield_mode moderator:manage:shoutouts]
```

## Revoking Access Tokens

Access tokens can be revoked with:

```
twitch token -r 0123456789abcdefghijABCDEFGHIJ
```

## Alternate IP for User Token Webserver

If you'd like to bind the webserver used for user tokens (`-u` flag), you can override it with the `--ip` flag. For example:

```
twitch token -u --ip 127.0.0.1"
```

## Alternate Port

Port 3000 on localhost is used by default when fetching User Access Tokens. The `-p` flag can be used to change to another port if another service is already occupying that port. For example:

```
twitch token -u -p 3030 -s "moderator:manage:shoutouts moderator:manage:shield_mode"
```

NOTE: You must update the first entry in the _OAuth Redirect URLs_ section of your app's management page in the [Developer's Application Console](https://dev.twitch.tv/console/apps) to match the new port number. Make sure there is no `/` at the end of the URL (e.g. use `http://localhost:3030` and not `http://localhost:3030/`) and that the URL is the first entry in the list if there is more than one.


## Alternate Host

If you'd like to change the hostname for one reason or another (e.g. binding to a local domain), you can use the `--redirect-host` to change the domain. You should _not_ prefix it with `http` or `https`.

Example: 

```
twitch token -u --redirect-host contoso.com
```

NOTE: You must update the first entry in the _OAuth Redirect URLs_ section of your app's management page in the [Developer's Application Console](https://dev.twitch.tv/console/apps) to match the new port number. Make sure there is no `/` at the end of the URL (e.g. use `http://localhost:3030` and not `http://localhost:3030/`) and that the URL is the first entry in the list if there is more than one.


## Errors

This error occurs when there's a problem with the OAuth Redirect URLs. Check in the app's management page in the [Developer's Application Console](https://dev.twitch.tv/console/apps) to ensure the first entry is set to `http://localhost:3000`. Specifically, verify that your using `http` and not `https` and that the URL does not end with a `/`. (If you've changed ports with the `-p` flag, ensure those numbers match as well)

```
Error! redirect_mismatch
Error Details: Parameter redirect_uri does not match registered URI
```

## Command Line Notes

**Args**

None.


**Flags**

| Flag              | Shorthand | Description                                                                                                    | Example                                      | Required? (Y/N) |
|-------------------|-----------|----------------------------------------------------------------------------------------------------------------|----------------------------------------------|-----------------|
| `--user-token`    | `-u`      | Whether to fetch a user token or not. Default is false.                                                        | `token -u`                                   | N               |
| `--scopes`        | `-s`      | The space separated scopes to use when getting a user token.                                                   | `-s "user:read:email user_read"`             | N               |
| `--revoke`        | `-r`      | Instead of generating a new token, revoke the one passed to this parameter.                                    | `-r 0123456789abcdefghijABCDEFGHIJ`          | N               |
| `--ip`            |           | Manually set the port to be used for the User Token web server. The default binds to all interfaces. (0.0.0.0) | `--ip 127.0.0.1`                             | N               |
| `--port`          | `-p`      | Override/manually set the port for token actions. (The default is 3000)                                        | `-p 3030`                                    | N               |
| `--client-id`     |           | Override/manually set client ID for token actions. By default client ID from CLI config will be used.          | `--client-id uo6dggojyb8d6soh92zknwmi5ej1q2` | N               |
| `--redirect-host` |           | Override/manually set the redirect host token actions. The default is `localhost`                              | `--redirect-host contoso.com`                | N               |

## Notes

- If you've already authorized the app, the webpage will redirect back immediately without requiring any interaction
- You'll be asked to fill in the Client ID and Client Secret if you run the `token` product without having already set them

# api

- [api](#api)
  - [Arguments](#arguments)
  - [get](#get)
  - [post](#post)
  - [put](#put)
  - [patch](#patch)
  - [delete](#delete)


The `api` product enables users to interact with the [Twitch API](https://dev.twitch.tv/docs/api) via CLI. It supports both query parameters and bodies for applicable endpoints, and all standard HTTP methods. 

The format is `api <method> <url> <flags>`.

## Arguments

All API commands accept one of two formats: 
1. The endpoint with a leading slash, for example: `twitch api get /users/follows`
2. The endpoint without slashes, such as `twitch api patch channels`

## get

Allows the user to make GET calls to endpoints on Helix. Requires a logged in token from the [`token`](token.md) command.

**Args**

[See Arguments, above.](#arguments)

**Flags**

| Flag             | Shorthand | Description                                                                                                                                                         | Example              | Required? (Y/N) |
|------------------|-----------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------|----------------------|-----------------|
| `--query-param`  | `-q`      | Query parameters for the endpoint in `key=value` format. Multiple can be entered to give multiple parameters.                                                       | `get -q login=ninja` | N               |
| `--unformatted`  | `-u`      | Whether to return unformatted responses. Default is `false`.                                                                                                        | `get -u`             | N               |
| `--autopaginate` | `-P`      | Whether to autopaginate the response from Twitch **WARNING** This flag can cause extremely large payloads and cause issues with some terminals. Default is `false`. | `get -P`             | N               |

**Examples**

```sh
twitch api get users follows -q from_id=44635596 // gets user follows from user ID 44635596
twitch api get /subscriptions -q broadcaster_id=44635596 // gets subscriptions to broadcaster 44635596
```

## post

Allows the user to make POST calls to endpoints on Helix. Requires a logged in token from the [`token`](token.md) command.

**Args**

[See Arguments, above.](#arguments)

**Flags**

| Flag             | Shorthand | Description                                                                                                   | Example               | Required? (Y/N) |
|------------------|-----------|---------------------------------------------------------------------------------------------------------------|-----------------------|-----------------|
| `--query-param`  | `-q`      | Query parameters for the endpoint in `key=value` format. Multiple can be entered to give multiple parameters. | `post -q login=ninja` | N               |
| `--body`         | `-b`      | Body for the request. Supports CURL-like references to files using the format of `@data.json `.               | `post -b @data.json`  | N               |
| `--pretty-print` | `-p`      | Whether to pretty-print API requests. Default is `true`.                                                      | `post -p`             | N               |

**Examples**

```sh
twitch api post users follows -q from_id=44635596 -q to_id=135093069
```

## put

Allows the user to make PUT calls to endpoints on Helix. Requires a logged in token from the [`token`](token.md) command.

**Args**

[See Arguments, above.](#arguments)

**Flags**

| Flag             | Shorthand | Description                                                                                                   | Example              | Required? (Y/N) |
|------------------|-----------|---------------------------------------------------------------------------------------------------------------|----------------------|-----------------|
| `--query-param`  | `-q`      | Query parameters for the endpoint in `key=value` format. Multiple can be entered to give multiple parameters. | `put -q login=ninja` | N               |
| `--body`         | `-b`      | Body for the request. Supports CURL-like references to files using the format of `@data.json `.               | `put -b @data.json`  | N               |
| `--pretty-print` | `-p`      | Whether to pretty-print API requests. Default is `true`.                                                      | `put -p`             | N               |

**Examples**

```sh
twitch api put users -q "description=hi mom" 
```

## patch

Allows the user to make PATCH calls to endpoints on Helix. Requires a logged in token from the [`token`](token.md) command.

**Args**

[See Arguments, above.](#arguments)

**Flags**

| Flag             | Shorthand | Description                                                                                                   | Example                | Required? (Y/N) |
|------------------|-----------|---------------------------------------------------------------------------------------------------------------|------------------------|-----------------|
| `--query-param`  | `-q`      | Query parameters for the endpoint in `key=value` format. Multiple can be entered to give multiple parameters. | `patch -q login=ninja` | N               |
| `--body`         | `-b`      | Body for the request. Supports CURL-like references to files using the format of `@data.json `.               | `patch -b @data.json`  | N               |
| `--pretty-print` | `-p`      | Whether to pretty-print API requests. Default is `true`.                                                      | `patch -p`             | N               |


**Examples**

```sh
twitch api patch channels -q broadcaster_id=44635596 -b '{"game_id":"394568"}' 
```

## delete

Allows the user to make DELETE calls to endpoints on Helix. Requires a logged in token from the [`token`](token.md) command.

**Args**

[See Arguments, above.](#arguments)

**Flags**

| Flag            | Shorthand | Description                                                                                                   | Example                | Required? (Y/N) |
|-----------------|-----------|---------------------------------------------------------------------------------------------------------------|------------------------|-----------------|
| `--query-param` | `-q`      | Query parameters for the endpoint in `key=value` format. Multiple can be entered to give multiple parameters. | `patch -q login=ninja` | N               |

**Examples**

```sh
twitch api delete users follows -q from_id=44635596 -q to_id=135093069
```
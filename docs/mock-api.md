# mock-api

- [mock-api](#mock-api)
  - [Description](#description)
  - [generate](#generate)
  - [start](#start)
    - [mock namespace](#mock-namespace)
    - [units namespace](#units-namespace)
    - [auth namespace](#auth-namespace)

## Description

The `mock-api` product has two primary functions. The first is to generate so-called `units`- those are core building blocks on Twitch. These are, for now:

* Application Clients
* Categories
* Streams
* Subscriptions
* Tags
* Teams
* Users

The second is the actual server used to mock the endpoints. In the next iteration, you will be able to edit these and add further ones manually (for example, making a user with specific attributes), but for the beta we won't be providing this functionality as the current `generate` feature will make all of these (and more).

As of the 1.1 release, this product is in an **open beta** and any bugs should be filed via GitHub Issues. Given the breadth of the tool, it is likely you may run across oddities; please fill out an issue if that is the case. 

All commands exit the program with a non-zero exit code when the command fails, including when the mock API server fails to start.

## generate

This command will generate a specified number of users with associated relationships (e.g. subscriptions/mods/blocks).

**Args**

None.

**Flags**

| Flag      | Shorthand | Description                                                                 | Example | Required? (Y/N) |
|-----------|-----------|-----------------------------------------------------------------------------|---------|-----------------|
| `--count` | `-c`      | Number of users to generate (and associated relationships). Defaults to 10. | `-c 25` | N               |


## start

The `start` function starts a new mock server for use with testing functionality. Currently, this replicates a large majority of the current API endpoints on the new API, but are omitting: 

* Extensions endpoints
* Code entitlement endpoints
* Websub endpoints
* EventSub endpoints

For many of these, we are exploring how to better integrate this with existing features (for example, allowing events to be triggered on unit creation or otherwise), and for others, the value is minimal compared to the docs. All other endpoints should be currently supported, however it is possible to be out of date- if so, [please raise an issue](https://github.com/twitchdev/twitch-cli/issues). 

To access these endpoints, you will point your code to `http://localhost:<port>/mock`, where port is either the default port (8080) or one you specified using the flag. For example, to access the users endpoint:

```sh
curl -i -H "Accept: application/json" http://localhost:8080/mock/users
```

For information on accessing those endpoints, please see [the documentation on the Developer site](https://dev.twitch.tv/docs/api/reference).

In total, there are three namespaces (top-level folder) that are used:

### mock namespace

Example URL: `http://localhost:8080/mock/users`

This namespace houses all mock endpoints. For information on accessing those endpoints, please see [the documentation on the Developer site](https://dev.twitch.tv/docs/api/reference).

### units namespace

Example URL: `http://localhost:8080/units/users`

This endpoint gives an unauthenticated peek into the list of units in the database- used for debugging or finding units for testing. Endpoints include:

* GET /categories
* GET /clients
* GET /streams
* GET /subscriptions
* GET /tags
* GET /teams
* GET /users
* GET /videos

More will be added in the future. 

### auth namespace

This endpoint is a light implementation of OAuth, without support for OIDC. These endpoints are used to generate either an app access token or user token. The two endpoints are below, with documentation and examples using cURL. All tokens expire after 24 hours. 

**POST /authorize**

This endpoint generates a user token, similar to OAuth authorization code. 

| Query Parameter | Description                                                          | Example                  | Required? (Y/N) | 
|-----------------|----------------------------------------------------------------------|--------------------------|-----------------|
| `client_id`     | Application client ID, which is output by the `generate` command.    | `?client_id=1234`        | Y               |  
| `client_secret` | Application client secret, which is output by the `generate` command | `?client_secret=1234`    | Y               |   
| `grant_type`    | Must be `user_token`                                                 | `?grant_type=user_token` | Y               |   
| `user_id`       | User to get the token for.                                           | `?user_id=1234`          | Y               |   
| `scope`         | Space separated list of scopes to request for the given user.        | `?scope=bits:read`       | N               |   

The response is identical to the OAuth `authorization_code` with the omission of a refresh token. 

Example request for user 78910 with no scopes:

```sh
curl -X POST http://localhost:8080/auth/authorize?client_id=123&client_secret=456&grant_type=user_token&user_id=78910
```

Example response:

```json
{
    "access_token": "ff4231a5befca12",
    "refresh_token": "",
    "expires_in": 86399,
    "scope": [],
    "token_type": "bearer"
}
```

Docs: https://dev.twitch.tv/docs/authentication/getting-tokens-oauth#oauth-authorization-code-flow

**POST /token**

This endpoint generates an app access token using the `client_credentials` flow as documented. 


| Query Parameter | Description                                                          | Example                          | Required? (Y/N) |   
|-----------------|----------------------------------------------------------------------|----------------------------------|-----------------|
| `client_id`     | Application client ID, which is output by the `generate` command.    | `?client_id=1234`                | Y               |  
| `client_secret` | Application client secret, which is output by the `generate` command | `?client_secret=1234`            | Y               |   
| `grant_type`    | Must be `client_credentials`                                         | `?grant_type=client_credentials` | Y               |   
| `scope`         | Space separated list of scopes to request for the given user.        | `?scope=bits:read`               | N               |   


The response is identical to the OAuth `client_credentials` flow with the omission of a refresh token. 

Example request with no scopes:

```sh
curl -X POST http://localhost:8080/auth/token?client_id=123&client_secret=456&grant_type=client_credentials
```

Example response:

```json
{
    "access_token": "4f5dce6cea626cb",
    "refresh_token": "",
    "expires_in": 86399,
    "scope": [],
    "token_type": "bearer"
}
```

Docs: https://dev.twitch.tv/docs/authentication/getting-tokens-oauth#oauth-client-credentials-flow

**Args**

None.

**Flags**

| Flag     | Shorthand | Description                              | Example   | Required? (Y/N) |
|----------|-----------|------------------------------------------|-----------|-----------------|
| `--port` | `-p`      | Port number to use with the mock server. | `-p 8000` | N               |



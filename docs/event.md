# Events

- [Events](#events)
  - [Description](#description)
  - [Trigger](#trigger)
  - [Retrigger](#retrigger)
  - [Verify-Subscription](#verify-subscription)
  - [Websocket](#websocket)

## Description

The `event` product contains commands to trigger mock events for local webhook testing or migration.

## Trigger

Used to either create or send mock events for use with local webhooks testing.

**Args**

This command can take either the Event or Alias listed as an argument. It is preferred that you work with the Event, but for backwards compatibility Aliases still work.

| Event                                                    | Alias                 | Description |
|----------------------------------------------------------|-----------------------|-------------|
| `channel.ban`                                            | `ban`                 | Channel ban event. |
| `channel.channel_points_custom_reward.add`               | `add-reward`          | Channel Points event for a Custom Reward being added. |
| `channel.channel_points_custom_reward.remove`            | `remove-reward`       | Channel Points event for a Custom Reward being removed. |
| `channel.channel_points_custom_reward.update`            | `update-reward`       | Channel Points event for a Custom Reward being updated. |
| `channel.channel_points_custom_reward_redemption.add`    | `add-redemption`      | Channel Points EventSub event for a redemption being performed. |
| `channel.channel_points_custom_reward_redemption.update` | `add-update`          | Channel Points EventSub event for a redemption being performed. |
| `channel.charity_campaign.donate`                        | `charity-donate`      | Charity campaign donation occurance event. |
| `channel.charity_campaign.progress`                      | `charity-progress`    | Charity campaign progress event. |
| `channel.charity_campaign.start`                         | `charity-start`       | Charity campaign start event. |
| `channel.charity_campaign.stop`                          | `charity-stop`        | Charity campaign stop event. |
| `channel.cheer`                                          | `cheer`               | Channel event for receiving cheers. |
| `channel.follow`                                         | `follow`              | Channel event for receiving a follow. |
| `channel.goal.begin`                                     | `goal-begin`          | Channel creator goal start event. |
| `channel.goal.end`                                       | `goal-end`            | Channel creator goal end event. |
| `channel.goal.progress`                                  | `goal-progress`       | Channel creator goal progress event. |
| `channel.hype_train.begin`                               | `hype-train-begin`    | Channel hype train start event. |
| `channel.hype_train.end`                                 | `hype-train-end`      | Channel hype train start event. |
| `channel.hype_train.progress`                            | `hype-train-progress` | Channel hype train start event. |
| `channel.moderator.add`                                  | `add-moderator`       | Channel moderator add event. |
| `channel.moderator.remove`                               | `remove-moderator`    | Channel moderator removal event. |
| `channel.poll.begin`                                     | `poll-begin`          | Channel poll begin event. |
| `channel.poll.end`                                       | `poll-end`            | Channel poll end event. |
| `channel.poll.progress`                                  | `poll-progress`       | Channel poll progress event. |
| `channel.prediction.begin`                               | `prediction-begin`    | Channel prediction begin event. |
| `channel.prediction.end`                                 | `prediction-end`      | Channel prediction end event. |
| `channel.prediction.lock`                                | `prediction-lock`     | Channel prediction lock event. |
| `channel.prediction.progress`                            | `prediction-progress` | Channel prediction progress event. |
| `channel.raid`                                           | `raid`                | Channel raid event with a random viewer count. |
| `channel.shield_mode.begin`                              | `shield-mode-begin`   | Channel Shield Mode activate event. |
| `channel.shield_mode.end`                                | `shield-mode-end`     | Channel Shield Mode deactivate event. |
| `channel.shoutout.create`                                | `shoutout-create`     | Channel shoutout created event. This is for outgoing shoutouts, from your channel to another. |
| `channel.shoutout.receive`                               | `shoutout-received`   | Channel shoutout created event. This is for incoming shoutouts, to your channel from anothers. |
| `channel.subscribe`                                      | `subscribe`           | A standard subscription event. Triggers a basic tier 1 sub, but can be flexible with --tier |
| `channel.subscribe`                                      | `gift`                | A gifted subscription event. Triggers a basic tier 1 sub, but can be flexible with --tier |
| `channel.subscription.end`                               | `unsubscribe`         | A standard subscription end event. Triggers a basic tier 1 sub, but can be flexible with --tier |
| `channel.subscription.gift`                              | `channel-gift`        | Channel gifting event; not to be confused with the `gift` event. This event is a description of the number of gifts given by a user. |
| `channel.subscription.message`                           | `subscribe-message`   | Subscription Message event. |
| `channel.unban`                                          | `unban`               | Channel unban event. |
| `channel.update`                                         | `stream-change`       | Channel update event. When a broadcaster updates channel properties. |
| `drop.entitlement.grant`                                 | `drop`                | Drop Entitlement event. |
| `extension.bits_transaction.create`                      | `transaction`         | Bits in Extensions transactions events. |
| `stream.offline`                                         | `streamdown`          | Stream offline event. |
| `stream.online`                                          | `streamup`            | Stream online event. |
| `user.authorization.grant`                               | `grant`               | Authorization grant event. |
| `user.authorization.revoke`                              | `revoke`              | User authorization revoke event. Uses local Client as set in `twitch configure` or generates one randomly. |




**Flags**

| Flag                      | Shorthand | Description                                                                                                                     | Example                                      | Required? (Y/N) |
|---------------------------|-----------|---------------------------------------------------------------------------------------------------------------------------------|----------------------------------------------|-----------------|
| `--anonymous`             | `-a`      | Denotes if the event is anonymous. Only applies to Gift and Sub events.                                                         | `-a`                                         | N               |
| `--charity-current-value` |           | For charity events, manually set the charity dollar value.                                                                      | `--charity-current-value 11000`              | N               |
| `--charity-target-value`  |           | For charity events, manually set the charity dollar value.                                                                      | `--charity-current-value 23400`              | N               |
| `--client-id`             |           | Manually set the Client ID used for revoke, grant, and bits transactions.                                                       | `--client-id 4ofh8m0706jqpholgk00u3xvb4spct` | N               |
| `--cost`                  | `-C`      | Amount of subscriptions, bits, or channel points redeemed/used in the event.                                                    | `-C 250`                                     | N               |
| `--count`                 | `-c`      | Count of events to fire. This can be used to simulate an influx of events.                                                      | `-c 100`                                     | N               |
| `--description`           | `-d`      | Title the stream should be updated/started with.                                                                                | `-d Awesome new title!`                      | N               |
| `--event-status`          | `-S`      | Status of the Event object (.event.status in JSON); Currently applies to channel points redemptions.                            | `-S fulfilled`                               | N               |
| `--forward-address`       | `-F`      | Web server address for where to send mock events.                                                                               | `-F https://localhost:8080`                  | N               |
| `--from-user`             | `-f`      | Denotes the sender's TUID of the event, for example the user that follows another user or the subscriber to a broadcaster.      | `-f 44635596`                                | N               |
| `--game-id`               | `-G`      | Game ID for Drop or other relevant events.                                                                                      | `-G 1234`                                    | N               |
| `--gift-user`             | `-g`      | Used only for subcription-based events, denotes the gifting user ID.                                                            | `-g 44635596`                                | N               |
| `--item-id`               | `-i`      | Manually set the ID of the event payload item (for example the reward ID in redemption events or game in stream events).        | `-i 032e4a6c-4aef-11eb-a9f5-1f703d1f0b92`    | N               |
| `--item-name`             | `-n`      | Manually set the name of the event payload item (for example the reward ID in redemption events or game name in stream events). | `-n "Science & Technology"`                  | N               |
| `--secret`                | `-s`      | Webhook secret. If defined, signs all forwarded events with the SHA256 HMAC and must be 10-100 characters in length.            | `-s testsecret`                              | N               |
| `--session`               |           | WebSocket session to target. Only used when forwarding to WebSocket servers with --transport=websocket                          | `--session e411cc1e_a2613d4e`                | N               |
| `--subscription-id`       | `-u`      | Manually set the subscription/event ID of the event itself.                                                                     | `-u 5d3aed06-d019-11ed-afa1-0242ac120002`    | N               |
| `--subscription-status`   | `-r`      | Status of the Subscription object (.subscription.status in JSON). Defaults to "enabled"                                         | `-r revoked`                                 | N               |
| `--tier`                  |           | Tier of the subscription.                                                                                                       | `--tier 3000`                                | N               |
| `--timestamp`             |           | Sets the timestamp to be used in payloads and headers. Must be in RFC3339Nano format.                                           | `--timestamp 2017-04-13T14:34:23`            | N               |
| `--to-user`               | `-t`      | Denotes the receiver's TUID of the event, usually the broadcaster.                                                              | `-t 44635596`                                | N               |
| `--transport`             | `-T`      | The method used to send events. Can either be `webhook` or `websocket`. Default is `webhook`.                                   | `-T webhook`                                 | N               |


```sh
twitch event trigger subscribe -F https://localhost:8080/ // triggers a randomly generated subscribe event and forwards to the localhost:8080 server
twitch event trigger cheer -f 1234 -t 4567 // generates JSON for a cheer event from user 1234 to user 4567
```

## Retrigger

Allows previous events to be refired based on the event ID. The ID is noted within the event itself, such as in the "subscription" payload of standard webhooks.

For example, for:

```json
{
  "subscription": {
    "id": "713f3254-0178-9757-7439-d779400c0999",
    "type": "channels.cheer",
    ...
  }
}
```

The resulting ID would be `713f3254-0178-9757-7439-d779400c0999`.

**Args**
None

**Flags**

| Flag                | Shorthand | Description                                                                                                                                                   | Example                     | Required? (Y/N) |
|---------------------|-----------|---------------------------------------------------------------------------------------------------------------------------------------------------------------|-----------------------------|-----------------|
| `--forward-address` | `-F`      | Web server address for where to send mock events.                                                                                                             | `-F https://localhost:8080` | N               |
| `--id`              | `-i`      | The ID of the event to refire.                                                                                                                                | `-i <id>`                   | Y               |
| `--secret`          | `-s`      | Webhook secret. If defined, signs all forwarded events with the SHA256 HMAC and must be 10-100 characters in length.                                          | `-s testsecret`             | N               |

**Examples**

```sh
twitch event retrigger -i "713f3254-0178-9757-7439-d779400c0999" -F https://localhost:8080/ // triggers the previous cheer event to localhost:8080
```

## Verify-Subscription

Allows you to test if your webserver responds to subscription requests properly.

**Args**

This command takes the same arguments as [Trigger](#trigger).

**Flags**

| Flag                | Shorthand | Description                                                                                                          | Example                     | Required? (Y/N) |
|---------------------|-----------|----------------------------------------------------------------------------------------------------------------------|-----------------------------|-----------------|
| `--forward-address` | `-F`      | Web server address for where to send mock subscription.                                                              | `-F https://localhost:8080` | Y               |
| `--secret`          | `-s`      | Webhook secret. If defined, signs all forwarded events with the SHA256 HMAC and must be 10-100 characters in length. | `-s testsecret`             | N               |
| `--transport`       | `-T`      | The method used to send events. Default is `eventsub`.                                                               | `-T eventsub`               | N               |

**Examples**

```sh
twitch event verify-subscription cheer -F https://localhost:8080/ // triggers a fake "cheer" EventSub subscription and validates if localhost responds properly
```

## WebSocket

Provides access to a mock EventSub WebSocket server. More information can be found on [Twitch Developers documentation](https://dev.twitch.tv/docs/cli/websocket-event-command/).

**Args**

| Arg          | Description |
|--------------|-------------|
| start-server | Attempts to start the websocket sever. Default port is 8080. |
| reconnect    | Server command. Starts reconnect testing on the active WebSocket server. See documentation for more info. |
| close        | Server command. Closes a specific client connection with the provided WebSocket close code. |
| subscription | Server command. Modifies an existing subscription on the WebSocket server. |

**Flags used with start-server**
| Flag                     | Shorthand | Description                                                                          | Example       |
|--------------------------|-----------|--------------------------------------------------------------------------------------|---------------|
| `--port`                 | `-p`      | Use to specify the port number to use in the localhost address. The default is 8080. | `--port=8080` |
| `--require-subscription` | `-S`      | 	Prevents the server from allowing subscriptions to be forwarded unless they have a subscription created. Also enables 10 second subscription requirement when a client connects. | `-S` |


**Flags used with all other sub-commands**
| Flag             | Shorthand | Description                                                                                                                  | Example |
|------------------|-----------|------------------------------------------------------------------------------------------------------------------------------|---------|
| `--session`      | `-s`      | Targets a specific client by the session_id given during its Welcome message.                                                | `twitch event websocket close --session=e411cc1e_a2613d4e` |
| `--reason`       |           | Specifies the Close message code you wish to close a client’s connection with. Only used with "twitch websocket close"       | `twitch event websocket close --reason=4006` |
| `--status`       |           | Specifies the Status code you wish to override an existing subscription’s status to. Only used with "twitch websocket close" | `twitch event websocket subscription --status=user_removed` |
| `--subscription` |           | Specifies the subscription ID you wish to target. Only used with “twitch websocket subscription”.	                          | `twitch event websocket subscription --subscription=48d3-b9a-f84c` |

**Examples**

```sh
twitch event websocket start-server
twitch event websocket reconnect
twitch event websocket close --session=e411cc1e_a2613d4e --reason=4006
twitch event websocket subscription --status=user_removed --subscription=82a855-fae8-93bff0
```
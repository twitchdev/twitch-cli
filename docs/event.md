# Events

- [Events](#events)
  - [Description](#description)
  - [Trigger](#trigger)
  - [Retrigger](#retrigger)
  - [Verify-Subscription](#verify-subscription)

## Description

The `event` product contains commands to trigger mock events for local webhook testing or migration.

## Trigger

Used to either create or send mock events for use with local webhooks testing.

**Args**

| Argument              | Description                                                                                                |
|-----------------------|------------------------------------------------------------------------------------------------------------|
| `subscribe`           | A standard subscription event. Triggers a basic tier 1 sub.                                                |
| `unsubscribe`         | A standard unsubscribe event. Triggers a basic tier 1 sub.                                                 |
| `gift`                | A gifted subscription event. Triggers a basic tier 1 sub.                                                  |
| `cheer`               | Only usable with the `eventsub` transport, shows Cheers from chat.                                         |
| `transaction`         | Bits in Extensions transactions events.                                                                    |
| `add-reward`          | Channel Points EventSub event for a Custom Reward being added.                                             |
| `update-reward`       | Channel Points EventSub event for a Custom Reward being updated.                                           |
| `remove-reward`       | Channel Points EventSub event for a Custom Reward being removed.                                           |
| `add-redemption`      | Channel Points EventSub event for a redemption being performed.                                            |
| `update-redemption`   | Channel Points EventSub event for a redemption being updated.                                              |
| `raid`                | Channel Raid event with a random viewer count.                                                             |
| `revoke`              | User authorization revoke event. Uses local Client as set in `twitch configure` or generates one randomly. |
| `stream-change`       | Stream Changed event.                                                                                      |
| `streamup`            | Stream online event.                                                                                       |
| `streamdown`          | Sstream offline event.                                                                                     |
| `add-moderator`       | Channel moderator add event.                                                                               |
| `remove-moderator`    | Channel moderator removal event.                                                                           |
| `ban`                 | Channel ban event.                                                                                         |
| `unban`               | Channel unban event.                                                                                       |
| `hype-train-begin`    | Channel hype train start event.                                                                            |
| `hype-train-progress` | Channel hype train progress event.                                                                         |
| `hype-train-end`      | Channel hype train end event.                                                                              |



**Flags**

| Flag                | Shorthand | Description                                                                                                                     | Example                                   | Required? (Y/N) |
|---------------------|-----------|---------------------------------------------------------------------------------------------------------------------------------|-------------------------------------------|-----------------|
| `--forward-address` | `-F`      | Web server address for where to send mock events.                                                                               | `-F https://localhost:8080`               | N               |
| `--transport`       | `-T`      | The method used to send events. Default is eventsub, but can send using websub.                                                 | `-T websub`                               | N               |
| `--to-user`         | `-t`      | Denotes the receiver's TUID of the event, usually the broadcaster.                                                              | `-t 44635596`                             | N               |
| `--from-user`       | `-f`      | Denotes the sender's TUID of the event, for example the user that follows another user or the subscriber to a broadcaster.      | `-f 44635596`                             | N               |
| `--gift-user`       | `-g`      | Used only for subcription-based events, denotes the gifting user ID                                                             | `-g 44635596`                             | N               |
| `--secret`          | `-s`      | Webhook secret. If defined, signs all forwarded events with the SHA256 HMAC.                                                    | `-s testsecret`                           | N               |
| `--count`           | `-c`      | Count of events to fire. This can be used to simulate an influx of subscriptions.                                               | `-c 100`                                  | N               |
| `--anonymous`       | `-a`      | If the event is anonymous. Only applies to `gift` and `cheer` events.                                                           | `-a`                                      | N               |
| `--status`          | `-S`      | Status of the event object, currently applies to channel points redemptions.                                                    | `-S fulfilled`                            | N               |
| `--item-id`         | `-i`      | Manually set the ID of the event payload item (for example the reward ID in redemption events or game in stream events).        | `-i 032e4a6c-4aef-11eb-a9f5-1f703d1f0b92` | N               |
| `--item-name`       | `-n`      | Manually set the name of the event payload item (for example the reward ID in redemption events or game name in stream events). | `-n "Science & Technology"`               | N               |
| `--cost`            | `-C`      | Amount of bits or channel points redeemed/used in the event.                                                                    | `-C 250`                                  | N               |
| `--description`     | `-d`      | Title the stream should be updated/started with.                                                                                | `-d Awesome new title!`                   | N               |


**Examples**

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

| Flag                | Shorthand | Description                                                                  | Example                     | Required? (Y/N) |
| ------------------- | --------- | ---------------------------------------------------------------------------- | --------------------------- | --------------- |
| `--forward-address` | `-F`      | Web server address for where to send mock events.                            | `-F https://localhost:8080` | N               |
| `--id`              | `-i`      | The ID of the event to refire.                                               | `-i <id>`                   | Y               |
| `--secret`          | `-s`      | Webhook secret. If defined, signs all forwarded events with the SHA256 HMAC. | `-s testsecret`             | N               |

**Examples**

```sh
twitch event retrigger -i "713f3254-0178-9757-7439-d779400c0999" -F https://localhost:8080/ // triggers the previous cheer event to localhost:8080
```

## Verify-Subscription

Allows you to test if your webserver responds to subscription requests properly.

**Args**

| Argument              | Description                                                                                                |
|-----------------------|------------------------------------------------------------------------------------------------------------|
| `subscribe`           | A standard subscription event. Triggers a basic tier 1 sub.                                                |
| `unsubscribe`         | A standard unsubscribe event. Triggers a basic tier 1 sub.                                                 |
| `gift`                | A gifted subscription event. Triggers a basic tier 1 sub.                                                  |
| `cheer`               | Only usable with the `eventsub` transport, shows Cheers from chat.                                         |
| `transaction`         | Bits in Extensions transactions events.                                                                    |
| `add-reward`          | Channel Points EventSub event for a Custom Reward being added.                                             |
| `update-reward`       | Channel Points EventSub event for a Custom Reward being updated.                                           |
| `remove-reward`       | Channel Points EventSub event for a Custom Reward being removed.                                           |
| `add-redemption`      | Channel Points EventSub event for a redemption being performed.                                            |
| `update-redemption`   | Channel Points EventSub event for a redemption being updated.                                              |
| `raid`                | Channel Raid event with a random viewer count.                                                             |
| `revoke`              | User authorization revoke event. Uses local Client as set in `twitch configure` or generates one randomly. |
| `stream-change`       | Stream Changed event.                                                                                      |
| `streamup`            | Stream online event.                                                                                       |
| `streamdown`          | Sstream offline event.                                                                                     |
| `add-moderator`       | Channel moderator add event.                                                                               |
| `remove-moderator`    | Channel moderator removal event.                                                                           |
| `ban`                 | Channel ban event.                                                                                         |
| `unban`               | Channel unban event.                                                                                       |
| `hype-train-begin`    | Channel hype train start event.                                                                            |
| `hype-train-progress` | Channel hype train progress event.                                                                         |
| `hype-train-end`      | Channel hype train end event.                                                                              |

**Flags**

| Flag                | Shorthand | Description                                                                     | Example                     | Required? (Y/N) |
| ------------------- | --------- | ------------------------------------------------------------------------------- | --------------------------- | --------------- |
| `--forward-address` | `-F`      | Web server address for where to send mock subscription.                         | `-F https://localhost:8080` | Y               |
| `--secret`          | `-s`      | Webhook secret. If defined, signs all forwarded events with the SHA256 HMAC.    | `-s testsecret`             | N               |
| `--transport`       | `-T`      | The method used to send events. Default is eventsub, but can send using websub. | `-T websub`                 | N               |

**Examples**

```sh
twitch event verify-subscription cheer -F https://localhost:8080/ // triggers a fake "cheer" EventSub subscription and validates if localhost responds properly
```

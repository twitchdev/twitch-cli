# Drops
- [Drops](#drops)
  - [Description](#description)
  - [Export](#export)

## Description 

The `drops` command contains sub-commands to interact with the Drops product. 

## Export

Used to export Drops entitlements into a provided CSV filename. 

**Args**

None.


**Flags**

| Flag         | Shorthand | Description                                                                                           | Example                     | Required? (Y/N) |
|--------------|-----------|-------------------------------------------------------------------------------------------------------|-----------------------------|-----------------|
| `--filename` | `-f`      | Name of the CSV file to be generated.                                                                 | `-f drops_entitlements.csv` | Y               |
| `--game-id`  | `-T`      | ID of the game to be filtered for. If unsure, use [api get games](./api.md) described in the example. | `-g websub`                 | N               |
| `--user-id`  | `-t`      | Denotes the user's TUID of the entitlement to filter for.                                             | `-u 44635596`               | N               |

**Examples**

```sh
twitch api get games -q "name=Fortnite" 
twitch drops export -f "fortnite_drops_entitlements.csv" -g 33214
```

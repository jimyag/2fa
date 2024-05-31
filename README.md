
# 2fa

2fa is a two-factor authentication command line tool.

## Install

```bash
go get github.com/jimyag/2fa@latest
```

or download the binary from <https://github.com/jimyag/2fa/releases>

## Usage

### Add

Add a 2fa key

`2fa add <name> <key>`

```bash
2fa add jimyag 4LDRN6EUDSF3RNV7
```

### Get

Get a 2fa key and copy to clipboard

`2fa get <name> [--copy/-c]`

```bash
2fa get jimyag
2fa get jimyag --copy
```

### List

List all 2fa keys and display in a table

`2fa list`

```bash
2fa list
+--------+--------+------------+-----------+
| NAME   | TOTP   | LIFETIME/S | NEXT TOTP |
+--------+--------+------------+-----------+
| jimyag | 056907 |      2     |   134552  |
+--------+--------+------------+-----------+

```

### Delete

Delete a 2fa key

`2fa del <name>`

```bash
2fa del jimyag
```

# Link Keeper Bot

This telegram bot stores links that you send and provides some operations with them

## Flags
- token - token for telegram bot, emitted by [BotFather](https://t.me/BotFather)
- mode (optional) - mode to run the bot.
  - keeper (default) - keep all stored links until they are seen and deleted with /clear command
  - scavenger - all seen links are deleted IMMEDIATELY
- host (optional) - host to collect bot updates, default is api.telegram.org
## Run locally
In project directory
```shell
go run ./cmd/bot/main.go -token "$token"
```

or

```shell
go build -o ./bin/main ./cmd/bot/main.go && ./bin/main -token "$token"
```

## Functionality

The list of functions, provided by bot:
- **Store links**

    Send bot any message. It will parse it's content and extract all given links, then store them for you


- **Pick random link**

    Returns random link from unseen. Is scavenger mode is on, deletes picked link


- **Show list**

    Show list of all stored links. Links that have been seen are marked with 👀


- **Clear list**

    Remove all links, that have been already seen

## Run In Docker

To build image
```shell
docker build -t storage_bot .
```
To run image
```shell
docker run --name storage_bot -v storage_bot_vol:/root/data -d storage_bot -token "$token" 
```

## Thanks

Basic bot architecture [inspired by this playlist](https://www.youtube.com/playlist?list=PLFAQFisfyqlWDwouVTUztKX2wUjYQ4T3l)

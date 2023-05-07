# Peristera - Telegram Bot for ORGs

This bot provide some tooks to interact with ORGs like christian churchs.

## Authors

- [@vinicius73](https://www.github.com/vinicius73)

## Development

### Requeriments

- [Go ~> 1.20](https://go.dev/dl/)
- [Task v3](https://taskfile.dev/)
- [upx](https://upx.github.io/)

### Workflow

This project uses [taskfile](https://taskfile.dev/) to automate actions and commands.

```sh
task setup  # install dependencies
task build  # build project
task run -- --help   # run cli with an argument "--help"
```

## CLI Commands

### Root arguments

The root arguments allow to define the main behaviours and configurations.

```sh
# define config location
peristera --config /path/to/config/config.yaml

# output debug information
peristera --debug

# hide the main banner
peristera --no-banner

# define log level
peristera --level error
```

### System command

The system command is usefull for small interactions about the host whom the bot will run.

```sh
# show host info
peristera system info

# instead of show the host info, send it to root users
peristera --no-banner system info --notify "Ya-Ha!"

# do a backup of interal database and send it to root users
peristera --no-banner system backup --force --notify yo
```

### Worker command

The `worker` command starts the telegram bot pooling. Its allows the bot to receive and handle the messages and commands.

```sh
# start the bot worker and cron actions
peristera --no-banner -c /peristera.yaml worker --cron
```

## Architecture and structure

Peristera bot allows you to customize the commands and output, but each one is restricted to some objectives.

> This customization is build on the [config file](./peristera.example.yml), check the example for more details.

### Handlers

There are 5 main handlers

#### `start`

`start` is the main handler when someone start to interact with Peristera. Excentialy it output the menu and the `description` from config file.

#### `pix`

The `pix` handler outputs a QRCode to send a PIX, also some extra informations based on the config file.

#### `videos`

The `videos` handler fetch the latest videos from `youtube` channel.

> You must have a [Youtube API Key](https://developers.google.com/youtube/registering_an_application) to feth this data.

#### `address`

The `address` handler just send a location messages based on the config `location` option.

#### `calendar`

The `calendar` handler use the content of `calendar` to send a message.

## Internal / Restricted bot commands

There are some bot commands whom are restricted to admins or root users.

Those commands are designed to do some administrative actions and special actions.

### Admin commands

#### `/me`

Collect data from the current user or the user from a message.

It's useful to retrieve the user ID for posterior usage.

#### `/system`

It will send a message with current data from the host whom the bot is running.

### `/cover`

Cover is a special tool to generate images from a input text.

This feature is based on [`diakonos` cover generate](https://github.com/comunidade-shallom/diakonos#diakonos-cover) feature.

> Check `cover` config to knows how to adjust some options.

```
/cover 1920x1080 Generates a imagem from given size and text input
```

```
/cover <WIDTH>x<HEIGTH> <INPUT TEXT>
```

You can also reply a previous generate image to create a new one.

## Root commands

### `/exec`

Execute a command in the host system

```
/exec free -h
/exec df -h
```

### `/backup`

Do a backup from the database file.

### `/load`

It is possible to load a backup file from a message when you answer a backup file with this command.

# Telegram bot

This repository contains the code for a Telegram bot written in Golang. The bot is designed to 
provide a simple interface for users to interact with, and can be easily extended with additional 
functionality.

## Installation
To install the bot, you will need to have Golang 1.20 installed on your system. You can download 
Golang from the official website: https://go.dev/dl/

Once you have Golang installed, you can clone the repository and build the bot:

```bash
git clone https://github.com/bonefabric/telegram-adviser-bot.git
cd telegram-adviser-bot
go build bonefabric/adviser/cmd/adviser
```
This will generate an executable file called telegram-bot that you can run to start the bot.

## Configuration
Before you can start using the bot, you will need to create a new bot and obtain an API token 
from Telegram. To do this, follow these steps:

- Open Telegram and search for the BotFather bot.
- Follow the on-screen instructions to create a new bot and obtain an API token.
- Copy the API token and paste it into the config.yaml file in the repository.

You can also configure other settings in the __config.yaml__ file, such as data store, credentials 
and other. For example repository contains __config.example.yaml__ file.

## Usage
To start the bot, simply run the executable file:

```bash
./adviser
```
The bot will connect to the Telegram API and start listening for incoming messages. 
You can send messages to the bot by searching for it on Telegram and sending a message.

By default, the bot will respond with help with the commands

## Contributing
If you would like to contribute to the project, please fork the repository and create a new branch 
for your changes. Once you have made your changes, submit a pull request and I will review your changes.

# Telegram notifier bot for GitHub releases
Bot check new releases and tags from user subscriptions and send notification to Telegram, like this:
> bot_name, [27.10.20 09:59]
> redis: 5.0.10 at 10-27-2020: https://github.com/redis/redis

## Getting started
* Generate personal GitHub token: https://github.com/settings/tokens
* Create new Telegram bot https://telegram.me/BotFather
* Build project "go build"
* Create two token files with tokens: github_token and telegram_token in same directory of binary location
* Run releasebot

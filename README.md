# taskuji-slackbot-socketmode
Mention to bot, bot select user randomly in the channel.  
Typical usage is when asking someone for a task.  
It uses [Socket Mode](https://api.slack.com/apis/connections/socket-implement)

## Usage in Slack
Mention bot user name then bot select a member in the channel randomly.
Press `OK` button if the member can do the task, otherwise press `NG` buttton.
When `NG` button is pressed, bot select a member randomly again.

## Preparation in Slack
-  [Create Slack App as internal integrations](https://api.slack.com/internal-integrations)
-  Turn on `Bots`
-  Turn on `Socket Mode`
-  Turn on `Interactivity`
-  Turn on `Event Subscriptions`
    - Add `Subscribe to bot events` : `app_mention` 
-  Add Permissions
    - channels:read
    - chat:write
-  Invite bot user to your channel.

## Run
To run this bot, you need to set the following env vars,

```bash
export APP_LEVEL_TOKEN="xapp-***"      // you can get this after turn on socketmode (via slack app management console)
export BOT_TOKEN="xoxb-***"      // you can get this after create a bot user (via slack app management console)
```

To run this, 

```bash
$ dep ensure
$ go build -o bot && ./bot
```
## Appedix
Reference:
- https://api.slack.com/apis/connections/socket-implement
- https://qiita.com/seratch/items/c7d9aeb60ead5c126c01
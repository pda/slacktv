slacktv
=======

slacktv is a [Slack bot] designed to run on a Mac plugged into a digital signage TV. When mentioned in a message with a URL, that URL is opened.


Usage
-----

* Create a Slack bot integration user at https://99designs.slack.com/services/new/bot
* Invite the bot to the desired channel.
* Run `SLACK_TOKEN=... slacktv`.
* Mention the bot's username along with a URL.

```sh
export SLACK_TOKEN="slack-bot-auth-token-here"
slacktv
```


Development
------------

slacktv is written in [Go], uses [godep] to manage dependencies, and [gorilla/websocket] for [Slack's Real Time Messaging API][RTM].

The `slack` subpackage is the beginnings of a generic Slack API/RTM client.

```sh
export SLACK_TOKEN=...

godep restore
go test ./...
go build ./cmd/slacktv
DEBUG=1 ./slacktv  # verbose logging to stdout

# restart on every source file change
go get github.com/skelterjohn/rerun
DEBUG=1 GOPATH=`godep path`:$GOPATH rerun github.com/99designs/slacktv/cmd/slacktv
```


[Slack bot]: https://api.slack.com/bot-users
[Go]: https://golang.org/
[godep]: https://github.com/tools/godep
[gorilla/websocket]: https://github.com/gorilla/websocket
[RTM]: https://api.slack.com/rtm

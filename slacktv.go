package slacktv

import (
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"

	"github.com/pda/slacktv/slack"
)

var (
	urlRegexp *regexp.Regexp
	tagRegexp *regexp.Regexp
)

func init() {
	urlRegexp = regexp.MustCompile(`<(https?://.+?)>`)
	tagRegexp = regexp.MustCompile(`<.*?>`)
}

func Run() {
	dbg("initializing RTM session")
	rtm, err := slack.RTMStart()
	if err != nil {
		panic(err)
	}

	dbg("listening for events")
	for ev := range rtm.Events {
		handleEvent(rtm, ev)
	}
}

func handleEvent(rtm *slack.RTMSession, ev slack.Event) {
	dbg("event: %#v", ev)
	switch ev["type"] {
	case "message":
		handleMessage(rtm, ev)
	}
}

func handleMessage(rtm *slack.RTMSession, ev slack.Event) {
	if ev["user"] == nil || ev["text"] == nil {
		return
	}

	user := rtm.User(ev["user"].(string))
	text := ev["text"].(string)
	chanId := ev["channel"].(string)

	if isMention(text, rtm.Self.Id) {
		if isGreeting(text) {
			err := chatPostMessage(chanId, fmt.Sprintf("hello <@%s>!", user.Id))
			if err != nil {
				log.Print(err)
			}
		}

		if url, ok := urlFromMessageText(text); ok {
			err := chatPostMessage(
				chanId,
				fmt.Sprintf("<@%s> changing the channel to %s", user.Id, url),
			)
			if err != nil {
				log.Print(err)
			}

			dbg("URL from %s: %s", user.Name, url)
			if err := open(url); err != nil {
				log.Print(err)
			}
		}
	}
}

func isMention(text, userId string) bool {
	return strings.Contains(text, fmt.Sprintf("<@%s>", userId))
}

func urlFromMessageText(text string) (string, bool) {
	m := urlRegexp.FindSubmatch([]byte(text))
	if m == nil {
		return "", false
	}
	return string(m[1]), true
}

func isGreeting(text string) bool {
	stripped := tagRegexp.ReplaceAllString(text, "")
	return strings.Contains(stripped, "hello")
}

func chatPostMessage(chanId, text string) (err error) {
	args := url.Values{
		"channel":      {chanId},
		"text":         {text},
		"as_user":      {"true"},
		"unfurl_links": {"false"},
	}
	dbg("chat.postMessage: %#v", args)
	call, err := slack.Method("chat.postMessage", args)
	if err != nil {
		log.Print(err)
	}
	var resp slack.ChatPostMessageResponse
	err = call.Exec(&resp)
	if err != nil {
		log.Print(err)
	}
	if !resp.Ok {
		log.Printf("chatPostMessage: %#v\n", resp)
	}
	return
}

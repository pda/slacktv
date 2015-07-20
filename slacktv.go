package slacktv

import (
	"fmt"
	"log"
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

	if isMention(text, rtm.Self.Id) {
		if url, ok := urlFromMessageText(text); ok {
			dbg("URL from %s: %s", user.Name, url)
			err := open(url)
			if err != nil {
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

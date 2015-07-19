package slacktv

import (
	"log"
	"regexp"

	"github.com/pda/slacktv/slack"
)

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

	re := regexp.MustCompile(`\A\s*<@(.+?)>.*<(.+?)>`)
	sm := re.FindSubmatch([]byte(text))
	if sm != nil && string(sm[1]) == rtm.Self.Id {
		url := string(sm[2])
		dbg("URL from %s: %s", user.Name, url)
		err := open(url)
		if err != nil {
			log.Print(err)
		}
	}
}

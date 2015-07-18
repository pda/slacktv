package slacktv

import (
	"fmt"
	"os"
	"regexp"
)

func Run() {
	token := mustGetToken()
	sess, err := Connect(token)
	if err != nil {
		panic(err)
	}

	for event := range sess.Events {
		handleEvent(sess, event)
	}
}

func mustGetToken() string {
	token := os.Getenv("AUTH_TOKEN")
	if token == "" {
		panic("mustGetToken() requires AUTH_TOKEN")
	}
	return token
}

func handleEvent(s *Session, ev Event) {
	if ev["type"] == "message" {
		handleMessage(s, ev)
	}
}

func handleMessage(s *Session, ev Event) {
	if ev["user"] == nil || ev["text"] == nil || ev["channel"] == nil {
		return
	}

	user := s.User(ev["user"].(string))
	channel := s.Channel(ev["channel"].(string))
	text := ev["text"].(string)

	fmt.Printf("#%s | <%s> %s\n", channel.Name, user.Name, text)

	re := regexp.MustCompile(`\A\s*<@(.+?)>.*<(.+?)>`)
	sm := re.FindSubmatch([]byte(text))
	if sm != nil && string(sm[1]) == s.Self.Id {
		url := sm[2]
		fmt.Printf("-- URL from %s: %s\n", user.Name, url)
	}
}

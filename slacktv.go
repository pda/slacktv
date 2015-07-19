package slacktv

import (
	"bytes"
	"os/exec"
	"regexp"

	"github.com/pda/slacktv/slack"
)

func Run() {
	dbg("initializing RTM session")
	rtm, err := slack.RTMStart()
	if err != nil {
		panic(err)
	}

	browser := make(chan string)
	go func() {
		for {
			url := <-browser
			dbg("opening %s", url)
			cmdStdout := &bytes.Buffer{}
			cmdStderr := &bytes.Buffer{}
			cmd := exec.Command("/usr/bin/open", "-n", "-a", "Google Chrome", "--args", "--kiosk", url)
			cmd.Stdout = cmdStdout
			cmd.Stderr = cmdStderr
			err = cmd.Start()
			if err != nil {
				panic(err)
			}
			err = cmd.Wait()
			if err != nil {
				dbg("cmd exit: %#v\nstdout: %s\nstderr: %s", err, cmdStdout, cmdStderr)
			}
		}
	}()

	dbg("listening for events")
	for ev := range rtm.Events {
		handleEvent(rtm, ev, browser)
	}
}

func handleEvent(rtm *slack.RTMSession, ev slack.Event, browser chan string) {
	dbg("event: %#v", ev)
	switch ev["type"] {
	case "message":
		handleMessage(rtm, ev, browser)
	}
}

func handleMessage(rtm *slack.RTMSession, ev slack.Event, browser chan string) {
	if ev["user"] == nil || ev["text"] == nil {
		return
	}

	user := rtm.User(ev["user"].(string))
	text := ev["text"].(string)

	re := regexp.MustCompile(`\A\s*<@(.+?)>.*<(.+?)>`)
	sm := re.FindSubmatch([]byte(text))
	if sm != nil && string(sm[1]) == rtm.Self.Id {
		url := string(sm[2])
		dbg("-- URL from %s: %s", user.Name, url)
		browser <- url
	}
}

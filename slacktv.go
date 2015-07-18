package slacktv

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
)

func Run() {
	token := mustGetToken()
	sess, err := Connect(token)
	if err != nil {
		panic(err)
	}

	browser := make(chan string)
	go func() {
		for {
			url := <-browser
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
				fmt.Printf("exit: %#v\nstdout: %s\nstderr: %s\n", err, cmdStdout, cmdStderr)
			}
		}
	}()

	for event := range sess.Events {
		handleEvent(sess, event, browser)
	}
}

func mustGetToken() string {
	token := os.Getenv("AUTH_TOKEN")
	if token == "" {
		panic("mustGetToken() requires AUTH_TOKEN")
	}
	return token
}

func handleEvent(s *Session, ev Event, browser chan string) {
	switch ev["type"] {
	case "message":
		handleMessage(s, ev, browser)
	}
}

func handleMessage(s *Session, ev Event, browser chan string) {
	if ev["user"] == nil || ev["text"] == nil {
		return
	}

	user := s.User(ev["user"].(string))
	text := ev["text"].(string)

	re := regexp.MustCompile(`\A\s*<@(.+?)>.*<(.+?)>`)
	sm := re.FindSubmatch([]byte(text))
	if sm != nil && string(sm[1]) == s.Self.Id {
		url := string(sm[2])
		fmt.Printf("-- URL from %s: %s\n", user.Name, url)
		browser <- url
	}
}

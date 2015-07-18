package slacktv

import (
	"fmt"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

func Run() {
	client := &SlackClient{Token: mustGetToken()}
	testAuth(client)
	rtm(client)
}

func mustGetToken() string {
	token := os.Getenv("AUTH_TOKEN")
	if token == "" {
		panic("mustGetToken() requires AUTH_TOKEN")
	}
	return token
}

func testAuth(client *SlackClient) {
	resp, err := client.AuthTest()
	if err != nil {
		panic(err)
	}

	if resp.Ok {
		fmt.Println("authentication successful")
	} else {
		fmt.Println("authentication failed:", resp.Error)
	}
}

func rtm(client *SlackClient) {
	resp, err := client.RtmStart()
	if err != nil {
		panic(err)
	}

	var chanNames []string
	for _, c := range resp.Channels {
		if c.IsMember {
			chanNames = append(chanNames, fmt.Sprintf("#%s", c.Name))
		}
	}

	fmt.Printf("Self: %#v\n", resp.Self)
	fmt.Printf("Team: %#v\n", resp.Team)
	fmt.Printf("Channels: %s", strings.Join(chanNames, " "))
	fmt.Println()

	conn, wsResp, err := websocket.DefaultDialer.Dial(resp.Url, nil)
	if err != nil {
		panic(err)
	}
	_ = wsResp

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			return
		}
		fmt.Println(messageType, string(p))
	}

}

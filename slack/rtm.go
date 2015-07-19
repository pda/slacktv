package slack

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

type RTMSession struct {
	Events chan Event
	Self   *SelfResource
	Team   *TeamResource
	Url    string
	Users  []UserResource

	conn    *websocket.Conn
	userMap map[string]UserResource
}

func (rtm *RTMSession) User(id string) *UserResource {
	u := rtm.userMap[id]
	if u.Id == "" {
		u = UserResource{Id: id, Name: fmt.Sprintf("(id:%s)", u.Id)}
	}
	return &u
}

type Event map[string]interface{}

func RTMStart() (rtm *RTMSession, err error) {
	resp, err := rtmStartResponse()
	if err != nil {
		return
	}

	rtm = rtmSessionFromResponse(resp)

	conn, _, err := websocket.DefaultDialer.Dial(rtm.Url, nil)
	if err != nil {
		return nil, err
	}
	rtm.conn = conn

	go rtmEmitEvents(rtm)

	return
}

func rtmStartResponse() (resp *RTMStartResponse, err error) {
	call, err := Method("rtm.start", url.Values{})
	if err != nil {
		return
	}
	err = call.Exec(&resp)
	if err != nil {
		return
	}
	if !resp.Ok {
		return nil, fmt.Errorf("%#v", resp)
	}
	return
}

func rtmSessionFromResponse(resp *RTMStartResponse) (rtm *RTMSession) {
	rtm = &RTMSession{
		Url:    resp.Url,
		Self:   resp.Self,
		Team:   resp.Team,
		Users:  resp.Users,
		Events: make(chan Event),
	}
	mapUsers(rtm)
	return
}

func mapUsers(rtm *RTMSession) {
	rtm.userMap = make(map[string]UserResource)
	for _, u := range rtm.Users {
		rtm.userMap[u.Id] = u
	}
}

func rtmEmitEvents(rtm *RTMSession) {
	for {
		messageType, p, err := rtm.conn.ReadMessage()
		if err != nil {
			log.Print(err)
			continue
		}
		if messageType == websocket.TextMessage {
			var ev Event
			err = json.Unmarshal(p, &ev)
			if err != nil {
				log.Print(err)
				continue
			}
			rtm.Events <- ev
		}
	}
}

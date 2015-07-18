package slacktv

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/websocket"
)

type Session struct {
	Token   string
	Self    *SelfResource
	Team    *TeamResource
	Users   []UserResource
	UserMap map[string]UserResource
	Events  chan Event
}

func (s *Session) User(id string) *UserResource {
	u := s.UserMap[id]
	if u.Id == "" {
		u = UserResource{Id: id, Name: fmt.Sprintf("(id:%s)", u.Id)}
	}
	return &u
}

type RtmStartResponse struct {
	Ok    bool           `json:"ok"`
	Url   string         `json:"url"`
	Self  *SelfResource  `json:"self"`
	Team  *TeamResource  `json:"team"`
	Users []UserResource `json:"users"`
}

type SelfResource struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type TeamResource struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	EmailDomain string `json:"email_domain"`
	Domain      string `json:"domain"`
}

type UserResource struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Event map[string]interface{}

func Connect(token string) (*Session, error) {
	resp, err := rtmInit(token)
	if err != nil {
		return nil, err
	}

	session := &Session{
		Token:  token,
		Self:   resp.Self,
		Team:   resp.Team,
		Users:  resp.Users,
		Events: make(chan Event),
	}

	mapUsers(session)

	conn, _, err := websocket.DefaultDialer.Dial(resp.Url, nil)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				return
			}
			if messageType == websocket.TextMessage {
				var ev Event
				err = json.Unmarshal(p, &ev)
				if err != nil {
					panic(err)
				}
				session.Events <- ev
			}
		}
	}()

	return session, nil
}

func rtmInit(token string) (resp *RtmStartResponse, err error) {
	url := rtmUrl(token)
	httpResp, err := http.Get(url)
	if err != nil {
		return
	}
	defer httpResp.Body.Close()
	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return
	}
	return
}

func rtmUrl(token string) string {
	var rtmUrl = "https://slack.com/api/rtm.start"
	return fmt.Sprintf("%s?token=%s", rtmUrl, token)
}

func mapUsers(s *Session) {
	s.UserMap = make(map[string]UserResource)
	for _, u := range s.Users {
		s.UserMap[u.Id] = u
	}
}

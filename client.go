package slacktv

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type SlackClient struct {
	Token string
}

func (*SlackClient) AuthTest() (resp AuthTestResponse, err error) {
	url := authUrl(mustGetToken())
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

func authUrl(token string) string {
	var authUrl = "https://slack.com/api/auth.test"
	return fmt.Sprintf("%s?token=%s", authUrl, token)
}

type AuthTestResponse struct {
	Ok    bool
	Error string
}

func (*SlackClient) RtmStart() (resp RtmStartResponse, err error) {
	url := rtmUrl(mustGetToken())
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

type RtmStartResponse struct {
	Ok       bool              `json:"ok"`
	Url      string            `json:"url"`
	Self     SelfResource      `json:"self"`
	Team     TeamResource      `json:"team"`
	Users    []UserResource    `json:"users"`
	Channels []ChannelResource `json:"channels"`
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

type ChannelResource struct {
	Id         string   `json:"id"`
	Name       string   `json:"name"`
	IsArchived bool     `json:"is_archived"`
	IsGeneral  bool     `json:"is_general"`
	Members    []string `json:"members"`
	IsMember   bool     `json:"is_member"`
}

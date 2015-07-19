package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type MethodCall struct {
	Url         string
	ContentType string
	Body        []byte
}

type Data map[string]interface{}

func (m *MethodCall) Exec(resp interface{}) (err error) {
	httpResp, err := http.Post(m.Url, m.ContentType, bytes.NewReader(m.Body))
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

func url(method string) string {
	base := os.Getenv("SLACK_URL")
	if base == "" {
		base = "https://slack.com/api"
	}
	return fmt.Sprint(base, "/", method)
}

func Method(name string, d Data) (mc *MethodCall, err error) {
	d["token"] = mustGetToken()
	json, err := json.Marshal(d)
	if err != nil {
		return
	}
	mc = &MethodCall{
		Url:         url(name),
		ContentType: "application/json",
		Body:        json,
	}
	return
}

func mustGetToken() string {
	token := os.Getenv("SLACK_TOKEN")
	if token == "" {
		panic("mustGetToken() requires SLACK_TOKEN")
	}
	return token
}

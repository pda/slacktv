package slack

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

type MethodCall struct {
	Url  string
	Data url.Values
}

func (m *MethodCall) Exec(resp interface{}) (err error) {
	httpResp, err := http.PostForm(m.Url, m.Data)
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

func apiUrl(method string) string {
	base := os.Getenv("SLACK_URL")
	if base == "" {
		base = "https://slack.com/api"
	}
	return fmt.Sprint(base, "/", method)
}

func Method(name string, d url.Values) (mc *MethodCall, err error) {
	d.Set("token", mustGetToken())
	mc = &MethodCall{
		Url:  apiUrl(name),
		Data: d,
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

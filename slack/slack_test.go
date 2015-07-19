package slack_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/pda/slacktv/slack"
)

func TestMain(m *testing.M) {
	os.Setenv("SLACK_TOKEN", "secret")
	os.Exit(m.Run())
}

func TestMethod(t *testing.T) {
	call, err := slack.Method("channel.join", slack.Data{"name": "test"})

	if err != nil {
		t.Error(err)
	}
	if call.Url != "https://slack.com/api/channel.join" {
		t.Fail()
	}
	if call.ContentType != "application/json" {
		t.Fail()
	}
	var d slack.Data
	err = json.Unmarshal(call.Body, &d)
	if d["token"] != "secret" {
		t.Fail()
	}
	if d["name"] != "test" {
		t.Fail()
	}
}

func TestApiCall(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"ok": true, "channel": {"id": "xyz", "name": "test"}}`)
	}))
	defer ts.Close()

	oldUrl := os.Getenv("SLACK_URL")
	os.Setenv("SLACK_URL", ts.URL)
	defer os.Setenv("SLACK_URL", oldUrl)

	call, err := slack.Method("channel.join", slack.Data{"name": "test"})
	if err != nil {
		t.Error(err)
	}
	var resp slack.ChannelResponse
	err = call.Exec(&resp)
	if err != nil {
		t.Error(err)
	}
	if resp.Ok != true {
		t.Logf("resp: %#v", resp)
		t.Fail()
	}
	if resp.Channel.Name != "test" || resp.Channel.Id != "xyz" {
		t.Logf(`resp.Channel: %#v`, resp.Channel)
		t.Fail()
	}
}

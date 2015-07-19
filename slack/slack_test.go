package slack_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/pda/slacktv/slack"
)

func TestMain(m *testing.M) {
	os.Setenv("SLACK_TOKEN", "secret")
	os.Exit(m.Run())
}

func TestMethod(t *testing.T) {
	call, err := slack.Method("channel.join", url.Values{"name": {"test"}})

	if err != nil {
		t.Error(err)
	}
	if call.Url != "https://slack.com/api/channel.join" {
		t.Error(call.Url)
	}
	if call.Data.Get("token") != "secret" {
		t.Errorf("%#v", call.Data)
	}
	if call.Data.Get("name") != "test" {
		t.Errorf("%#v", call.Data)
	}
}

func TestApiCall(t *testing.T) {
	var req *http.Request

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req = r
		fmt.Fprintln(w, `{"ok": true, "channel": {"id": "xyz", "name": "test"}}`)
	}))
	defer ts.Close()

	oldUrl := os.Getenv("SLACK_URL")
	os.Setenv("SLACK_URL", ts.URL)
	defer os.Setenv("SLACK_URL", oldUrl)

	call, err := slack.Method("channel.join", url.Values{"name": {"test"}})
	if err != nil {
		t.Error(err)
	}
	var resp slack.ChannelResponse
	err = call.Exec(&resp)
	if err != nil {
		t.Error(err)
	}

	ct := req.Header.Get("Content-Type")
	if ct != "application/x-www-form-urlencoded" {
		t.Error(`Content-Type:`, ct)
	}

	if resp.Ok != true {
		t.Errorf("resp: %#v", resp)
	}

	if resp.Channel.Name != "test" || resp.Channel.Id != "xyz" {
		t.Errorf(`resp.Channel: %#v`, resp.Channel)
	}
}

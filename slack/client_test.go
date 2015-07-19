package slack

import "testing"

func TestUrl(t *testing.T) {
	if url("api.test") != "https://slack.com/api/api.test" {
		t.Fail()
	}
}

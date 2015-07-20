package slacktv

import "testing"

func TestIsMentionSingleMention(t *testing.T) {
	if !isMention("<@U02QFPOAN> hi", "U02QFPOAN") {
		t.Fail()
	}
}

func TestIsMentionMultipleMentions(t *testing.T) {
	if !isMention("ping <@U02QFPOAN> <@U05JWPFMA>", "U02QFPOAN") {
		t.Fail()
	}
}

func TestIsMentionNope(t *testing.T) {
	if isMention("ping <@U08PQHGUD> <@U05JWPFMA>", "U02QFPOAN") {
		t.Fail()
	}
}

func TestUrlFromText(t *testing.T) {
	type result struct {
		url string
		ok  bool
	}

	table := map[string]result{
		"hello":                                          {"", false},
		"<@U07RFLAPN>: hello":                            {"", false},
		"https://example.org/":                           {"", false},
		"<https://example.org/>":                         {"https://example.org/", true},
		"<@U07RFLAPN>: <https://example.org/>":           {"https://example.org/", true},
		"<http://example.org/1> <https://example.org/2>": {"http://example.org/1", true},
	}

	for text, expected := range table {
		t.Logf("text: %#v", text)
		url, ok := urlFromMessageText(text)
		if ok != expected.ok {
			t.Errorf("expected ok = %#v, got %#v:", expected.ok, ok)
		}
		if url != expected.url {
			t.Errorf("expected url = %#v, got %#v:", expected.url, url)
		}
	}

}

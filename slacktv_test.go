package slacktv

import "testing"

func TestIsMention(t *testing.T) {
	type expectation struct {
		text   string
		userId string
		result bool
	}

	table := []expectation{
		{text: "<@U02QFPOAN> hi", userId: "U02QFPOAN", result: true},
		{text: "ping <@U02QFPOAN> <@U05JWPFMA>", userId: "U02QFPOAN", result: true},
		{text: "ping <@U08PQHGUD> <@U05JWPFMA>", userId: "U02QFPOAN", result: false},
	}

	for _, exp := range table {
		if isMention(exp.text, exp.userId) != exp.result {
			t.Errorf("%#v", exp)
		}
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

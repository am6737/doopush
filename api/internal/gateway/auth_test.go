package gateway

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestParseHandshakeParams_OK(t *testing.T) {
	r := httptest.NewRequest("GET", "/ws?appid=42&installation_token=dp_ins_test", nil)
	p, err := parseHandshakeParams(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.AppID != 42 || p.InstallationToken != "dp_ins_test" {
		t.Fatalf("got %+v", p)
	}
}

func TestParseHandshakeParams_Missing(t *testing.T) {
	cases := []url.Values{
		{"installation_token": {"dp_ins_test"}},
		{"appid": {"1"}},
		{"appid": {"1"}, "appkey": {"x"}, "token": {"y"}},
		{"appid": {"abc"}, "installation_token": {"dp_ins_test"}},
	}
	for i, q := range cases {
		r := httptest.NewRequest("GET", "/ws?"+q.Encode(), nil)
		if _, err := parseHandshakeParams(r); err == nil {
			t.Fatalf("case %d expected error, got nil", i)
		}
	}
}

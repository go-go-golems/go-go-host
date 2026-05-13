package cmds

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestPollDeviceTokenPendingSlowDownThenSuccess(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Fatalf("parse form: %v", err)
		}
		if r.Form.Get("grant_type") != deviceCodeGrantType || r.Form.Get("client_id") != "go-go-host-cli" || r.Form.Get("device_code") != "dev_123" {
			t.Fatalf("unexpected form: %v", r.Form)
		}
		attempts++
		switch attempts {
		case 1:
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(oauthErrorResponse{Error: "authorization_pending"})
		case 2:
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(oauthErrorResponse{Error: "slow_down"})
		default:
			_ = json.NewEncoder(w).Encode(tokenResponse{AccessToken: "access", RefreshToken: "refresh", TokenType: "Bearer", ExpiresIn: 300})
		}
	}))
	defer server.Close()

	var sleeps []time.Duration
	sleeper := func(ctx context.Context, d time.Duration) error {
		sleeps = append(sleeps, d)
		return nil
	}
	tok, err := pollDeviceTokenWithSleeper(context.Background(), server.URL, "go-go-host-cli", deviceAuthorizationResponse{DeviceCode: "dev_123", Interval: 1, ExpiresIn: 600}, sleeper)
	if err != nil {
		t.Fatalf("poll: %v", err)
	}
	if tok.AccessToken != "access" || attempts != 3 {
		t.Fatalf("unexpected token/attempts: %#v attempts=%d", tok, attempts)
	}
	if len(sleeps) != 3 || sleeps[0] != time.Second || sleeps[1] != time.Second || sleeps[2] != 6*time.Second {
		t.Fatalf("unexpected sleeps: %#v", sleeps)
	}
}

func TestPollDeviceTokenStopsOnDenied(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(oauthErrorResponse{Error: "access_denied", ErrorDescription: "no"})
	}))
	defer server.Close()
	sleeper := func(ctx context.Context, d time.Duration) error { return nil }
	_, err := pollDeviceTokenWithSleeper(context.Background(), server.URL, "go-go-host-cli", deviceAuthorizationResponse{DeviceCode: "dev_123", Interval: 1, ExpiresIn: 600}, sleeper)
	if err == nil {
		t.Fatalf("expected denial error")
	}
}

func TestScopesFromString(t *testing.T) {
	got := scopesFromString("openid,profile email openid", []string{"fallback"})
	want := []string{"openid", "profile", "email"}
	if len(got) != len(want) {
		t.Fatalf("unexpected scopes: %#v", got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("unexpected scopes: %#v", got)
		}
	}
	fallback := []string{"openid"}
	got = scopesFromString("", fallback)
	if len(got) != 1 || got[0] != "openid" {
		t.Fatalf("expected fallback, got %#v", got)
	}
}

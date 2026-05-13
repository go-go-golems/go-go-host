package cmds

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const deviceCodeGrantType = "urn:ietf:params:oauth:grant-type:device_code"

var errAuthorizationPending = errors.New("authorization pending")

type publicConfigResponse struct {
	DevAuth bool              `json:"devAuth"`
	OIDC    *publicOIDCConfig `json:"oidc"`
}

type publicOIDCConfig struct {
	Issuer             string   `json:"issuer"`
	ClientID           string   `json:"clientId"`
	DeviceClientID     string   `json:"deviceClientId"`
	Scopes             []string `json:"scopes"`
	RedirectPath       string   `json:"redirectPath"`
	LogoutRedirectPath string   `json:"logoutRedirectPath"`
}

type oidcDiscovery struct {
	Issuer                      string   `json:"issuer"`
	TokenEndpoint               string   `json:"token_endpoint"`
	DeviceAuthorizationEndpoint string   `json:"device_authorization_endpoint"`
	RevocationEndpoint          string   `json:"revocation_endpoint"`
	GrantTypesSupported         []string `json:"grant_types_supported"`
}

type deviceAuthorizationResponse struct {
	DeviceCode              string `json:"device_code"`
	UserCode                string `json:"user_code"`
	VerificationURI         string `json:"verification_uri"`
	VerificationURIComplete string `json:"verification_uri_complete"`
	ExpiresIn               int    `json:"expires_in"`
	Interval                int    `json:"interval"`
}

type tokenResponse struct {
	AccessToken      string `json:"access_token"`
	IDToken          string `json:"id_token"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	Scope            string `json:"scope"`
}

type oauthErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func (e oauthErrorResponse) ErrorString() string {
	if e.ErrorDescription != "" {
		return e.Error + ": " + e.ErrorDescription
	}
	return e.Error
}

type pollSleeper func(context.Context, time.Duration) error

func realPollSleeper(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}

func fetchPublicConfig(ctx context.Context, apiURL string) (publicConfigResponse, error) {
	var cfg publicConfigResponse
	if err := getJSON(strings.TrimRight(apiURL, "/"), "/api/v1/config", &cfg); err != nil {
		return cfg, err
	}
	if cfg.OIDC == nil || strings.TrimSpace(cfg.OIDC.Issuer) == "" {
		return cfg, fmt.Errorf("server %s did not publish OIDC configuration; use --dev-user for dev auth or --bearer-token for manual auth", apiURL)
	}
	return cfg, nil
}

func discoverOIDC(ctx context.Context, issuer string) (oidcDiscovery, error) {
	endpoint := strings.TrimRight(issuer, "/") + "/.well-known/openid-configuration"
	var doc oidcDiscovery
	if err := getJSONAbsolute(ctx, endpoint, &doc); err != nil {
		return doc, fmt.Errorf("discover OIDC provider: %w", err)
	}
	if doc.TokenEndpoint == "" {
		return doc, fmt.Errorf("OIDC discovery document does not include token_endpoint")
	}
	if doc.DeviceAuthorizationEndpoint == "" {
		return doc, fmt.Errorf("OIDC discovery document does not include device_authorization_endpoint")
	}
	return doc, nil
}

func getJSONAbsolute(ctx context.Context, endpoint string, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return httpResponseError("GET", endpoint, resp)
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

func startDeviceAuthorization(ctx context.Context, endpoint, clientID string, scopes []string) (deviceAuthorizationResponse, error) {
	form := url.Values{}
	form.Set("client_id", clientID)
	if len(scopes) > 0 {
		form.Set("scope", strings.Join(scopes, " "))
	}
	var out deviceAuthorizationResponse
	if err := postForm(ctx, endpoint, form, &out); err != nil {
		return out, err
	}
	if out.Interval <= 0 {
		out.Interval = 5
	}
	if out.ExpiresIn <= 0 {
		out.ExpiresIn = 600
	}
	if out.DeviceCode == "" || out.UserCode == "" || out.VerificationURI == "" {
		return out, fmt.Errorf("device authorization response is missing device_code, user_code, or verification_uri")
	}
	return out, nil
}

func pollDeviceToken(ctx context.Context, tokenEndpoint, clientID string, device deviceAuthorizationResponse) (tokenResponse, error) {
	return pollDeviceTokenWithSleeper(ctx, tokenEndpoint, clientID, device, realPollSleeper)
}

func pollDeviceTokenWithSleeper(ctx context.Context, tokenEndpoint, clientID string, device deviceAuthorizationResponse, sleep pollSleeper) (tokenResponse, error) {
	interval := time.Duration(device.Interval) * time.Second
	if interval <= 0 {
		interval = 5 * time.Second
	}
	deadline := time.Now().Add(time.Duration(device.ExpiresIn) * time.Second)
	if device.ExpiresIn <= 0 {
		deadline = time.Now().Add(10 * time.Minute)
	}
	for {
		if time.Now().After(deadline) {
			return tokenResponse{}, fmt.Errorf("device authorization expired before approval; run go-go-host login again")
		}
		if err := sleep(ctx, interval); err != nil {
			return tokenResponse{}, err
		}
		tok, oauthErr, err := requestDeviceToken(ctx, tokenEndpoint, clientID, device.DeviceCode)
		if err != nil {
			interval *= 2
			if interval > 30*time.Second {
				interval = 30 * time.Second
			}
			continue
		}
		if oauthErr == nil {
			return tok, nil
		}
		switch oauthErr.Error {
		case "authorization_pending":
			continue
		case "slow_down":
			interval += 5 * time.Second
			continue
		case "access_denied":
			return tokenResponse{}, fmt.Errorf("device authorization was denied in the browser")
		case "expired_token":
			return tokenResponse{}, fmt.Errorf("device code expired; run go-go-host login again")
		default:
			return tokenResponse{}, fmt.Errorf("device token polling failed: %s", oauthErr.ErrorString())
		}
	}
}

func requestDeviceToken(ctx context.Context, tokenEndpoint, clientID, deviceCode string) (tokenResponse, *oauthErrorResponse, error) {
	form := url.Values{}
	form.Set("grant_type", deviceCodeGrantType)
	form.Set("client_id", clientID)
	form.Set("device_code", deviceCode)
	return postTokenForm(ctx, tokenEndpoint, form)
}

func refreshOIDCToken(ctx context.Context, session *CLIOIDCSession) (*CLIOIDCSession, error) {
	if session == nil || session.RefreshToken == "" {
		return session, fmt.Errorf("no refresh token is available; run go-go-host login again")
	}
	discovery, err := discoverOIDC(ctx, session.Issuer)
	if err != nil {
		return session, err
	}
	form := url.Values{}
	form.Set("grant_type", "refresh_token")
	form.Set("client_id", session.ClientID)
	form.Set("refresh_token", session.RefreshToken)
	tok, oauthErr, err := postTokenForm(ctx, discovery.TokenEndpoint, form)
	if err != nil {
		return session, err
	}
	if oauthErr != nil {
		return session, fmt.Errorf("refresh token failed: %s", oauthErr.ErrorString())
	}
	return sessionWithToken(session.Issuer, session.ClientID, session.Scopes, tok, session.RefreshToken), nil
}

func revokeOIDCToken(ctx context.Context, session *CLIOIDCSession) (bool, error) {
	if session == nil || session.RefreshToken == "" {
		return false, nil
	}
	discovery, err := discoverOIDC(ctx, session.Issuer)
	if err != nil {
		return false, err
	}
	if discovery.RevocationEndpoint == "" {
		return false, nil
	}
	form := url.Values{}
	form.Set("client_id", session.ClientID)
	form.Set("token", session.RefreshToken)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, discovery.RevocationEndpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return false, httpResponseError("POST", discovery.RevocationEndpoint, resp)
	}
	return true, nil
}

func postTokenForm(ctx context.Context, endpoint string, form url.Values) (tokenResponse, *oauthErrorResponse, error) {
	var tok tokenResponse
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return tok, nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return tok, nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		if err := json.NewDecoder(resp.Body).Decode(&tok); err != nil {
			return tok, nil, err
		}
		if tok.AccessToken == "" {
			return tok, nil, fmt.Errorf("token response did not include access_token")
		}
		return tok, nil, nil
	}
	var oauthErr oauthErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&oauthErr); err != nil || oauthErr.Error == "" {
		return tok, nil, httpResponseError("POST", endpoint, resp)
	}
	return tok, &oauthErr, nil
}

func postForm(ctx context.Context, endpoint string, form url.Values, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var oauthErr oauthErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&oauthErr); err == nil && oauthErr.Error != "" {
			return fmt.Errorf("oauth error: %s", oauthErr.ErrorString())
		}
		return httpResponseError("POST", endpoint, resp)
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

func httpResponseError(method, endpoint string, resp *http.Response) error {
	data, _ := io.ReadAll(resp.Body)
	message := strings.TrimSpace(string(data))
	if message == "" {
		message = resp.Status
	}
	return fmt.Errorf("%s %s: unexpected status %s: %s", method, endpoint, resp.Status, message)
}

func scopesFromString(raw string, fallback []string) []string {
	if strings.TrimSpace(raw) == "" {
		return fallback
	}
	parts := strings.FieldsFunc(raw, func(r rune) bool { return r == ',' || r == ' ' || r == '\n' || r == '\t' })
	out := make([]string, 0, len(parts))
	seen := map[string]struct{}{}
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if _, ok := seen[part]; ok {
			continue
		}
		seen[part] = struct{}{}
		out = append(out, part)
	}
	if len(out) == 0 {
		return fallback
	}
	return out
}

func defaultScopes(scopes []string) []string {
	if len(scopes) > 0 {
		return scopes
	}
	return []string{"openid", "profile", "email"}
}

func sessionWithToken(issuer, clientID string, scopes []string, tok tokenResponse, previousRefreshToken string) *CLIOIDCSession {
	refreshToken := tok.RefreshToken
	if refreshToken == "" {
		refreshToken = previousRefreshToken
	}
	expiresAt := time.Time{}
	if tok.ExpiresIn > 0 {
		expiresAt = time.Now().UTC().Add(time.Duration(tok.ExpiresIn) * time.Second)
	}
	tokenType := tok.TokenType
	if tokenType == "" {
		tokenType = "Bearer"
	}
	return &CLIOIDCSession{
		Issuer:       issuer,
		ClientID:     clientID,
		Scopes:       scopes,
		AccessToken:  tok.AccessToken,
		IDToken:      tok.IDToken,
		RefreshToken: refreshToken,
		TokenType:    tokenType,
		ExpiresAt:    expiresAt,
	}
}

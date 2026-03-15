package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	nethttp "net/http"
	"net/url"
	"strings"
)

const (
	defaultLINEAuthorizeURL = "https://access.line.me/oauth2/v2.1/authorize"
	defaultLINETokenURL     = "https://api.line.me/oauth2/v2.1/token"
	defaultLINEVerifyURL    = "https://api.line.me/oauth2/v2.1/verify"
)

var (
	ErrInvalidOAuthState     = errors.New("auth: invalid oauth state")
	ErrMissingOAuthCode      = errors.New("auth: missing oauth code")
	ErrLineTokenExchange     = errors.New("auth: line token exchange failed")
	ErrLineTokenVerification = errors.New("auth: line token verification failed")
)

type LINEConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

type LineTokenResponse struct {
	IDToken string `json:"id_token"`
}

type LineIdentity struct {
	Subject     string `json:"sub"`
	DisplayName string `json:"name"`
	Nonce       string `json:"nonce"`
}

type LineProvider interface {
	BuildAuthorizationURL(state string, nonce string) (string, error)
	ExchangeCode(ctx context.Context, code string) (LineTokenResponse, error)
	VerifyIDToken(ctx context.Context, idToken string, nonce string) (LineIdentity, error)
}

type HTTPLineProvider struct {
	client       *nethttp.Client
	clientID     string
	clientSecret string
	redirectURI  string
	authorizeURL string
	tokenURL     string
	verifyURL    string
}

func NewHTTPLineProvider(client *nethttp.Client, cfg LINEConfig) *HTTPLineProvider {
	if client == nil {
		client = nethttp.DefaultClient
	}

	return &HTTPLineProvider{
		client:       client,
		clientID:     cfg.ClientID,
		clientSecret: cfg.ClientSecret,
		redirectURI:  cfg.RedirectURI,
		authorizeURL: defaultLINEAuthorizeURL,
		tokenURL:     defaultLINETokenURL,
		verifyURL:    defaultLINEVerifyURL,
	}
}

func (p *HTTPLineProvider) BuildAuthorizationURL(state string, nonce string) (string, error) {
	redirectURL, err := url.Parse(p.authorizeURL)
	if err != nil {
		return "", fmt.Errorf("parse authorize url: %w", err)
	}

	query := redirectURL.Query()
	query.Set("response_type", "code")
	query.Set("client_id", p.clientID)
	query.Set("redirect_uri", p.redirectURI)
	query.Set("scope", "openid profile")
	query.Set("state", state)
	query.Set("nonce", nonce)
	redirectURL.RawQuery = query.Encode()

	return redirectURL.String(), nil
}

func (p *HTTPLineProvider) ExchangeCode(ctx context.Context, code string) (LineTokenResponse, error) {
	form := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {p.redirectURI},
		"client_id":     {p.clientID},
		"client_secret": {p.clientSecret},
	}

	request, err := nethttp.NewRequestWithContext(ctx, nethttp.MethodPost, p.tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return LineTokenResponse{}, fmt.Errorf("%w: build request: %v", ErrLineTokenExchange, err)
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := p.client.Do(request)
	if err != nil {
		return LineTokenResponse{}, fmt.Errorf("%w: %v", ErrLineTokenExchange, err)
	}
	defer response.Body.Close()

	if response.StatusCode != nethttp.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(response.Body, 1024))
		return LineTokenResponse{}, fmt.Errorf("%w: status %d: %s", ErrLineTokenExchange, response.StatusCode, strings.TrimSpace(string(body)))
	}

	var tokenResponse LineTokenResponse
	if err := json.NewDecoder(response.Body).Decode(&tokenResponse); err != nil {
		return LineTokenResponse{}, fmt.Errorf("%w: decode response: %v", ErrLineTokenExchange, err)
	}

	if tokenResponse.IDToken == "" {
		return LineTokenResponse{}, fmt.Errorf("%w: missing id_token", ErrLineTokenExchange)
	}

	return tokenResponse, nil
}

func (p *HTTPLineProvider) VerifyIDToken(ctx context.Context, idToken string, nonce string) (LineIdentity, error) {
	form := url.Values{
		"id_token":  {idToken},
		"client_id": {p.clientID},
	}

	request, err := nethttp.NewRequestWithContext(ctx, nethttp.MethodPost, p.verifyURL, strings.NewReader(form.Encode()))
	if err != nil {
		return LineIdentity{}, fmt.Errorf("%w: build request: %v", ErrLineTokenVerification, err)
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := p.client.Do(request)
	if err != nil {
		return LineIdentity{}, fmt.Errorf("%w: %v", ErrLineTokenVerification, err)
	}
	defer response.Body.Close()

	if response.StatusCode != nethttp.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(response.Body, 1024))
		return LineIdentity{}, fmt.Errorf("%w: status %d: %s", ErrLineTokenVerification, response.StatusCode, strings.TrimSpace(string(body)))
	}

	var identity LineIdentity
	if err := json.NewDecoder(response.Body).Decode(&identity); err != nil {
		return LineIdentity{}, fmt.Errorf("%w: decode response: %v", ErrLineTokenVerification, err)
	}

	if identity.Subject == "" || identity.DisplayName == "" || identity.Nonce != nonce {
		return LineIdentity{}, fmt.Errorf("%w: invalid verified identity", ErrLineTokenVerification)
	}

	return identity, nil
}

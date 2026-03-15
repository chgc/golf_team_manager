package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

const OAuthStateCookieName = "gtm_line_oauth"

type OAuthState struct {
	State string `json:"state"`
	Nonce string `json:"nonce"`
}

type UpsertLineUserInput struct {
	DisplayName string
	Subject     string
}

type UserStore interface {
	UpsertLineUser(ctx context.Context, input UpsertLineUserInput) (User, error)
}

type LineAuthService struct {
	provider   LineProvider
	users      UserStore
	tokens     TokenManager
	jwtTTL     time.Duration
	randReader io.Reader
	now        func() time.Time
}

func NewLineAuthService(provider LineProvider, users UserStore, tokens TokenManager, jwtTTL time.Duration) *LineAuthService {
	return &LineAuthService{
		provider:   provider,
		users:      users,
		tokens:     tokens,
		jwtTTL:     jwtTTL,
		randReader: rand.Reader,
		now:        time.Now,
	}
}

func (s *LineAuthService) NewLoginFlow() (OAuthState, string, error) {
	flow := OAuthState{}

	state, err := randomToken(s.randReader, 24)
	if err != nil {
		return OAuthState{}, "", fmt.Errorf("generate oauth state: %w", err)
	}
	flow.State = state

	nonce, err := randomToken(s.randReader, 24)
	if err != nil {
		return OAuthState{}, "", fmt.Errorf("generate oauth nonce: %w", err)
	}
	flow.Nonce = nonce

	redirectURL, err := s.provider.BuildAuthorizationURL(flow.State, flow.Nonce)
	if err != nil {
		return OAuthState{}, "", fmt.Errorf("build authorize url: %w", err)
	}

	return flow, redirectURL, nil
}

func (s *LineAuthService) CompleteCallback(ctx context.Context, code string, flow OAuthState, returnedState string) (string, User, error) {
	if code == "" {
		return "", User{}, ErrMissingOAuthCode
	}

	if flow.State == "" || flow.Nonce == "" || returnedState == "" || flow.State != returnedState {
		return "", User{}, ErrInvalidOAuthState
	}

	tokenResponse, err := s.provider.ExchangeCode(ctx, code)
	if err != nil {
		return "", User{}, err
	}

	identity, err := s.provider.VerifyIDToken(ctx, tokenResponse.IDToken, flow.Nonce)
	if err != nil {
		return "", User{}, err
	}

	user, err := s.users.UpsertLineUser(ctx, UpsertLineUserInput{
		DisplayName: identity.DisplayName,
		Subject:     identity.Subject,
	})
	if err != nil {
		return "", User{}, fmt.Errorf("upsert line user: %w", err)
	}

	token, err := s.tokens.Sign(user.Claims(s.now().UTC(), s.jwtTTL))
	if err != nil {
		return "", User{}, fmt.Errorf("sign app jwt: %w", err)
	}

	return token, user, nil
}

func EncodeOAuthState(flow OAuthState) (string, error) {
	marshaled, err := json.Marshal(flow)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(marshaled), nil
}

func DecodeOAuthState(value string) (OAuthState, error) {
	decoded, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil {
		return OAuthState{}, err
	}

	var flow OAuthState
	if err := json.Unmarshal(decoded, &flow); err != nil {
		return OAuthState{}, err
	}

	return flow, nil
}

func randomToken(reader io.Reader, length int) (string, error) {
	buffer := make([]byte, length)
	if _, err := io.ReadFull(reader, buffer); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(buffer), nil
}

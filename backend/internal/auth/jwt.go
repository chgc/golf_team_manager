package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrInvalidToken = errors.New("auth: invalid token")
	ErrExpiredToken = errors.New("auth: expired token")
)

type Claims struct {
	Subject         string   `json:"sub"`
	PlayerID        string   `json:"player_id,omitempty"`
	Provider        Provider `json:"provider"`
	ProviderSubject string   `json:"provider_subject"`
	Role            Role     `json:"role"`
	DisplayName     string   `json:"display_name"`
	ExpiresAt       int64    `json:"exp"`
	IssuedAt        int64    `json:"iat"`
}

func (c Claims) Principal() Principal {
	return Principal{
		DisplayName: c.DisplayName,
		PlayerID:    c.PlayerID,
		Provider:    c.Provider,
		Role:        c.Role,
		Subject:     c.ProviderSubject,
		UserID:      c.Subject,
	}
}

type TokenManager interface {
	Sign(claims Claims) (string, error)
	Validate(token string) (Claims, error)
}

type HMACTokenManager struct {
	secret []byte
	now    func() time.Time
}

func NewHMACTokenManager(secret string) *HMACTokenManager {
	return &HMACTokenManager{
		secret: []byte(secret),
		now:    time.Now,
	}
}

func (m *HMACTokenManager) Sign(claims Claims) (string, error) {
	if err := validateClaims(claims); err != nil {
		return "", err
	}

	headerSegment, err := encodeJWTSegment(map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	})
	if err != nil {
		return "", fmt.Errorf("encode header: %w", err)
	}

	claimsSegment, err := encodeJWTSegment(claims)
	if err != nil {
		return "", fmt.Errorf("encode claims: %w", err)
	}

	signingInput := headerSegment + "." + claimsSegment
	signature := signJWT(signingInput, m.secret)

	return signingInput + "." + signature, nil
}

func (m *HMACTokenManager) Validate(token string) (Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return Claims{}, ErrInvalidToken
	}

	var header struct {
		Algorithm string `json:"alg"`
		Type      string `json:"typ"`
	}
	if err := decodeJWTSegment(parts[0], &header); err != nil {
		return Claims{}, fmt.Errorf("%w: decode header", ErrInvalidToken)
	}

	if header.Algorithm != "HS256" || header.Type != "JWT" {
		return Claims{}, ErrInvalidToken
	}

	expectedSignature := signJWT(parts[0]+"."+parts[1], m.secret)
	if !hmac.Equal([]byte(parts[2]), []byte(expectedSignature)) {
		return Claims{}, ErrInvalidToken
	}

	var claims Claims
	if err := decodeJWTSegment(parts[1], &claims); err != nil {
		return Claims{}, fmt.Errorf("%w: decode claims", ErrInvalidToken)
	}

	if err := validateClaims(claims); err != nil {
		return Claims{}, err
	}

	if m.now().Unix() >= claims.ExpiresAt {
		return Claims{}, ErrExpiredToken
	}

	return claims, nil
}

func validateClaims(claims Claims) error {
	switch {
	case claims.Subject == "":
		return ErrInvalidToken
	case claims.Provider != ProviderDevelopmentStub && claims.Provider != ProviderLINEOAuth:
		return ErrInvalidToken
	case claims.ProviderSubject == "":
		return ErrInvalidToken
	case claims.Role != RoleManager && claims.Role != RolePlayer:
		return ErrInvalidToken
	case claims.DisplayName == "":
		return ErrInvalidToken
	case claims.IssuedAt <= 0 || claims.ExpiresAt <= claims.IssuedAt:
		return ErrInvalidToken
	default:
		return nil
	}
}

func encodeJWTSegment(value any) (string, error) {
	marshaled, err := json.Marshal(value)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(marshaled), nil
}

func decodeJWTSegment(segment string, target any) error {
	decoded, err := base64.RawURLEncoding.DecodeString(segment)
	if err != nil {
		return err
	}

	return json.Unmarshal(decoded, target)
}

func signJWT(signingInput string, secret []byte) string {
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(signingInput))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

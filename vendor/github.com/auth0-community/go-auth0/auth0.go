package auth0

import (
	"errors"
	"net/http"
	"time"

	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

// SecretProvider will provide everything
// needed retrieve the secret.
type SecretProvider interface {
	GetSecret(token *jwt.JSONWebToken) (interface{}, error)
}

// SecretProviderFunc simple wrappers to provide
// secret with functions.
type SecretProviderFunc func(token *jwt.JSONWebToken) (interface{}, error)

// GetSecret implements the SecretProvider interface.
func (f SecretProviderFunc) GetSecret(token *jwt.JSONWebToken) (interface{}, error) {
	return f(token)
}

// NewKeyProvider provide a simple passphrase key provider.
func NewKeyProvider(key interface{}) SecretProvider {
	return SecretProviderFunc(func(_ *jwt.JSONWebToken) (interface{}, error) {
		return key, nil
	})
}

var (
	// ErrNoJWTHeaders is returned when there are no headers in the JWT.
	ErrNoJWTHeaders = errors.New("No headers in the token")
)

// Configuration contains
// all the information about the
// Auth0 service.
type Configuration struct {
	secretProvider SecretProvider
	expectedClaims jwt.Expected
	signIn         jose.SignatureAlgorithm
}

// NewConfiguration creates a configuration for server
func NewConfiguration(provider SecretProvider, audience []string, issuer string, method jose.SignatureAlgorithm) Configuration {
	return Configuration{
		secretProvider: provider,
		expectedClaims: jwt.Expected{Issuer: issuer, Audience: audience},
		signIn:         method,
	}
}

// NewConfigurationTrustProvider creates a configuration for server with no enforcement for token sig alg type, instead trust provider
func NewConfigurationTrustProvider(provider SecretProvider, audience []string, issuer string) Configuration {
	return Configuration{
		secretProvider: provider,
		expectedClaims: jwt.Expected{Issuer: issuer, Audience: audience},
	}
}

// JWTValidator helps middleware
// to validate token
type JWTValidator struct {
	config    Configuration
	extractor RequestTokenExtractor
}

// NewValidator creates a new
// validator with the provided configuration.
func NewValidator(config Configuration, extractor RequestTokenExtractor) *JWTValidator {
	if extractor == nil {
		extractor = RequestTokenExtractorFunc(FromHeader)
	}
	return &JWTValidator{config, extractor}
}

// ValidateRequest validates the token within
// the http request.
// A default leeway value of one minute is used to compare time values.
func (v *JWTValidator) ValidateRequest(r *http.Request) (*jwt.JSONWebToken, error) {
	return v.validateRequestWithLeeway(r, jwt.DefaultLeeway)
}

// ValidateRequestWithLeeway validates the token within
// the http request.
// The provided leeway value is used to compare time values.
func (v *JWTValidator) ValidateRequestWithLeeway(r *http.Request, leeway time.Duration) (*jwt.JSONWebToken, error) {
	return v.validateRequestWithLeeway(r, leeway)
}

func (v *JWTValidator) validateRequestWithLeeway(r *http.Request, leeway time.Duration) (*jwt.JSONWebToken, error) {
	token, err := v.extractor.Extract(r)
	if err != nil {
		return nil, err
	}

	if err := v.validateTokenWithLeeway(token, leeway); err != nil {
		return nil, err
	}

	return token, nil
}

func (v *JWTValidator) ValidateToken(token *jwt.JSONWebToken) error {
	return v.validateTokenWithLeeway(token, jwt.DefaultLeeway)
}

func (v *JWTValidator) ValidateTokenWithLeeway(token *jwt.JSONWebToken, leeway time.Duration) error {
	return v.validateTokenWithLeeway(token, leeway)
}

func (v *JWTValidator) validateTokenWithLeeway(token *jwt.JSONWebToken, leeway time.Duration) error {
	if len(token.Headers) < 1 {
		return ErrNoJWTHeaders
	}

	// trust secret provider when sig alg not configured and skip check
	if v.config.signIn != "" {
		header := token.Headers[0]
		if header.Algorithm != string(v.config.signIn) {
			return ErrInvalidAlgorithm
		}
	}

	claims := jwt.Claims{}
	key, err := v.config.secretProvider.GetSecret(token)
	if err != nil {
		return err
	}

	if err = token.Claims(key, &claims); err != nil {
		return err
	}

	expected := v.config.expectedClaims.WithTime(time.Now())
	err = claims.ValidateWithLeeway(expected, leeway)
	return err
}

// Claims unmarshall the claims of the provided token
func (v *JWTValidator) Claims(token *jwt.JSONWebToken, values ...interface{}) error {
	key, err := v.config.secretProvider.GetSecret(token)
	if err != nil {
		return err
	}
	return token.Claims(key, values...)
}
package utils

import (
	"net/url"

	"github.com/auth0-community/go-auth0"
	"github.com/pkg/errors"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

func ParseAndValidateSignedJWT(domain *url.URL, audience, rawJWT string) (*jwt.JSONWebToken, map[string]interface{}, error) {
	if domain == nil || domain.String() == "" {
		return nil, nil, errors.New("domain is empty")
	}
	if audience == "" {
		return nil, nil, errors.New("audience is empty")
	}

	jwksURL := *domain
	jwksURL.Path = "/.well-known/jwks.json"
	client := auth0.NewJWKClient(auth0.JWKClientOptions{URI: jwksURL.String()}, nil)

	configuration := auth0.NewConfiguration(client, []string{audience}, domain.String(), jose.RS256)
	validator := auth0.NewValidator(configuration, nil)

	token, err := jwt.ParseSigned(rawJWT)
	if err != nil {
		return nil, nil, errors.New("couldn't parse token")
	}
	err = validator.ValidateToken(token)
	if err != nil {
		return nil, nil, errors.New("couldn't validate token")
	}

	claims := make(map[string]interface{})
	err = validator.Claims(token, &claims)
	if err != nil {
		return nil, nil, errors.New("couldn't unmarshal claims")
	}

	return token, claims, nil
}

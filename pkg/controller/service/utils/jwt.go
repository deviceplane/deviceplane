package utils

import (
	"github.com/auth0-community/go-auth0"
	"github.com/pkg/errors"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

func ParseAndValidateSignedJWT(domain, audience, rawJWT string) (*jwt.JSONWebToken, map[string]interface{}, error) {
	client := auth0.NewJWKClient(auth0.JWKClientOptions{URI: domain + ".well-known/jwks.json"}, nil)
	configuration := auth0.NewConfiguration(client, []string{audience}, domain, jose.RS256)
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

| **Build**                                                                                                                           | **Coverage**                                                                                                                                                                 | **Godoc**                                                                                                                           | **Go report**                                                                                                                     | **License**                                                                            |
| ----------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------- |
| [![Build Status](https://travis-ci.org/auth0-community/go-auth0.svg?branch=master)](https://travis-ci.org/auth0-community/go-auth0) | [![Coverage Status](https://coveralls.io/repos/github/auth0-community/go-auth0/badge.svg?branch=master)](https://coveralls.io/github/auth0-community/go-auth0?branch=master) | [![GoDoc](https://godoc.org/github.com/auth0-community/go-auth0?status.png)](https://godoc.org/github.com/auth0-community/go-auth0) | [![Report Cart](http://goreportcard.com/badge/auth0-community/go-auth0)](http://goreportcard.com/report/auth0-community/go-auth0) | [![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE) |

### Contributors

Thanks goes to these wonderful people who contribute(d) or maintain(ed) this repo ([emoji key](https://allcontributors.org/docs/en/emoji-key)):

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore -->
<table>
  <tr>
    <td align="center"><a href="https://twitter.com/beardaway"><img src="https://avatars3.githubusercontent.com/u/11062800?v=4" width="100px;" alt="Conrad Sopala"/><br /><sub><b>Conrad Sopala</b></sub></a><br /><a href="#maintenance-beardaway" title="Maintenance">🚧</a> <a href="#review-beardaway" title="Reviewed Pull Requests">👀</a></td>
    <td align="center"><a href="http://blog.yageek.net"><img src="https://avatars2.githubusercontent.com/u/170917?v=4" width="100px;" alt="Yannick Heinrich"/><br /><sub><b>Yannick Heinrich</b></sub></a><br /><a href="#maintenance-yageek" title="Maintenance">🚧</a> <a href="https://github.com/auth0-community/auth0-go/commits?author=yageek" title="Code">💻</a></td>
    <td align="center"><a href="http://tuanhao.github.io"><img src="https://avatars2.githubusercontent.com/u/18233972?v=4" width="100px;" alt="Hao Chau"/><br /><sub><b>Hao Chau</b></sub></a><br /><a href="https://github.com/auth0-community/auth0-go/commits?author=Tuanhao" title="Code">💻</a></td>
    <td align="center"><a href="http://www.looplab.se"><img src="https://avatars0.githubusercontent.com/u/821518?v=4" width="100px;" alt="Max Ekman"/><br /><sub><b>Max Ekman</b></sub></a><br /><a href="https://github.com/auth0-community/auth0-go/commits?author=maxekman" title="Code">💻</a></td>
    <td align="center"><a href="https://www.linkedin.com/in/tong-tony-jia-a248227b/"><img src="https://avatars1.githubusercontent.com/u/9013477?v=4" width="100px;" alt="Tony Jia"/><br /><sub><b>Tony Jia</b></sub></a><br /><a href="https://github.com/auth0-community/auth0-go/commits?author=ttjiaa" title="Code">💻</a></td>
  </tr>
</table>

<!-- ALL-CONTRIBUTORS-LIST:END -->

# Go-Auth0

<img src="https://img.shields.io/badge/community-driven-brightgreen.svg"/> <br>

go-auth0 is a package helping to authenticate using the [Auth0](https://auth0.com) service.

This repo is supported and maintained by Community Developers, not Auth0. For more information about different support levels check https://auth0.com/docs/support/matrix .

## Getting started

### Installation

```
go get github.com/auth0-community/go-auth0
```

## Usage

### Example

### Gin

Using [Gin](https://github.com/gin-gonic/gin) and the [Auth0 Authorization Extension](https://auth0.com/docs/extensions/authorization-extension), you
may want to implement the authentication auth like the following:

```go
var auth.AdminGroup string = "my_admin_group"

// Access Control Helper function.
func shouldAccess(wantedGroups []string, groups []interface{}) bool {
 /* Fill depending on your needs */
}

// Wrapping a Gin endpoint with Auth0 Groups.
func Auth0Groups(wantedGroups ...string) gin.HandlerFunc {

	return gin.HandlerFunc(func(c *gin.Context) {

		tok, err := validator.ValidateRequest(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			log.Println("Invalid token:", err)
			return
		}

		claims := map[string]interface{}{}
		err = validator.Claims(tok, &claims)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
			c.Abort()
			log.Println("Invalid claims:", err)
			return
		}

		metadata, okMetadata := claims["app_metadata"].(map[string]interface{})
		authorization, okAuthorization := metadata["authorization"].(map[string]interface{})
		groups, hasGroups := authorization["groups"].([]interface{})
		if !okMetadata || !okAuthorization || !hasGroups || !shouldAccess(wantedGroups, groups) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "need more privileges"})
			c.Abort()
			log.Println("Need more provileges")
			return
		}
		c.Next()
	})
}

// Use it
r.PUT("/news", auth.Auth0Groups(auth.AdminGroup), api.GetNews)
```

For a sample usage, take a look inside the `example` directory.

### Usage

#### Client Credentials - HS256

Using HS256, the validation key is the secret you retrieve in the dashboard.

```go
// Creates a configuration with the Auth0 information
secret, _ := base64.URLEncoding.DecodeString(os.Getenv("AUTH0_CLIENT_SECRET"))
secretProvider := auth0.NewKeyProvider(secret)
audience := os.Getenv("AUTH0_CLIENT_ID")

configuration := auth0.NewConfiguration(secretProvider, []string{audience}, "https://mydomain.eu.auth0.com/", jose.HS256)
validator := auth0.NewValidator(configuration, nil)

token, err := validator.ValidateRequest(r)

if err != nil {
    fmt.Println("Token is not valid:", token)
}
```

#### Client Credentials - RS256

Using RS256, the validation key is the certificate you find in advanced settings

```go
// Extracted from https://github.com/square/go-jose/blob/master/utils.go
// LoadPublicKey loads a public key from PEM/DER-encoded data.
// You can download the Auth0 pem file from `applications -> your_app -> scroll down -> Advanced Settings -> certificates -> download`
func LoadPublicKey(data []byte) (interface{}, error) {
	input := data

	block, _ := pem.Decode(data)
	if block != nil {
		input = block.Bytes
	}

	// Try to load SubjectPublicKeyInfo
	pub, err0 := x509.ParsePKIXPublicKey(input)
	if err0 == nil {
		return pub, nil
	}

	cert, err1 := x509.ParseCertificate(input)
	if err1 == nil {
		return cert.PublicKey, nil
	}

	return nil, fmt.Errorf("square/go-jose: parse error, got '%s' and '%s'", err0, err1)
}
// Create a configuration with the Auth0 information
pem, err := ioutil.ReadFile("path/to/your/cert.pem")
if err != nil {
	panic(err)
}
secret, err := LoadPublicKey(sharedKey)
if err != nil {
	panic(err)
}
secretProvider := auth0.NewKeyProvider(secret)
audience := os.Getenv("AUTH0_CLIENT_ID")

configuration := auth0.NewConfiguration(secretProvider, []string{audience}, "https://mydomain.eu.auth0.com/", jose.RS256)
validator := auth0.NewValidator(configuration, nil)

token, err := validator.ValidateRequest(r)

if err != nil {
    fmt.Println("Token is not valid:", token)
}
```

#### API with JWK

```go
client := NewJWKClient(JWKClientOptions{URI: "https://mydomain.eu.auth0.com/.well-known/jwks.json"}, nil)
audience := os.Getenv("AUTH0_CLIENT_ID")
configuration := NewConfiguration(client, []string{audience}, "https://mydomain.eu.auth0.com/", jose.RS256)
validator := NewValidator(configuration, nil)

token, err := validator.ValidateRequest(r)

if err != nil {
    fmt.Println("Token is not valid:", token)
}
```

#### Support interface for configurable key cacher

```go
opts := JWKClientOptions{URI: "https://mydomain.eu.auth0.com/.well-known/jwks.json"}
// Creating key cacher with max age of 100sec and max size of 5 entries.
// Defaults to persistent key cacher if not specified when creating a client.
keyCacher := NewMemoryKeyCacher(time.Duration(100) * time.Second, 5)
client := NewJWKClientWithCache(opts, nil, keyCacher)

searchedKey, err := client.GetKey("KEY_ID")

if err != nil {
	fmt.Println("Cannot get key because of", err)
}
```

#### Validating a token outside an HTTP request

Sometimes a token is received from something that is not an HTTP request (such as a GRPC call)
when a token must be validated in those cases, use `ValidateToken` directly

```go
client := NewJWKClient(JWKClientOptions{URI: "https://mydomain.eu.auth0.com/.well-known/jwks.json"}, nil)
audience := os.Getenv("AUTH0_CLIENT_ID")
configuration := NewConfiguration(client, []string{audience}, "https://mydomain.eu.auth0.com/", jose.RS256)
validator := NewValidator(configuration, nil)

token = retrieveToken() // Retrieves the token through custom logic, such as for example fetching it from GRPC metadata
if err := validator.ValidateToken(token); err != nil {
	fmt.Println("Cannot validate token because of", err)
}
```

## Contribute

Feel like contributing to this repo? We're glad to hear that! Before you start contributing please visit our [Contributing Guideline](https://github.com/auth0-community/getting-started/blob/master/CONTRIBUTION.md) .

Here you can also find the [PR template](https://github.com/auth0-community/go-auth0/blob/master/PULL_REQUEST_TEMPLATE.md) to fill once creating a PR. It will automatically appear once you open a pull request.

## Issues Reporting

Spotted a bug or any other kind of issue? We're just humans and we're always waiting for constructive feedback! Check our section on how to [report issues](https://github.com/auth0-community/getting-started/blob/master/CONTRIBUTION.md#issues)!

Here you can also find the [Issue template](https://github.com/auth0-community/go-auth0/blob/master/ISSUE_TEMPLATE.md) to fill once opening a new issue. It will automatically appear once you create an issue.

## Repo Community

Feel like PRs and issues are not enough? Want to dive into further discussion about the tool? We created topics for each Auth0 Community repo so that you can join discussion on stack available on our repos. Here it is for this one: [go-auth0](https://community.auth0.com/t/auth0-community-oss-go-auth0/15969)

<a href="https://community.auth0.com/">
<img src="/assets/join_auth0_community_badge.png"/>
</a>

## License

This project is licensed under the MIT license. See the [LICENSE](https://github.com/auth0-community/go-auth0/blob/master/LICENSE) file for more info.

## What is Auth0?

Auth0 helps you to:

- Add authentication with [multiple authentication sources](https://docs.auth0.com/identityproviders), either social like

  - Google
  - Facebook
  - Microsoft
  - Linkedin
  - GitHub
  - Twitter
  - Box
  - Salesforce
  - etc.

  **or** enterprise identity systems like:

  - Windows Azure AD
  - Google Apps
  - Active Directory
  - ADFS
  - Any SAML Identity Provider

- Add authentication through more traditional [username/password databases](https://docs.auth0.com/mysql-connection-tutorial)
- Add support for [linking different user accounts](https://docs.auth0.com/link-accounts) with the same user
- Support for generating signed [JSON Web Tokens](https://docs.auth0.com/jwt) to call your APIs and create user identity flow securely
- Analytics of how, when and where users are logging in
- Pull data from other sources and add it to user profile, through [JavaScript rules](https://docs.auth0.com/rules)

## Create a free Auth0 account

- Go to [Auth0 website](https://auth0.com/signup)
- Hit the **SIGN UP** button in the upper-right corner

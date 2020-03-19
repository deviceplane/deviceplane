package models

type SsoJWT struct {
	Email    string
	Name     string
	Provider string
	Subject  string
	Claims   map[string]interface{}
}

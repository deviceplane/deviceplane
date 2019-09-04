package validator

import (
	"regexp"
	"sync"

	"gopkg.in/go-playground/validator.v9"
)

var (
	vldr          = validator.New()
	standardRegex = regexp.MustCompile("^[a-zA-Z0-9_-]+$")
	once          sync.Once
)

func Validate(s interface{}) error {
	once.Do(func() {
		vldr.RegisterValidation("standard", func(fl validator.FieldLevel) bool {
			return standardRegex.Match([]byte(fl.Field().String()))
		})
		vldr.RegisterAlias("id", "required,min=1,max=32,standard")
		vldr.RegisterAlias("name", "required,min=1,max=100,standard")
		vldr.RegisterAlias("labelkey", "required,min=1,max=100,standard")
		vldr.RegisterAlias("labelvalue", "required,min=1,max=100")
		vldr.RegisterAlias("password", "required,min=8,max=100")
		vldr.RegisterAlias("config", "required,min=1,max=5000")
		vldr.RegisterAlias("description", "max=5000")
	})
	return vldr.Struct(s)
}

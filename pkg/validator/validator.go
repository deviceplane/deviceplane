package validator

import (
	"regexp"
	"sync"

	"gopkg.in/go-playground/validator.v9"
)

var (
	vldr = validator.New()
	once sync.Once

	internalTitleRegex       = regexp.MustCompile(`^[a-zA-Z0-9]+_[a-zA-Z0-9]+$`)
	userTitleRegex           = regexp.MustCompile(`^[a-zA-Z0-9-]+$`)
	environmentVariableRegex = regexp.MustCompile(`^[a-zA-Z]+[a-zA-Z0-9_]*$`)
)

func Validate(s interface{}) error {
	once.Do(func() {
		vldr.RegisterValidation("internaltitle", func(fl validator.FieldLevel) bool {
			return internalTitleRegex.Match([]byte(fl.Field().String()))
		})
		vldr.RegisterValidation("usertitle", func(fl validator.FieldLevel) bool {
			return userTitleRegex.Match([]byte(fl.Field().String()))
		})
		vldr.RegisterValidation("environmentvariable", func(fl validator.FieldLevel) bool {
			return environmentVariableRegex.Match([]byte(fl.Field().String()))
		})

		vldr.RegisterAlias("id", "required,min=1,max=32,internaltitle")
		vldr.RegisterAlias("name", "required,min=1,max=100,usertitle")
		vldr.RegisterAlias("labelkey", "required,min=1,max=100,usertitle")
		vldr.RegisterAlias("labelvalue", "required,min=1,max=100")
		vldr.RegisterAlias("environmentvariablekey", "required,min=1,max=100,environmentvariable")
		vldr.RegisterAlias("environmentvariablevalue", "required,min=1,max=500")
		vldr.RegisterAlias("password", "required,min=8,max=100")
		vldr.RegisterAlias("config", "required,min=1,max=5000")
		vldr.RegisterAlias("description", "max=5000")
		vldr.RegisterAlias("protocol", "eq=tcp|eq=http")
		vldr.RegisterAlias("port", "required,min=1,max=65535")
	})
	return vldr.Struct(s)
}

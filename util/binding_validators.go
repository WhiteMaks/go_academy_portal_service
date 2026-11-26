package util

import (
	"fmt"
	"reflect"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// RegisterBindingValidators registers project-specific gin binding rules.
func RegisterBindingValidators() error {
	engine, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return fmt.Errorf("binding validator engine is not *validator.Validate")
	}
	return engine.RegisterValidation("e164", func(fl validator.FieldLevel) bool {
		if fl.Field().Kind() != reflect.String {
			return false
		}
		return IsPhoneE164(fl.Field().String())
	})
}

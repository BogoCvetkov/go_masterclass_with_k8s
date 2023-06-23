package controller

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if cur, ok := fl.Field().Interface().(string); ok {
		switch cur {
		case "USD", "EUR", "CAD":
			return true
		}
	}

	return false
}

func RegisterValidation() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("validcurrency", validCurrency)
	}

}

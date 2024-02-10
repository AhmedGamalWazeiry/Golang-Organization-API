package utils

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func BindAndValidate(c *gin.Context, model interface{}) (int, string) {
	if err := c.ShouldBindJSON(model); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			for _, e := range errs {
				switch e.Tag() {
				case "required":
					return http.StatusBadRequest, fmt.Sprintf("%s is required", e.Field())
				case "email":
					return http.StatusBadRequest, fmt.Sprintf("%s is not a valid email", e.Field())
				case "min":
					return http.StatusBadRequest, fmt.Sprintf("%s must be at least %s characters", e.Field(), e.Param())
				}
			}
		} else {
			// If the error is not of type validator.ValidationErrors, return a generic error message
			return http.StatusBadRequest, "Invalid JSON format"
		}
	}
	return http.StatusOK, ""
}

package utils

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fiqrioemry/asset_management_system_app/server/errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func BindAndValidateJSON[T any](c *gin.Context, req *T) bool {
	if err := c.ShouldBindJSON(req); err != nil {

		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			validationErr := buildValidationError(validationErrors)
			HandleError(c, validationErr)
			return false
		}

		if jsonErr, ok := err.(*json.UnmarshalTypeError); ok {
			parseErr := errors.NewBadRequest("Invalid data type for field").WithContext("field", jsonErr.Field).WithContext("expected_type", jsonErr.Type.String())
			HandleError(c, parseErr)
			return false
		}

		if syntaxErr, ok := err.(*json.SyntaxError); ok {

			parseErr := errors.NewBadRequest("Invalid JSON syntax").WithContext("offset", syntaxErr.Offset)
			HandleError(c, parseErr)
			return false
		}

		parseErr := errors.NewBadRequest("Invalid JSON format")
		HandleError(c, parseErr)
		return false
	}
	return true
}
func BindAndValidateForm[T any](c *gin.Context, req *T) bool {
	if err := c.ShouldBind(req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			validationErr := buildValidationError(validationErrors)
			HandleError(c, validationErr)
			return false
		}

		formErr := errors.NewBadRequest("Invalid form data format")
		HandleError(c, formErr)
		return false
	}
	return true
}

func buildValidationError(validationErrors validator.ValidationErrors) *errors.AppError {
	errorDetails := make(map[string]string)

	for _, fieldError := range validationErrors {
		fieldName := strings.ToLower(fieldError.Field())

		switch fieldError.Tag() {
		case "required":
			errorDetails[fieldName] = fmt.Sprintf("%s is required", fieldName)
		case "email":
			errorDetails[fieldName] = "Please provide a valid email address"
		case "min":
			errorDetails[fieldName] = fmt.Sprintf("%s must be at least %s characters", fieldName, fieldError.Param())
		case "max":
			errorDetails[fieldName] = fmt.Sprintf("%s must be at most %s characters", fieldName, fieldError.Param())
		case "len":
			errorDetails[fieldName] = fmt.Sprintf("%s must be exactly %s characters", fieldName, fieldError.Param())
		case "numeric":
			errorDetails[fieldName] = fmt.Sprintf("%s must be numeric", fieldName)
		case "alpha":
			errorDetails[fieldName] = fmt.Sprintf("%s must contain only letters", fieldName)
		case "alphanum":
			errorDetails[fieldName] = fmt.Sprintf("%s must contain only letters and numbers", fieldName)
		case "url":
			errorDetails[fieldName] = fmt.Sprintf("%s must be a valid URL", fieldName)
		case "uuid":
			errorDetails[fieldName] = fmt.Sprintf("%s must be a valid UUID", fieldName)
		default:
			errorDetails[fieldName] = fmt.Sprintf("%s is invalid", fieldName)
		}
	}

	err := errors.NewBadRequest("Validation failed")
	err.WithContext("validation_errors", errorDetails)
	err.WithContext("failed_fields", len(errorDetails))

	return err
}

func ValidateStruct(s any) error {
	validate := validator.New()
	return validate.Struct(s)
}

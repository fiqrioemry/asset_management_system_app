package utils

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fiqrioemry/go-api-toolkit/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func BindAndValidateJSON[T any](c *gin.Context, req *T) bool {
	if err := c.ShouldBindJSON(req); err != nil {

		// Handle validation errors
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			validationErr := buildValidationError(validationErrors)
			response.Error(c, validationErr)
			return false
		}
		// Handle specific JSON parsing errors
		if jsonErr, ok := err.(*json.UnmarshalTypeError); ok {
			parseErr := response.NewBadRequest("Invalid data type for field").WithContext("field", jsonErr.Field).WithContext("expected_type", jsonErr.Type.String())
			response.Error(c, parseErr)
			return false
		}

		// Handle syntax errors in JSON
		if syntaxErr, ok := err.(*json.SyntaxError); ok {
			parseErr := response.NewBadRequest("Invalid JSON syntax").WithContext("offset", syntaxErr.Offset)
			response.Error(c, parseErr)
			return false
		}

		// Handle other binding errors
		parseErr := response.NewBadRequest("Invalid JSON format")
		response.Error(c, parseErr)
		return false
	}
	return true
}

func BindAndValidateForm[T any](c *gin.Context, req *T) bool {
	if err := c.ShouldBind(req); err != nil {

		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			// Handle validation errors
			validationErr := buildValidationError(validationErrors)
			response.Error(c, validationErr)
			return false
		}

		formErr := response.NewBadRequest(fmt.Sprintf("Invalid form data format: %v", err))
		response.Error(c, formErr)
		return false
	}
	return true
}
func buildValidationError(validationErrors validator.ValidationErrors) *response.AppError {
	errorDetails := make(map[string]any)

	// Iterate over validation errors and build a user-friendly error message
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

	err := response.NewBadRequest("Validation failed")
	err.WithContext("errors", errorDetails)

	return err
}

func ValidateStruct(s any) error {
	validate := validator.New()
	return validate.Struct(s)
}

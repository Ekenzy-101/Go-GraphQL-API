package service

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	lowercaseRegex        = regexp.MustCompile(`[a-z]+`)
	nameRegex             = regexp.MustCompile(`^[a-zA-z ]+$`)
	numberRegex           = regexp.MustCompile(`\d+`)
	specialCharacterRegex = regexp.MustCompile(`\W+`)
	uppercaseRegex        = regexp.MustCompile(`[A-Z]+`)
)

var validate *validator.Validate

var rules = map[string]string{
	"content":  "required,max=2000",
	"email":    "lowercase,email,max=255",
	"id":       "uuid4",
	"name":     "required,name,max=100",
	"password": "required,min=8,max=128,password",
	"title":    "required,max=255",
}

func init() {
	validate = validator.New()
	err := validate.RegisterValidation("name", validateName)
	if err != nil {
		log.Fatalf("RegisterValidation [name] %v", err)
	}

	err = validate.RegisterValidation("password", validatePassword)
	if err != nil {
		log.Fatalf("RegisterValidation [password] %v", err)
	}
}

func validateName(fl validator.FieldLevel) bool {
	return nameRegex.MatchString(fl.Field().String())
}

func validatePassword(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return specialCharacterRegex.MatchString(value) &&
		lowercaseRegex.MatchString(value) &&
		uppercaseRegex.MatchString(value) &&
		numberRegex.MatchString(value)
}

func (s *service) ValidateArgs(ctx context.Context, args map[string]interface{}) error {
	for field, value := range args {
		err := validate.VarCtx(ctx, value, rules[field])
		fieldErrors, ok := err.(validator.ValidationErrors)
		if err != nil && !ok {
			return err
		}

		if err != nil {
			return transformFieldError(field, fieldErrors[0])
		}
	}
	return nil
}

func transformFieldError(field string, err validator.FieldError) error {
	switch err.ActualTag() {
	case "name":
		return fmt.Errorf("%v should contain only letters and spaces", strings.Title(field))
	case "email":
		return fmt.Errorf("%v is not a valid email address", strings.Title(field))
	case "gt":
		return fmt.Errorf("%v should be greater than %v", strings.Title(field), err.Param())
	case "lte":
		return fmt.Errorf("%v should not be greater than %v", strings.Title(field), err.Param())
	case "max":
		return fmt.Errorf("%v should not be greater than %v characters", strings.Title(field), err.Param())
	case "min":
		return fmt.Errorf("%v should not be less than %v characters", strings.Title(field), err.Param())
	case "oneof":
		return fmt.Errorf("%v should be in these category %v", strings.Title(field), err.Param())
	case "password":
		return fmt.Errorf("%v should be a mix of uppercase, lowercase, numeric and special characters", strings.Title(field))
	case "required":
		return fmt.Errorf("%v is required", strings.Title(field))
	default:
		return fmt.Errorf("%v is in an invalid format", strings.Title(field))
	}
}

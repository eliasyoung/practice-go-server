package main

import (
	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}

func validateCreatePostPayload(data *CreatePostPayload) error {
	rules := map[string]string{
		"Title":   "required,max=100",
		"Content": "required,max=1000",
	}

	Validate.RegisterStructValidationMapRules(rules, CreatePostPayload{})

	err := Validate.Struct(data)

	return err
}

func validateUpdatePostPayload(data *UpdatePostPayload) error {
	rules := map[string]string{
		"Title":   "omitempty,max=100",
		"Content": "omitempty,max=1000",
	}

	Validate.RegisterStructValidationMapRules(rules, CreatePostPayload{})

	err := Validate.Struct(data)

	return err
}

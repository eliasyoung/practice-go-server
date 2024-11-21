package main

import (
	"github.com/eliasyoung/go-backend-server-practice/internal/db"
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

func validateCreateUserPayload(data *db.CreateUserParams) error {
	rules := map[string]string{
		"Username": "required,max=20",
		"Password": "required,max=20",
		"Email":    "required,email,max=320",
	}

	Validate.RegisterStructValidationMapRules(rules, db.CreateUserParams{})

	err := Validate.Struct(data)

	return err
}

func validateFollowUserPayload(data *FollowUser) error {
	rules := map[string]string{
		"UserID": "required",
	}

	Validate.RegisterStructValidationMapRules(rules, FollowUser{})

	err := Validate.Struct(data)

	return err
}

func validatePagiationQuery(data *db.PaginatedFeedQuery) error {
	rules := map[string]string{
		"Limit":  "gte=1,lte=20",
		"Offset": "gte=0",
		"Sort":   "oneof=asc desc",
		"Tags":   "max=5",
		"Search": "max=100",
		"Since":  "max=200",
		"Until":  "max=200",
	}

	Validate.RegisterStructValidationMapRules(rules, db.PaginatedFeedQuery{})

	err := Validate.Struct(data)

	return err
}

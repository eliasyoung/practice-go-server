package utils

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Password struct {
	text *string
	hash []byte
}

func (p *Password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	p.text = &text
	p.hash = hash

	return nil
}

func (p *Password) GetHash() ([]byte, error) {
	if p.hash != nil {
		return p.hash, nil
	}
	return nil, errors.New("hash is empty!")
}

func PasswordHasher(text string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	return hash, nil
}

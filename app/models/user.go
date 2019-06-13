package models

import (
	"github.com/revel/revel"
)

type User struct {
	UserId             int
	Name               string
	Username, Password string
	HashedPassword     []byte
	Admin              bool
}

func ValidatePassword(v *revel.Validation, password string) *revel.ValidationResult {
	return v.Check(password,
		revel.Required{},
		revel.MinSize{6},
		revel.MaxSize{20},
	)
}

package helpers

import (
	"github.com/zeeshanahmad0201/todo-mongo/common"
	"golang.org/x/crypto/bcrypt"
)

// Encrypt the password before it is stored in the DB
func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		common.HandleError(err, common.ErrorHandlerConfig{
			Exit: true,
		})
	}

	return string(bytes)
}

// Validates and verifies the password in the DB
func VerifyPassword(userPassword *string, providedPassword *string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(*providedPassword), []byte(*userPassword))
	if err != nil {
		common.HandleError(err)
		return false, err
	}
	return true, nil
}

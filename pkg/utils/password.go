package utils

import (
    "golang.org/x/crypto/bcrypt"
)

// PasswordUtils provides utility functions for password handling
type PasswordUtils struct {
    cost int
}

func NewPasswordUtils(cost int) *PasswordUtils {
    if cost == 0 {
        cost = bcrypt.DefaultCost
    }
    return &PasswordUtils{cost: cost}
}

// HashPassword creates a bcrypt hash from password
func (p *PasswordUtils) HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), p.cost)
    if err != nil {
        return "", err
    }
    return string(bytes), nil
}

// CheckPassword compares password with hash
func (p *PasswordUtils) CheckPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
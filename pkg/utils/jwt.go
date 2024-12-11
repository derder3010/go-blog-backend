package utils

import (
    "errors"
    "time"
    "github.com/golang-jwt/jwt/v4"
)

type JWTUtils struct {
    secretKey []byte
    expires   time.Duration
}

type JWTClaims struct {
    UserID string `json:"user_id"`
    Email  string `json:"email"`
    jwt.RegisteredClaims
}

func NewJWTUtils(secretKey string, expires time.Duration) *JWTUtils {
    return &JWTUtils{
        secretKey: []byte(secretKey),
        expires:   expires,
    }
}

// GenerateToken creates a new JWT token
func (j *JWTUtils) GenerateToken(userID, email string) (string, error) {
    claims := JWTClaims{
        userID,
        email,
        jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.expires)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(j.secretKey)
}

// ValidateToken validates the JWT token and returns claims
func (j *JWTUtils) ValidateToken(tokenString string) (*JWTClaims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return j.secretKey, nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
        return claims, nil
    }

    return nil, errors.New("invalid token")
}
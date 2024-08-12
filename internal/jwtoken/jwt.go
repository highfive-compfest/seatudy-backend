package jwtoken

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/highfive-compfest/seatudy-backend/internal/config"
	"time"
)

type AccessClaims struct {
	jwt.RegisteredClaims
	Email           string `json:"email"`
	IsEmailVerified bool   `json:"is_email_verified"`
	Name            string `json:"name"`
	Role            string `json:"role"`
}

func CreateAccessJWT(id, email string, isEmailVerified bool, name string, role string) (string, error) {
	claims := AccessClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   id,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.Env.JwtAccessDuration)),
			Issuer:    "seatudy-backend-accesstoken",
		},
		Email:           email,
		IsEmailVerified: isEmailVerified,
		Name:            name,
		Role:            role,
	}

	unsignedJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedJWT, err := unsignedJWT.SignedString(config.Env.JwtAccessSecret)
	if err != nil {
		return "", err
	}

	return signedJWT, nil
}

func CreateRefreshJWT(id string) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   id,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.Env.JwtRefreshDuration)),
		Issuer:    "seatudy-backend-refreshtoken",
	}

	unsignedJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedJWT, err := unsignedJWT.SignedString(config.Env.JwtRefreshSecret)
	if err != nil {
		return "", err
	}

	return signedJWT, nil
}

func DecodeAccessJWT(tokenString string) (*AccessClaims, error) {
	var claims AccessClaims

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return config.Env.JwtAccessSecret, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return &claims, nil
}

func DecodeRefreshJWT(tokenString string) (*jwt.RegisteredClaims, error) {
	var claims jwt.RegisteredClaims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return config.Env.JwtRefreshSecret, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return &claims, nil
}

package service

import (
	"Agora/constants"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type (
	IJWTService interface {
		GenerateToken(userID string, role string) (string, string, error)
		ValidateToken(token string) (*jwt.Token, *jwtCustomClaims, error)
	}

	jwtCustomClaims struct {
		UserID string `json:"user_id"`
		Role   string `json:"role"`
		jwt.RegisteredClaims
	}

	JWTService struct {
		secretKey string
		issuer    string
	}
)

func getSecretKey() string {
	key := os.Getenv("JWT_SECRET")
	if key == "" {
		key = "TemplateSecretKey"
	}
	return key
}

func NewJWTService() *JWTService {
	return &JWTService{
		secretKey: getSecretKey(),
		issuer:    "Agora",
	}
}

func (j *JWTService) GenerateToken(userID, role string) (string, string, error) {
	// Access Token
	accessClaims := jwtCustomClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)), // bisa disesuaikan
			Issuer:    j.issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(j.secretKey))
	if err != nil {
		return "", "", constants.ErrGenerateAccessToken
	}

	// Refresh Token
	refreshClaims := jwtCustomClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			Issuer:    j.issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(j.secretKey))
	if err != nil {
		return "", "", constants.ErrGenerateRefreshToken
	}

	return accessToken, refreshToken, nil
}

func (j *JWTService) ValidateToken(tokenString string) (*jwt.Token, *jwtCustomClaims, error) {
	claims := &jwtCustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, constants.ErrUnexpectedSigningMethod
		}
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return nil, nil, constants.ErrValidateToken
	}
	if !token.Valid {
		return nil, nil, constants.ErrTokenInvalid
	}

	return token, claims, nil
}

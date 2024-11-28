package utils

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUUID   string
	RefreshUUID  string
	AtExpires    int64
	RtExpires    int64
}

type CustomClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateTokenPair(userID, email, role string) (*TokenDetails, error) {
	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUUID = uuid.New().String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUUID = uuid.New().String()

	var err error
	// Access Token
	atClaims := CustomClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(td.AtExpires, 0)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        td.AccessUUID,
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = accessToken.SignedString([]byte(getJWTSecret()))
	if err != nil {
		return nil, err
	}

	// Refresh Token
	rtClaims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Unix(td.RtExpires, 0)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ID:        td.RefreshUUID,
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = refreshToken.SignedString([]byte(getJWTSecret()))
	if err != nil {
		return nil, err
	}

	return td, nil
}

func ValidateToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			// Verify the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(getJWTSecret()), nil
		},
	)

	if err != nil {
		return nil, err
	}

	// Extract claims
	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// Additional validation
	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("token is expired")
	}

	return claims, nil
}

func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// Fallback to a default secret (not recommended for production)
		secret = "your-secret-key"
	}
	return secret
}

func RefreshAccessToken(refreshToken string) (*TokenDetails, error) {
	// Validate the refresh token
	claims, err := ValidateToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// In a real-world scenario, you'd also check if the refresh token is in your
	// blacklist or database of valid refresh tokens

	// Generate a new token pair
	return GenerateTokenPair(
		claims.UserID,
		claims.Email,
		claims.Role,
	)
}

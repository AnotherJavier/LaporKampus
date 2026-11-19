package helper

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/vyerru/lapor-kampus/app/model"
)

var (
	ErrInvalidToken = errors.New("token tidak valid")
	ErrExpiredToken = errors.New("token sudah kedaluwarsa")
)

const (
	tokenTypeAccess  = "access"
	tokenTypeRefresh = "refresh"
)

type accessClaims struct {
	Username  string `json:"username"`
	Role      string `json:"role"`
	TokenType string `json:"typ"`
	jwt.RegisteredClaims
}

type refreshClaims struct {
	TokenID   string `json:"token_id"`
	TokenType string `json:"typ"`
	jwt.RegisteredClaims
}

type JWTManager struct {
	secret    []byte
	issuer    string
	accessTTL time.Duration
}

func NewJWTManager(secret, issuer string, accessTTL time.Duration) *JWTManager {
	return &JWTManager{secret: []byte(secret), issuer: issuer, accessTTL: accessTTL}
}

func (m *JWTManager) AccessTTL() time.Duration { return m.accessTTL }

func (m *JWTManager) GenerateAccess(userID int, username, role string) (string, error) {
	now := time.Now()
	claims := accessClaims{
		Username:  username,
		Role:      role,
		TokenType: tokenTypeAccess,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.Itoa(userID),
			Issuer:    m.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.accessTTL)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

// GenerateRefresh membuat token refresh. tokenID disimpan di database dalam
// bentuk hash sehingga token dapat dicabut satu per satu saat logout.
func (m *JWTManager) GenerateRefresh(userID int, tokenID string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := refreshClaims{
		TokenID:   tokenID,
		TokenType: tokenTypeRefresh,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.Itoa(userID),
			Issuer:    m.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

func (m *JWTManager) Parse(tokenString string) (model.AuthUser, error) {
	claims := &accessClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("algoritma tidak diharapkan: %v", t.Header["alg"])
			}
			return m.secret, nil
		},
		jwt.WithIssuer(m.issuer),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return model.AuthUser{}, ErrExpiredToken
		}
		return model.AuthUser{}, ErrInvalidToken
	}
	if !token.Valid || claims.TokenType != tokenTypeAccess {
		return model.AuthUser{}, ErrInvalidToken
	}
	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return model.AuthUser{}, ErrInvalidToken
	}
	return model.AuthUser{
		UserID:   userID,
		Username: claims.Username,
		Role:     claims.Role,
	}, nil
}

func (m *JWTManager) ParseRefresh(tokenString string) (int, string, error) {
	claims := &refreshClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("algoritma tidak diharapkan: %v", t.Header["alg"])
			}
			return m.secret, nil
		},
		jwt.WithIssuer(m.issuer),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return 0, "", ErrExpiredToken
		}
		return 0, "", ErrInvalidToken
	}
	if !token.Valid || claims.TokenType != tokenTypeRefresh {
		return 0, "", ErrInvalidToken
	}
	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return 0, "", ErrInvalidToken
	}
	return userID, claims.TokenID, nil
}

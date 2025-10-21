// internal/auth/auth.go
package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Tokens struct {
	Access   string
	Refresh  string
	JTIAcc   string
	JTIRef   string
	ExpAcc   time.Time
	ExpRef   time.Time
	UserID   string
	Issuer   string
	Audience string
}

func IssueTokens(userID string) (*Tokens, error) {
	now := time.Now().UTC()
	t := &Tokens{
		UserID:   userID,
		JTIAcc:   uuid.NewString(),
		JTIRef:   uuid.NewString(),
		ExpAcc:   now.Add(15 * time.Minute),
		ExpRef:   now.Add(7 * 24 * time.Hour),
		Issuer:   "intelliquiz-app",
		Audience: "intelliquiz-client",
	}

	acc := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   userID,
		ID:        t.JTIAcc,
		Issuer:    t.Issuer,
		Audience:  jwt.ClaimStrings{t.Audience},
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(t.ExpAcc),
	})

	ref := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   userID,
		ID:        t.JTIRef,
		Issuer:    t.Issuer,
		Audience:  jwt.ClaimStrings{t.Audience},
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(t.ExpRef),
	})

	var err error
	t.Access, err = acc.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, err
	}
	t.Refresh, err = ref.SignedString([]byte(os.Getenv("JWT_REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}
	return t, nil
}

func ParseAccess(tokenStr string) (*jwt.RegisteredClaims, error) {
	secret := os.Getenv("JWT_SECRET")
	return parseWithSecret(tokenStr, secret)
}

func ParseRefresh(tokenStr string) (*jwt.RegisteredClaims, error) {
	secret := os.Getenv("JWT_REFRESH_SECRET")
	return parseWithSecret(tokenStr, secret)
}

func parseWithSecret(tokenStr, secret string) (*jwt.RegisteredClaims, error) {
	if secret == "" {
		return nil, errors.New("jwt secret not configured")
	}

	parser := jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	token, err := parser.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		// Extra safety: ensure HMAC family
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims.ExpiresAt != nil && time.Now().After(claims.ExpiresAt.Time) {
		return nil, errors.New("token expired")
	}

	return claims, nil
}

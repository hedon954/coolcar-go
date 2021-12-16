package shared_token

import (
	"fmt"
	"log"

	"github.com/dgrijalva/jwt-go"

	"crypto/rsa"
)

// JWTVerifier verifies access tokens.
type JWTVerifier struct {
	PublicKey *rsa.PublicKey
}

// Verify verifies a token and return accountID
func (v *JWTVerifier) Verify(token string) (string, error) {
	// parse token
	t, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return v.PublicKey, nil
	})
	if err != nil {
		return "", fmt.Errorf("cannot parse token, error: %v", err)
	}

	// check validation
	if !t.Valid {
		return "", fmt.Errorf("token is not valid")
	}

	// cast to *jwt.StandardClaims type
	clm, ok := t.Claims.(*jwt.StandardClaims)
	if !ok {
		return "", fmt.Errorf("token claim is not StandardClaim, error")
	}

	// check token is useful or not
	if err = clm.Valid(); err != nil {
		return "", fmt.Errorf("claim not valid, error: %v", err)
	}

	log.Printf("token is ok")
	log.Printf("accountID: %v", clm.Subject)

	// pass
	// subject == accountID
	return clm.Subject, nil
}

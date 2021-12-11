package token

import (
	"crypto/rsa"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTGen struct {
	privateKey *rsa.PrivateKey
	issuer     string
	nowFunc    func() time.Time
}

// NewJWTGen 工厂方法
func NewJWTGen(issuer string, privateKey *rsa.PrivateKey) *JWTGen {
	return &JWTGen{
		issuer:     issuer,
		nowFunc:    time.Now,
		privateKey: privateKey,
	}
}

// GenerateToken 生成 JWT
func (t *JWTGen) GenerateToken(accountID string, expire time.Duration) (string, error) {

	nowSec := t.nowFunc().Unix()

	// 调用 jwt 库生成 JWT
	// 采用 RS512 加密算法，非对称
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.StandardClaims{
		Issuer:    t.issuer,
		IssuedAt:  nowSec,
		ExpiresAt: nowSec + int64(expire.Seconds()),
		Subject:   accountID,
	})

	// 签名
	return token.SignedString(t.privateKey)
}

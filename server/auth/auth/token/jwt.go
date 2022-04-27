package token

import (
	"crypto/rsa"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTTokenGen struct {
	Issuer     string
	nowFunc    func() time.Time
	privateKey *rsa.PrivateKey
}

func NewJWTTokenGen(issuer string, privateKey *rsa.PrivateKey) *JWTTokenGen {
	return &JWTTokenGen{
		Issuer:     issuer,
		nowFunc:    time.Now,
		privateKey: privateKey,
	}
}

// GenerateToken 生成签名token
func (j *JWTTokenGen) GenerateToken(accountID string, expire time.Duration) (string, error) {
	tkn := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.StandardClaims{
		Issuer:    j.Issuer,
		IssuedAt:  j.nowFunc().Unix(), //颁发时间
		ExpiresAt: j.nowFunc().Unix() + int64(expire.Seconds()),
		Subject:   accountID,
	})
	signedString, err := tkn.SignedString(j.privateKey)
	if err != nil {
		return "", err
	}
	return signedString, nil
}

package token

import (
	"crypto/rsa"
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

// JWTTokenVerifier verify jwt access token
type JWTTokenVerifier struct {
	PublicKey *rsa.PublicKey
}

// Verify 用于生成accountID
func (v *JWTTokenVerifier) Verify(token string) (string, error) {
	//解析token
	parseWithClaims, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(*jwt.Token) (interface{}, error) {
		//验证签名
		return v.PublicKey, nil
	})
	if err != nil {
		return "", fmt.Errorf("cannot parse tonken:%v\n", err)
	}
	//
	if !parseWithClaims.Valid {
		return "", fmt.Errorf("token not valid\n")
	}
	//将解析出的类型放到新的claim中
	clm, ok := parseWithClaims.Claims.(*jwt.StandardClaims)
	if !ok {
		return "", fmt.Errorf("token claim is not standardclaim\n")

	}
	if err := clm.Valid(); err != nil {
		return "", fmt.Errorf("claim not valid:%v\n", err)
	}
	return clm.Subject, nil
}

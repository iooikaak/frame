package jwt

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"

	"github.com/iooikaak/frame/metadata"

	"github.com/iooikaak/frame/utils"
)

type RegisteredClaims = jwt.RegisteredClaims

type SigningMethod jwt.SigningMethod

var SigningMethodHS256 = jwt.SigningMethodHS256
var SigningMethodHS384 = jwt.SigningMethodHS384
var SigningMethodHS512 = jwt.SigningMethodHS512

var SigningMethodES256 = jwt.SigningMethodES256
var SigningMethodES384 = jwt.SigningMethodES384
var SigningMethodES512 = jwt.SigningMethodES512

var SigningMethodRS256 = jwt.SigningMethodRS256
var SigningMethodRS384 = jwt.SigningMethodRS384
var SigningMethodRS512 = jwt.SigningMethodRS512

var SigningMethodEdDSA = jwt.SigningMethodEdDSA

var ParseRSAPrivateKeyFromPEM = jwt.ParseRSAPrivateKeyFromPEM
var ParseRSAPublicKeyFromPEM = jwt.ParseRSAPublicKeyFromPEM
var ParseECPrivateKeyFromPEM = jwt.ParseECPrivateKeyFromPEM
var ParseECPublicKeyFromPEM = jwt.ParseECPublicKeyFromPEM
var ParseEdPrivateKeyFromPEM = jwt.ParseEdPrivateKeyFromPEM
var ParseEdPublicKeyFromPEM = jwt.ParseEdPublicKeyFromPEM

var NewNumericDate = jwt.NewNumericDate

type CustomClaims struct {
	UserInfo *UserInfo `json:"user_info"`
	RegisteredClaims
}

type UserInfo struct {
	UserId   int64  `json:"user_id"`
	RoleId   int64  `json:"role_id"`
	UserName string `json:"user_name"`
}

func GetJwtToken(signMethod SigningMethod, signKey interface{}, claims *CustomClaims) (r string, err error) {
	claims.ID = utils.UUIDv4()
	claims.NotBefore = NewNumericDate(time.Now())
	claims.IssuedAt = NewNumericDate(time.Now())

	token := jwt.NewWithClaims(signMethod, claims)
	r, err = token.SignedString(signKey)
	return
}

func ParseJwtToken(mySigningKey interface{}, tokenString string) (claims *CustomClaims, err error) {
	claims = &CustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})
	if err != nil {
		return
	}

	if token.Valid {
		return
	}
	return nil, fmt.Errorf("invalid token")
}

// nolint
func GetBeUserCustomClaims(ctx context.Context) *CustomClaims {
	claims := &CustomClaims{}
	c := ctx.Value(metadata.HeaderBeToken)
	claims = c.(*CustomClaims)
	return claims
}

// nolint
func GetBuUserCustomClaims(ctx context.Context) *CustomClaims {
	claims := &CustomClaims{}
	c := ctx.Value(metadata.HeaderBuToken)
	claims = c.(*CustomClaims)
	return claims
}

// nolint
func GetFrUserCustomClaims(ctx context.Context) *CustomClaims {
	claims := &CustomClaims{}
	c := ctx.Value(metadata.HeaderFrToken)
	claims = c.(*CustomClaims)
	return claims
}

package helper

import (
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/pajri/personal-backend/adapter/cerror"
	"github.com/pajri/personal-backend/global"
)

type JWTWrapper struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type JWTHelper struct {
}

func (j JWTHelper) CreateTokenPair(accessTokenParam, refreshTokenParam jwt.MapClaims) (*JWTWrapper, error) {
	var (
		token *JWTWrapper
		err   error
	)

	token = new(JWTWrapper)
	token.AccessToken, err = j.CreateToken(accessTokenParam)
	if err != nil {
		return nil, err
	}

	token.RefreshToken, err = j.CreateToken(refreshTokenParam)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (j JWTHelper) CreateToken(claims jwt.MapClaims) (string, error) {
	jwtWithClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := jwtWithClaims.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", cerror.NewAndPrintWithTag("CTH00", err, global.FRIENDLY_MESSAGE)
	}

	return token, nil
}

func (j JWTWrapper) ParseToken(tokenString string) (*jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

	if err != nil {
		return nil, cerror.NewAndPrintWithTag("PTH00", err, global.FRIENDLY_MESSAGE)
	}

	return &claims, err
}

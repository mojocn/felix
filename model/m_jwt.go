package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
)

var AppSecret = ""
var AppIss = "github.com/dejavuzhou/felix"
var ExpireTime = time.Hour * 24

type userStdClaims struct {
	jwt.StandardClaims
	*User
}

func (c userStdClaims) Valid() (err error) {
	err = c.StandardClaims.Valid()
	if err != nil {
		return err
	}
	if c.User.Id < 1 {
		return errors.New("invalid user in jwt")
	}
	return
}

func jwtGenerateToken(m *User) (string, error) {
	m.Password = ""
	expireAfterTime := time.Hour * 24
	expireTime := time.Now().Add(expireAfterTime)
	stdClaims := jwt.StandardClaims{
		ExpiresAt: expireTime.Unix(),
		IssuedAt:  time.Now().Unix(),
		Id:        fmt.Sprintf("%d", m.Id),
		Issuer:    AppIss,
	}

	uClaims := userStdClaims{
		StandardClaims: stdClaims,
		User:           m,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, uClaims)
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(AppSecret))
	if err != nil {
		logrus.WithError(err).Fatal("config is wrong, can not generate jwt")
	}
	return tokenString, err
}


//JwtParseUser
func JwtParseUser(tokenString string) (*User, error) {
	if tokenString == "" {
		return nil, errors.New("no token is found in Authorization Bearer")
	}
	claims := userStdClaims{}
	_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(AppSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims.VerifyExpiresAt(time.Now().Unix(), true) == false {
		return nil, errors.New("token is expired")
	}
	if !claims.VerifyIssuer(AppIss, true) {
		return nil, errors.New("token's issuer is wrong")
	}
	return claims.User, err
}

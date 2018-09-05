package services

import (
	"net/http"
	"github.com/kataras/iris/core/errors"
	"github.com/dgrijalva/jwt-go"
	"fmt"
	"fpdxIm/config"
)

// 用户认证
func Auth(req *http.Request) (string, error) {
	authorization := req.Header.Get("Authorization")
	if len(authorization) < 6 || authorization[0:6] != "bearer" {
		return "", errors.New("用户认证:authorization不合法")
	}
	token, err := jwt.Parse(authorization[0:6], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("用户认证:unexpected signing method: %v", token.Header["alg"])
		}
		return config.Config.App["JwtSecret"], nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["foo"], claims["nbf"])
		return claims["uid"].(string), nil
	} else {
		return "", fmt.Errorf("用户认证:%v", err.Error())
	}
}

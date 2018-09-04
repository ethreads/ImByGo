package services

import (
	"net/http"
	"github.com/kataras/iris/core/errors"
)

// 用户认证
func Auth(req *http.Request) (int, error) {
	if authorization := req.Header.Get("Authorization"); authorization == "" {
		return 0, errors.New("认证失败")
	}
	return 98047, nil
}

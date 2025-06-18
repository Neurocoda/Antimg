package models

import "errors"

var (
	ErrUserNotFound       = errors.New("用户不存在")
	ErrUserExists         = errors.New("用户已存在")
	ErrInvalidCredentials = errors.New("用户名或密码错误")
)

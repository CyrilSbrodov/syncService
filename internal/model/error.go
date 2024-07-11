package model

import "errors"

var (
	ErrorUserConflict = errors.New("user or email already exists")
	ErrorUserNotFound = errors.New("user not found")
)

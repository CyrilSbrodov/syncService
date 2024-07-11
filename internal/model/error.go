package model

import "errors"

var (
	ErrorClientConflict = errors.New("client name already exists")
)

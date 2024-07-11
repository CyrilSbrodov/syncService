package model

import "errors"

var (
	ErrorClientConflict = errors.New("client name already exists")
	ErrorNoClients      = errors.New("no one clients")
)

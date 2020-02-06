package domain

import "errors"

var (
	ErrRecordNotFound  = errors.New("record not found")
	ErrCloseConnection = errors.New("close connection")
)

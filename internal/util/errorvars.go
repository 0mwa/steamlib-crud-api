package util

import "errors"

var (
	ErrSWW    = errors.New("Something went wrong")
	ErrNI     = errors.New("Id is not int")
	ErrMethod = errors.New("405 - Wrong method")
)

package domain

import "errors"

var (
	ErrNotFound      = errors.New("not found")
	ErrUnknownType   = errors.New("unknown type")
	ErrInvalidType   = errors.New("bad type of obj")
	ErrAlreadyExists = errors.New("obj already exist")
	ErrForeignKey    = errors.New("sub obj dont exist")
)

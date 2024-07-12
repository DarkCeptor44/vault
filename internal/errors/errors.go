package vaultErrors

import "errors"

var (
	ErrInvalidKey   = errors.New("corrupt data or incorrect key")
	ErrDataTooShort = errors.New("data too short")
)

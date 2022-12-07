package gserror

import "github.com/jfy0o0/goStealer/errors/gscode"

// iCode is the interface for Code feature.
type iCode interface {
	Error() string
	Code() gscode.Code
}

// iStack is the interface for Stack feature.
type iStack interface {
	Error() string
	Stack() string
}

// iCause is the interface for Cause feature.
type iCause interface {
	Error() string
	Cause() error
}

// iCurrent is the interface for Current feature.
type iCurrent interface {
	Error() string
	Current() error
}

// iNext is the interface for Next feature.
type iNext interface {
	Error() string
	Next() error
}

package gherror

import "goHero/errors/ghcode"

// iCode is the interface for Code feature.
type iCode interface {
	Error() string
	Code() ghcode.Code
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

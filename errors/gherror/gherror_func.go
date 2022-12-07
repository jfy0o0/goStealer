package gherror

import "goHero/errors/ghcode"

// Code returns the error code of current error.
// It returns CodeNil if it has no error code neither it does not implement interface Code.
func Code(err error) ghcode.Code {
	if err == nil {
		return ghcode.CodeNil
	}
	if e, ok := err.(iCode); ok {
		return e.Code()
	}
	if e, ok := err.(iNext); ok {
		return Code(e.Next())
	}
	return ghcode.CodeNil
}

// Cause returns the root cause error of `err`.
func Cause(err error) error {
	if err == nil {
		return nil
	}
	if e, ok := err.(iCause); ok {
		return e.Cause()
	}
	if e, ok := err.(iNext); ok {
		return Cause(e.Next())
	}
	return err
}

// Stack returns the stack callers as string.
// It returns the error string directly if the `err` does not support stacks.
func Stack(err error) string {
	if err == nil {
		return ""
	}
	if e, ok := err.(iStack); ok {
		return e.Stack()
	}
	return err.Error()
}

// Current creates and returns the current level error.
// It returns nil if current level error is nil.
func Current(err error) error {
	if err == nil {
		return nil
	}
	if e, ok := err.(iCurrent); ok {
		return e.Current()
	}
	return err
}

// Next returns the next level error.
// It returns nil if current level error or the next level error is nil.
func Next(err error) error {
	if err == nil {
		return nil
	}
	if e, ok := err.(iNext); ok {
		return e.Next()
	}
	return nil
}

// HasStack checks and returns whether `err` implemented interface `iStack`.
func HasStack(err error) bool {
	_, ok := err.(iStack)
	return ok
}

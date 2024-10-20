package fake

import "errors"

var (
	ErrFakeNotInitialized   = errors.New("fake not initialized")
	ErrMethodNotInitialized = errors.New("method not initialized")
	ErrExpectationsMissing  = errors.New("expectations missing")
	ErrContextNil           = errors.New("received nil context")
	ErrArgumentMismatch     = errors.New("argument mismatch")
	ErrReturnMissing        = errors.New("missing return value(s) for expectation")
)

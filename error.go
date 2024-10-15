package fake

import "errors"

var (
	ErrFakeNotInitialized   = errors.New("fake not initialized")
	ErrMethodNotInitialized = errors.New("method not initialized")
	ErrContextNil           = errors.New("received nil context")
	ErrExpectationsMissing  = errors.New("expectations missing")
	ErrArgumentMismatch     = errors.New("argument mismatch")
	ErrReturnMissing        = errors.New("missing return value(s) for expectation")
)

var (
	errFakeNotInitialized   = "fake: " + ErrFakeNotInitialized.Error()
	errMethodNotInitialized = "fake: '%s.%s': " + ErrMethodNotInitialized.Error()
	errContextNil           = "fake: '%s.%s': " + ErrContextNil.Error()
	errContextCancel        = "fake: '%s.%s': %s"
	errExpectationsMissing  = "fake: '%s.%s': expectations missing: called '%d' time(s), '%d' expectation(s) registered"
	errArgumentMismatch     = "fake: '%s.%s': " + ErrArgumentMismatch.Error() + ": '%s': got '%+v', want '%+v'"
	errReturnMissing        = "fake: '%s.%s': " + ErrReturnMissing.Error() + ": '%d'"
)

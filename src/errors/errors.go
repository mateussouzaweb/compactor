package errors

import (
	"errors"
)

// Executes function and appends its error to the main error pointer
func Join(errPtr *error, errFunc func() error) {
	if err := errFunc(); err != nil {
		*errPtr = errors.Join(*errPtr, err)
	}
}

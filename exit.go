// Package exit implements an error-based alternative to os.Exit.
package exit

import (
	"fmt"
	"os"
)

var terminate = os.Exit

// Wrap wraps the given error with the given exit code.
func Wrap(code int, err error) error {
	return &wrapper{
		error: err,
		code:  code,
	}
}

// Wrapf behaves like fmt.Errorf while at the same time wrapping the returned
// error with the given exit code.
func Wrapf(code int, format string, a ...interface{}) error {
	return Wrap(code, fmt.Errorf(format, a...))
}

type wrapper struct {
	error
	code int
}

func (w *wrapper) Unwrap() error {
	return w.error
}

// Is reports whether the given error carries an exit code.
func Is(err error) (has bool) {
	_, has = err.(*wrapper)
	return
}

// HasCode reports whether the given error carries the given exit code.
func HasCode(err error, code int) bool {
	e, ok := err.(*wrapper)
	return ok && e.code == code
}

// With calls os.Exit with an appropriate exit code for the given error.
//
// If the given error is nil, os.Exit will be called with 0.
//
// Alternatively, and should the given error not implement the Error interface,
// os.Exit will be called with 1.
//
// In any other case, os.Exit will be called with the result of the given
// error's ExitCode function.
func With(err error) {
	var code int

	if err != nil {
		if ec, ok := err.(*wrapper); !ok {
			code = 1
		} else {
			code = ec.code
		}
	}

	terminate(code)
}

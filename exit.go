// Package exit implements an error-based alternative to os.Exit.
package exit

import (
	"errors"
	"fmt"
	"os"
)

var terminate = os.Exit

// Wrap wraps the err with code.
func Wrap(code int, err error) error {
	return &wrapper{
		error: err,
		int:   code,
	}
}

// Wrapf behaves like fmt.Errorf while at the same time wrapping the returned
// error with code.
func Wrapf(code int, format string, a ...interface{}) error {
	return Wrap(code, fmt.Errorf(format, a...))
}

type wrapper struct {
	error
	int
}

func (w *wrapper) Unwrap() error {
	return w.error
}

// Is reports whether any error in err's chain is an exit one.
func Is(err error) bool {
	var w *wrapper
	return errors.As(err, &w)
}

// Code returns the first exit code err's chain carries.
//
// In case err's chain does not carry an exit code, carries will be unset.
func Code(err error) (code int, carries bool) {
	var w *wrapper
	if carries = errors.As(err, &w); carries {
		code = w.int
	}

	return
}

// Carries reports whether the first exit error in err's chain carries code.
func Carries(err error, code int) bool {
	c, has := Code(err)
	return has && code == c
}

// With calls os.Exit with an exit code appropriate for err.
//
// Should err be nil, os.Exit will be called with 0.
//
// Alternatively, should err or its chain not carry an exit code, os.Exit will
// be called with 1.
//
// In all other cases, os.Exit will be called with the code err carries.
func With(err error) {
	var code int
	if err != nil {
		var carries bool
		if code, carries = Code(err); !carries {
			code = 1
		}
	}

	terminate(code)
}

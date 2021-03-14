// Package exit implements an error-based alternative to os.Exit.
package exit

import (
	"fmt"
	"os"
)

// Error wraps the set of exit errors.
type Error interface {
	error

	// ExitCode reports the exit code the error denotes.
	ExitCode() int
}

var terminate = os.Exit

// Code provides a built-in implementation of Error.
type Code int

// Error implements error for Code.
func (c Code) Error() string { return fmt.Sprintf("exit.Code(%d)", c) }

// ExitCode implements Error for code.
func (c Code) ExitCode() int { return int(c) }

var _ Error = Code(3) // compile time check that code implements Error

// With calls os.Exit with an appropriate for the given error exit code.
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
		if ec, ok := err.(Error); !ok {
			code = 1
		} else {
			code = ec.ExitCode()
		}
	}

	terminate(code)
}

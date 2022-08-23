package exit

import (
	"errors"
	"io"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWith(t *testing.T) {
	cases := []struct {
		err error
		exp int
	}{
		0: {nil, 0},
		1: {io.EOF, 1},
		2: {errors.New("error"), 1},
		3: {Wrap(2, io.EOF), 2},
		4: {Wrapf(4, "failed doing something: %w", io.EOF), 4},
	}

	for caseIndex := range cases {
		kase := cases[caseIndex]

		t.Run(strconv.Itoa(caseIndex), func(t *testing.T) {
			var got int
			defer capture(&got)()

			With(kase.err)
			assert.Equal(t, kase.exp, got)
		})
	}
}

func capture(into *int) func() {
	current := terminate

	terminate = func(code int) {
		*into = code
	}

	return func() { terminate = current }
}

func TestUnwrap(t *testing.T) {
	exp := errors.New("some error")

	assert.Same(t, exp, errors.Unwrap(Wrap(12, exp)))
}

func TestIs(t *testing.T) {
	cases := []struct {
		err error
		exp bool
	}{
		{nil, false},
		{io.EOF, false},
		{Wrap(2, io.EOF), true},
		{Wrap(3, Wrap(1, io.ErrUnexpectedEOF)), true},
	}

	for caseIndex := range cases {
		kase := cases[caseIndex]

		t.Run(strconv.Itoa(caseIndex), func(t *testing.T) {
			assert.Equal(t, kase.exp, Is(kase.err))
		})
	}
}

func TestCarries(t *testing.T) {
	const code = 5

	cases := []struct {
		err error
		exp bool
	}{
		{nil, false},
		{io.EOF, false},
		{Wrap(code, io.EOF), true},
		{Wrapf(code, "failed with: %w", io.EOF), true},
		{Wrap(code-1, io.EOF), false},
		{Wrapf(code+1, "failed with: %w", io.EOF), false},
	}

	for caseIndex := range cases {
		kase := cases[caseIndex]

		t.Run(strconv.Itoa(caseIndex), func(t *testing.T) {
			assert.Equal(t, kase.exp, Carries(kase.err, code))
		})
	}
}

func TestFail(t *testing.T) {
	assert.False(t, true)
}

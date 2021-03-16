package exit

import (
	"errors"
	"io"
	"strconv"
	"testing"
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
			if got != kase.exp {
				t.Errorf("got %d, expected %d", got, kase.exp)
			}
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
	err := Wrap(12, exp)

	if !errors.Is(err, exp) {
		t.Fatal("is reported false")
	}
}

func TestIs(t *testing.T) {
	cases := []struct {
		err error
		exp bool
	}{
		{nil, false},
		{io.EOF, false},
		{Wrap(2, io.EOF), true},
	}

	for caseIndex := range cases {
		kase := cases[caseIndex]

		t.Run(strconv.Itoa(caseIndex), func(t *testing.T) {
			if got := Is(kase.err); got != kase.exp {
				t.Errorf("got %t, expected %t", got, kase.exp)
			}
		})
	}
}

func TestHasCode(t *testing.T) {
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
			if got := HasCode(kase.err, code); got != kase.exp {
				t.Errorf("got %t, expected %t", got, kase.exp)
			}
		})
	}
}

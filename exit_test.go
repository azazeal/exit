package exit

import (
	"errors"
	"io"
	"reflect"
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
			if kase.exp != got {
				t.Errorf("\nexp: %q\ngot: %q", kase.exp, got)
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
	p1 := reflect.ValueOf(exp).Pointer()

	got := errors.Unwrap(Wrap(12, exp))
	p2 := reflect.ValueOf(got).Pointer()

	if p1 != p2 {
		t.Fatalf("\nexp: %p %#v\ngot: %p %#v", exp, exp, got, got)
	}
}

func TestIs(t *testing.T) {
	cases := []struct {
		err error
		exp bool
	}{
		0: {nil, false},
		1: {io.EOF, false},
		2: {Wrap(2, io.EOF), true},
		3: {Wrap(3, Wrap(1, io.ErrUnexpectedEOF)), true},
	}

	for caseIndex := range cases {
		kase := cases[caseIndex]

		t.Run(strconv.Itoa(caseIndex), func(t *testing.T) {
			if got := Is(kase.err); kase.exp != got {
				t.Errorf("\nexp: %t\ngot: %t", kase.exp, got)
			}
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
			if got := Carries(kase.err, code); kase.exp != got {
				t.Errorf("\nexp: %t\ngot: %t", kase.exp, got)
			}
		})
	}
}

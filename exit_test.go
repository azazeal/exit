package exit

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"testing"
	"testing/quick"
)

var cases = []struct {
	err error
	exp int
}{
	0: {nil, 0},
	1: {io.EOF, 1},
	2: {errors.New("error"), 1},
	3: {Code(0), 0},
	4: {Code(1), 1},
	5: {Code(-1), -1},
}

func Test(t *testing.T) {
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

func TestCodeError(t *testing.T) {
	fn := func(v Code) bool {
		exp := fmt.Sprintf("exit.Code(%d)", v)

		return v.Error() == exp
	}

	if err := quick.Check(fn, nil); err != nil {
		t.Fatal(err)
	}
}

package literal

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/google/go-cmp/cmp"
	"testing"
)

type evalTest struct {
	in  string
	out interface{}
	err bool
}

func Test(t *testing.T) {
	cases := []evalTest{
		{"42", int64(42), false},
		{"-42", int64(42), false},
	}
	for _, c := range cases {
		runEvalTest(t, c)
	}
}

func runEvalTest(t *testing.T, test evalTest) {
	res, err := Eval(test.in)
	if err != nil {
		if test.err {
			return
		}
		t.Fatal(err, spew.Sprint(test))
	}
	if test.err {
		t.Fatal("expected error", spew.Sprint(test))
	}
	if !cmp.Equal(res, test.out) {
		t.Fatal("got:", spew.Sprint(res), "expected:", spew.Sprint(test.out), spew.Sprint(test))
	}
}

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

func TestCorrect(t *testing.T) {
	cases := []evalTest{
		{"42", int64(42), false},
		{"-42", int64(-42), false},
		{"+42", int64(42), false},
		{"42.", 42., false},
		{"-42.", -42., false},
		{"+42.", 42., false},
		{"1j", complex(0, 1), false},
		{"-42j", complex(0, -42), false},
		{"+42j", complex(0, 42), false},
		{"'foo'", "foo", false},
		{"\"foo\"", "foo", false},
		{"...", Ellipsis{}, false},
		{"[]", []interface{}{}, false},
		{"[1, 'foo', ..., 'bar']", []interface{}{int64(1), "foo", Ellipsis{}, "bar"}, false},
		{"(1, 'foo', ..., 'bar')", []interface{}{int64(1), "foo", Ellipsis{}, "bar"}, false},
		{"(1, 'foo', ..., 'bar',)", []interface{}{int64(1), "foo", Ellipsis{}, "bar"}, false},
		{"{1, 'foo', ..., 'bar',}", Set{[]interface{}{int64(1), "foo", Ellipsis{}, "bar"}}, false},
		{"{1, 1}", Set{[]interface{}{int64(1)}}, false},
		{"{}", Dict{}, false},
		{"{'foo': 'bar'}", Dict{[]interface{}{"foo"}, []interface{}{"bar"}}, false},
		{"{'foo': 'bar', 'foo': 'foo'}", Dict{[]interface{}{"foo"}, []interface{}{"foo"}}, false},
		{"{('foo',1): 'bar'}", Dict{[]interface{}{[]interface{}{"foo", int64(1)}}, []interface{}{"bar"}}, false},
		{"None", nil, false},
		{"True", true, false},
		{"False", false, false},
		{"b'foo'", []byte("foo"), false},
		//{"1 + 1j", complex(1, 1), false}, // TODO currently only imaginary numbers are supported
	}
	for _, c := range cases {
		runEvalTest(t, c)
	}
}

func TestWrong(t *testing.T) {
	cases := []string{
		// TODO
	}
	for _, c := range cases {
		runEvalTest(t, evalTest{c, nil, true})
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

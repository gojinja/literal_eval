package literal

import (
	"fmt"
	"github.com/gojinja/literal_eval/ast"
	"github.com/google/go-cmp/cmp"
)

// Set represents Python's dict. It can't be Go's map as in Python more elements can be keys.
type Set struct {
	Elems []interface{}
}

func newSet(keys []interface{}, rawKeys []ast.Expr) (res Set, err error) {
	for i, k := range keys {
		if !isHashable(rawKeys[i]) {
			return res, fmt.Errorf("set key is not hashable")
		}
		in := false
		for _, oK := range keys[i+1:] {
			if cmp.Equal(k, oK) {
				in = true
				break
			}
		}
		if !in {
			res.Elems = append(res.Elems, k)
		}
	}
	return
}

func (s *Set) Has(key interface{}) bool {
	for _, k := range s.Elems {
		if cmp.Equal(key, k) {
			return true
		}
	}
	return false
}

func (s *Set) Len() int {
	return len(s.Elems)
}

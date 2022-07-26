package literal

import (
	"fmt"
	"github.com/gojinja/literal_eval/ast"
	"github.com/google/go-cmp/cmp"
)

// Dict represents Python's dict. It can't be Go's map as in Python more elements can be keys.
type Dict struct {
	Keys   []interface{}
	Values []interface{}
}

func newDict(keys []interface{}, values []interface{}, rawKeys []ast.Expr) (res Dict, err error) {
	for i, k := range keys {
		if !isHashable(rawKeys[i]) {
			return res, fmt.Errorf("dict key is not hashable")
		}
		in := false
		for _, oK := range keys[i+1:] {
			if cmp.Equal(k, oK) {
				in = true
				break
			}
		}
		if !in {
			res.Keys = append(res.Keys, k)
			res.Values = append(res.Values, values[i])
		}
	}
	return
}

func (d *Dict) Get(key interface{}) (interface{}, bool) {
	for i, k := range d.Keys {
		if cmp.Equal(key, k) {
			return d.Values[i], true
		}
	}
	return nil, false
}

func (d *Dict) Len() int {
	return len(d.Keys)
}

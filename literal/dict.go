package literal

import "github.com/google/go-cmp/cmp"

// Dict represents Python's dict. It can't be Go's map as in Python more elements can be keys.
type Dict struct {
	Keys  []interface{}
	Value []interface{}
}

func newDict(keys []interface{}, values []interface{}) (res Dict) {
	for i, k := range keys {
		in := false
		for _, oK := range keys[i+1:] {
			if cmp.Equal(k, oK) {
				in = true
				break
			}
		}
		if !in {
			res.Keys = append(res.Keys, k)
			res.Value = append(res.Keys, values[i])
		}
	}
	return
}

func (d *Dict) Get(key interface{}) (interface{}, bool) {
	for i, k := range d.Keys {
		if cmp.Equal(key, k) {
			return d.Value[i], true
		}
	}
	return nil, false
}

func (d *Dict) Len() int {
	return len(d.Keys)
}

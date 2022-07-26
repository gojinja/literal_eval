// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package builtin_test

import (
	"testing"

	"github.com/gojinja/literal_eval/pytest"
)

func TestVm(t *testing.T) {
	pytest.RunTests(t, "tests")
}

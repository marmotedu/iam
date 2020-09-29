// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pumps

import (
	"testing"
)

func TestGetPumpByName(t *testing.T) {
	name := "dummy"
	pmpType, err := GetPumpByName(name)

	if err != nil || pmpType == nil {
		t.Fail()
	}

	name2 := "xyz"
	pmpType2, err2 := GetPumpByName(name2)

	if err2 == nil || pmpType2 != nil {
		t.Fail()
	}
}

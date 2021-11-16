// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package reflect

import (
	"reflect"
	"testing"
)

func TestGetObjFieldsMap(t *testing.T) {
	type Obj struct {
		A int
		B int
		C int
	}

	org := &Obj{
		A: 1,
		B: 2,
		C: 3,
	}

	m := GetObjFieldsMap(org, []string{})
	if !reflect.DeepEqual(m, map[string]interface{}{
		"A": 1,
		"B": 2,
		"C": 3,
	}) {
		t.Fatalf("not equal")
	}

	m = GetObjFieldsMap(org, []string{"A"})
	if !reflect.DeepEqual(m, map[string]interface{}{
		"A": 1,
	}) {
		t.Fatalf("not equal")
	}
}

func TestCopyObj(t *testing.T) {
	type Obj struct {
		A int
		B int
		C int
	}

	org := &Obj{
		A: 1,
		B: 2,
		C: 3,
	}

	des := &Obj{
		A: 4,
		B: 5,
		C: 6,
	}

	changed, err := CopyObj(org, des, []string{"A"})
	if err != nil {
		t.Fatalf(err.Error())
	}

	if !changed {
		t.Fatalf("expect changed")
	}

	if des.A != org.A {
		t.Fatalf("A not copy")
	}

	if des.B != 5 || des.C != 6 {
		t.Fatalf("B and C changed")
	}

	des.A = org.A
	changed, err = CopyObj(org, des, []string{"A"})
	if err != nil {
		t.Fatalf(err.Error())
	}

	if changed {
		t.Fatalf("expect not changed")
	}
}

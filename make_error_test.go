// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                /[make_error_test.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"testing"
)

// makeError(id uint32, a ...interface{}) error
//
// go test -run Test_makeError_*

func Test_makeError_1(t *testing.T) {
	for _, it := range []struct {
		id uint32
		a  []interface{}
		//
		expect string
	}{
		{0, nil, "ERROR 0x000000:"},
		{1, nil, "ERROR 0x000001:"},
		{1, []interface{}{"a", 1, 2, "c", 3}, "ERROR 0x000001: a 1 2 c 3"},
		{0xE12345, []interface{}{"failed"}, "ERROR 0xE12345: failed"},
		{0xE12345, []interface{}{"failed", 123}, "ERROR 0xE12345: failed 123"},
	} {
		err := makeError(it.id, it.a...)
		got := err.Error()
		if got != it.expect {
			t.Errorf("0xE60E1A"+" makeError(%X, %#v)"+
				"\n expect: %#v"+
				"\n    got: %#v",
				it.id, it.a, it.expect, got)
		}
	}
} //                                                            Test_makeError_1

func Test_makeError_(t *testing.T) {
	a := makeError(0xEE5D1E, "the error message")
	b := makeError(0xE4C1F0, a.Error())
	c := makeError(0xED10D5, b.Error())
	got := c.Error()
	if got != "ERROR 0x"+"ED10D5: the error message" {
		t.Error("0xEE38BE")
	}
} //                                                            Test_makeError_2

// end

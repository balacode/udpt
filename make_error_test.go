// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                /[make_error_test.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"testing"
)

// makeError(id uint32, args ...interface{}) error
//
// go test -run Test_makeError_
//
func Test_makeError_(t *testing.T) {
	for _, it := range []struct {
		id   uint32
		args []interface{}
		//
		expect string
	}{
		{0, nil, "ERROR 0x000000:"},
		{1, nil, "ERROR 0x000001:"},
		{1, []interface{}{"a", 1, 2, "c", 3}, "ERROR 0x000001: a 1 2 c 3"},
		{0xE12345, []interface{}{"failed"}, "ERROR 0xE12345: failed"},
		{0xE12345, []interface{}{"failed", 123}, "ERROR 0xE12345: failed 123"},
	} {
		err := makeError(it.id, it.args...)
		got := err.Error()
		if got != it.expect {
			t.Errorf("makeError(%X, %#v)"+
				"\n expect: %#v"+
				"\n    got: %#v",
				it.id, it.args, it.expect, got)
		}
	}
} //                                                             Test_makeError_

// end

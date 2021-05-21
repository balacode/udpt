// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                  /[get_hash_test.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"fmt"
	"strings"
	"testing"
)

// getHash(data []byte) []byte
//
// go test -run _getHash_
//
func Test_getHash_getHash_(t *testing.T) {
	{
		got := fmt.Sprintf("%X", getHash(nil))
		if got != "E3B0C44298FC1C149AFBF4C8996FB924"+
			"27AE41E4649B934CA495991B7852B855" {
			t.Error("0xE00CE9")
		}
	}
	func() {
		defer func() {
			msg := fmt.Sprint(recover())
			if !strings.Contains(msg, "badHash.Write") {
				t.Error("0xE1A10A")
			}
		}()
		_ = getHashDI(nil, &badHash{})
		t.Error("0xE20F56") // previous line should panic and never reach here
	}()
	func() {
		defer func() {
			msg := fmt.Sprint(recover())
			if !strings.Contains(msg, "invalid hash size") {
				t.Error("0xE4A45E")
			}
		}()
		_ = getHashDI([]byte{1, 2, 3}, &badHash{})
		t.Error("0xEE2EE8") // previous line should panic and never reach here
	}()
} //                                                       Test_getHash_getHash_

// -----------------------------------------------------------------------------

type badHash struct{}

func (*badHash) BlockSize() int       { return 64 }
func (*badHash) Reset()               {}
func (*badHash) Size() int            { return 32 }
func (*badHash) Sum(in []byte) []byte { return nil }

func (*badHash) Write(data []byte) (int, error) {
	var err error
	if data == nil {
		err = makeError(0xED9CD4, "badHash.Write")
	}
	return len(data), err
} //                                                                       Write

// end

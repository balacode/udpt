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

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
// getHash(data []byte) []byte
//
// go test -run Test_getHash_*

func Test_getHash_1(t *testing.T) {
	got := fmt.Sprintf("%X", getHash(nil))
	if got != "E3B0C44298FC1C149AFBF4C8996FB924"+
		"27AE41E4649B934CA495991B7852B855" {
		t.Error("0xE00CE9")
	}
} //                                                              Test_getHash_1

func Test_getHash_2(t *testing.T) {
	func() {
		defer func() {
			msg := fmt.Sprint(recover())
			if !strings.Contains(msg, "failed mockHash.Write") {
				t.Error("0xE1A10A")
			}
		}()
		_ = getHashDI(nil, &mockHash{failWrite: true})
		t.Error("0xE20F56") // previous line should panic and never reach here
	}()
} //                                                              Test_getHash_2

func Test_getHash_3(t *testing.T) {
	func() {
		defer func() {
			msg := fmt.Sprint(recover())
			if !strings.Contains(msg, "invalid hash size") {
				t.Error("0xE4A45E")
			}
		}()
		_ = getHashDI([]byte{1, 2, 3}, &mockHash{})
		t.Error("0xEE2EE8") // previous line should panic and never reach here
	}()
} //                                                              Test_getHash_3

// -----------------------------------------------------------------------------

// mockHash is a mock hash.Hash with a Write method you can make fail.
type mockHash struct{ failWrite bool }

func (*mockHash) BlockSize() int { return 64 }

func (*mockHash) Reset() {}

func (*mockHash) Size() int { return 32 }

func (*mockHash) Sum(in []byte) []byte { return nil }

func (mk *mockHash) Write(data []byte) (int, error) {
	var err error
	if mk.failWrite {
		return 0, makeError(0xED9CD4, "failed mockHash.Write")
	}
	return len(data), err
} //                                                                       Write

// end

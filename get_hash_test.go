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

// test if zero-length hash is correct
func Test_getHash_1(t *testing.T) {
	got := fmt.Sprintf("%X", getHash(nil))
	if got != "E3B0C44298FC1C149AFBF4C8996FB924"+
		"27AE41E4649B934CA495991B7852B855" {
		t.Error("0xE00CE9")
	}
}

// must panic if Hash.Write fails (this should never happen)
func Test_getHash_2(t *testing.T) {
	func() {
		defer func() {
			s := fmt.Sprint(recover())
			if !strings.Contains(s, "failed mockHash.Write") {
				t.Error("0xE1A10A")
			}
		}()
		_ = getHashDI(nil, &mockHash{failWrite: true})
		t.Error("0xE20F56") // previous line should panic and never reach here
	}()
}

// must panic if Hash.Sum returns a result that is not 32 bytes long
// (this should never happen)
func Test_getHash_3(t *testing.T) {
	func() {
		defer func() {
			s := fmt.Sprint(recover())
			if !strings.Contains(s, "invalid hash size") {
				t.Error("0xE4A45E")
			}
		}()
		_ = getHashDI([]byte{1, 2, 3}, &mockHash{})
		t.Error("0xEE2EE8") // previous line should panic and never reach here
	}()
}

// -----------------------------------------------------------------------------

// mockHash is a mock SHA-256 hash that implements hash.Hash.
type mockHash struct{ failWrite bool }

// BlockSize must return 64: the underlying block size of SHA-256 hash.
func (*mockHash) BlockSize() int { return 64 }

// Reset resets the hash to its initial state.
// It does nothing here, included just to implement the interface.
func (*mockHash) Reset() {}

// Size must return 32: the number of bytes Sum will return.
func (*mockHash) Size() int { return 32 }

// Sum appends the current hash to b and returns the resulting slice. It
// should append 32 bytes, but you can set sumBytes to test wrong values.
func (*mockHash) Sum(in []byte) []byte { return nil }

// Write adds more data to the running hash using io.Writer and never errors.
func (mk *mockHash) Write(data []byte) (int, error) {
	var err error
	if mk.failWrite {
		return 0, makeError(0xED9CD4, "failed mockHash.Write")
	}
	return len(data), err
}

// end

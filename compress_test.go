// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                  /[compress_test.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"bytes"
	"strings"
	"testing"
)

// go test --run Test_compress_uncompress_
func Test_compress_uncompress_(t *testing.T) {
	input := []byte(strings.Repeat(
		"The quick brown fox jumps over the lazy dog!", 7,
	))
	comp, err := compress(input)
	if err != nil {
		t.Error("0xE26CD5 compress() failed:", err)
	}
	uncomp, err := uncompress(comp)
	if err != nil {
		t.Error("0xE38FD2 uncompress() failed:", err)
	}
	if !bytes.Equal(input, uncomp) {
		t.Error("0xE58D5D compress()/uncompress() failed")
	}
} //                                                   Test_compress_uncompress_

// end

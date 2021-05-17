// -----------------------------------------------------------------------------
// github.com/balacode/udpt                           /[zlib_compressor_test.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"bytes"
	"strings"
	"testing"
)

// go test -run _ZLibCompressor_
//
func Test_ZLibCompressor_(t *testing.T) {
	zc := zlibCompressor{}
	input := []byte(strings.Repeat(
		"The quick brown fox jumps over the lazy dog!", 7,
	))
	comp, err := zc.Compress(input)
	if err != nil {
		t.Error("0xE26CD5 Compress() failed:", err)
	}
	uncomp, err := zc.Uncompress(comp)
	if err != nil {
		t.Error("0xE38FD2 Uncompress() failed:", err)
	}
	if !bytes.Equal(input, uncomp) {
		t.Error("0xEB0E80 Compress()/Uncompress() failed")
	}
} //                                                        Test_ZLibCompressor_

// end

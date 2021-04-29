// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                       /[compress.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"bytes"
	"compress/zlib"
)

// compress compresses data using zlib and returns the compressed bytes.
// If there was an error compressing, returns nil and the error description.
func compress(data []byte) ([]byte, error) {
	var cbuf bytes.Buffer
	wr := zlib.NewWriter(&cbuf)
	_, err := wr.Write(data)
	if err != nil {
		defer func() {
			err2 := wr.Close()
			if err2 != nil {
				_ = logError(0xE0A6F2, "(Close):", err2)
			}
		}()
		return nil, logError(0xE5F7D3, "(Write):", err)
	}
	err = wr.Close()
	if err != nil {
		return nil, logError(0xE39D8B, "(Close):", err)
	}
	ret := cbuf.Bytes()
	return ret, nil
} //                                                                    compress

// end

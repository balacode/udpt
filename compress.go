// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                       /[compress.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
)

// compress compresses data using zlib and returns the compressed bytes.
// If there was an error compressing, returns nil and the error description.
func compress(data []byte) ([]byte, error) {
	var cbuf bytes.Buffer
	wr := zlib.NewWriter(&cbuf)
	_, err := wr.Write(data)
	if err != nil {
		defer func() {
			_ = wr.Close()
		}()
		return nil, makeError(0xE5F7D3, "(Write):", err)
	}
	err = wr.Close()
	if err != nil {
		return nil, makeError(0xE39D8B, "(Close):", err)
	}
	ret := cbuf.Bytes()
	//
	// write the uncompressed size after the compressed data
	nc := make([]byte, 4)
	binary.LittleEndian.PutUint32(nc, uint32(len(data)))
	ret = append(ret, nc...)
	//
	return ret, nil
} //                                                                    compress

// end

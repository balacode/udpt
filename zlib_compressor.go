// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                    /[compression.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"io"
)

// zlibCompressor implements the Compression interface to compress and
// uncompress byte slices using the zlib format specified in RFC-1950.
type zlibCompressor struct{}

// Compress compresses 'data' using zlib and returns the compressed bytes.
// If there was an error, returns nil and the error value.
func (*zlibCompressor) Compress(data []byte) ([]byte, error) {
	var cbuf bytes.Buffer
	wr := zlib.NewWriter(&cbuf)
	_, err := wr.Write(data)
	if err != nil {
		defer func() {
			_ = wr.Close()
		}()
		return nil, makeError(0xE5F7D3, err)
	}
	err = wr.Close()
	if err != nil {
		return nil, makeError(0xE39D8B, err)
	}
	ret := cbuf.Bytes()
	//
	// write the size of uncompressed data after the compressed bytes
	nc := make([]byte, 4)
	binary.LittleEndian.PutUint32(nc, uint32(len(data)))
	ret = append(ret, nc...)
	//
	return ret, nil
} //                                                                    Compress

// Uncompress uncompresses 'compressed' bytes using zlib and returns the
// uncompressed bytes. If there was an error, returns nil and the error value.
func (*zlibCompressor) Uncompress(compressed []byte) ([]byte, error) {
	nc := len(compressed)
	if len(compressed) <= 4 {
		return nil, makeError(0xE8A8A9, "invalid compressed")
	}
	// read the uncompressed data size (stored at the end of compressed bytes)
	// to know the number of bytes to allocate for the result
	nu := int64(binary.LittleEndian.Uint32(compressed[nc-4:]))
	compressed = compressed[:nc-4]
	//
	reader, err := zlib.NewReader(bytes.NewReader(compressed))
	if err != nil {
		return nil, makeError(0xE54F4B, err)
	}
	buf := bytes.NewBuffer(make([]byte, 0, nu))
	_, err = io.CopyN(buf, reader, nu)
	if err != nil {
		return nil, makeError(0xE6A29D, err)
	}
	err = reader.Close()
	if err != nil {
		return nil, makeError(0xE45AF8, err)
	}
	ret := buf.Bytes()
	return ret, nil
} //                                                                  Uncompress

// end

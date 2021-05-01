// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                     /[uncompress.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"io"
)

// uncompress uncompresses compressedData using zlib and returns
// the uncompressed bytes. If compressedData is invalid,
// returns nil and the error description.
func uncompress(compressedData []byte) ([]byte, error) {
	nc := len(compressedData)
	if len(compressedData) <= 4 {
		return nil, logError(0xE8A8A9, "invalid compreseData")
	}
	// extract the uncompressed size from the end of compressedData
	// to know the exact number of bytes to allocate for the result
	nu := int64(binary.LittleEndian.Uint32(compressedData[nc-4:]))
	compressedData = compressedData[:nc-4]
	//
	reader, err := zlib.NewReader(bytes.NewReader(compressedData))
	if err != nil {
		return nil, logError(0xE54F4B, "(NewReader):", err)
	}
	buf := bytes.NewBuffer(make([]byte, 0, nu))
	_, err = io.CopyN(buf, reader, nu)
	if err != nil {
		return nil, logError(0xE6A29D, "(CopyN):", err)
	}
	err = reader.Close()
	if err != nil {
		return nil, logError(0xE45AF8, "(Close):", err)
	}
	ret := buf.Bytes()
	return ret, nil
} //                                                                  uncompress

// end

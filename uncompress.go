// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                     /[uncompress.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"bytes"
	"compress/zlib"
	"io"
)

// uncompress uncompresses compressedData using zlib and returns
// the uncompressed bytes. If compressedData is invalid,
// returns nil and the error description.
func uncompress(compressedData []byte) ([]byte, error) {
	reader, err := zlib.NewReader(bytes.NewReader(compressedData))
	if err != nil {
		return nil, logError(0xE54F4B, "(NewReader):", err)
	}
	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, reader)
	if err != nil {
		return nil, logError(0xE6A29D, "(Copy):", err)
	}
	err = reader.Close()
	if err != nil {
		return nil, logError(0xE45AF8, "(Close):", err)
	}
	ret := buf.Bytes()
	return ret, nil
} //                                                                  uncompress

// end

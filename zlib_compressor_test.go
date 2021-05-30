// -----------------------------------------------------------------------------
// github.com/balacode/udpt                           /[zlib_compressor_test.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

// to run all tests in this file:
// go test -v -run Test_zlibCompressor_*

// -----------------------------------------------------------------------------

// must succeed compressing and uncompressing data:
func Test_zlibCompressor_1(t *testing.T) {
	comp := zCompress(t)
	zc := zlibCompressor{}
	uncomp, err := zc.Uncompress(comp) // must succeed
	if err != nil {
		t.Error("0xE38FD2", err)
	}
	if !bytes.Equal(zInput(), uncomp) {
		t.Error("0xEB0E80", "data corrupted")
	}
}

// uncompressing less than 4 bytes must fail:
func Test_zlibCompressor_2(t *testing.T) {
	zc := zlibCompressor{}
	comp, err := zc.Uncompress([]byte{1, 2, 3})
	if comp != nil {
		t.Error("0xEE55FC")
	}
	if !matchError(err, "invalid 'comp'") {
		t.Error("0xE4AB37", "wrong error:", err)
	}
}

// uncompressing must fail when wr.Write() fails in compressDI():
func Test_zlibCompressor_3(t *testing.T) {
	zc := zlibCompressor{}
	cbuf := bytes.Buffer{}
	wrc := &mockWriteCloser{failWrite: true}
	comp, err := zc.compressDI(zInput(), wrc, &cbuf)
	if comp != nil {
		t.Error("0xE27BB2")
	}
	if !matchError(err, "failed mockWriteCloser.Write") {
		t.Error("0xEE9C54", "wrong error:", err)
	}
}

// uncompressing must fail when wr.Close() fails in compressDI():
func Test_zlibCompressor_4(t *testing.T) {
	zc := zlibCompressor{}
	cbuf := bytes.Buffer{}
	wrc := &mockWriteCloser{failClose: true}
	comp, err := zc.compressDI(zInput(), wrc, &cbuf)
	if comp != nil {
		t.Error("0xE6F6A5")
	}
	if !matchError(err, "failed mockWriteCloser.Close") {
		t.Error("0xE4DF92", "wrong error:", err)
	}
}

// uncompressing must fail when reader.Read or io.Copy fail in uncompressDI():
func Test_zlibCompressor_5(t *testing.T) {
	comp, zc := zCompress(t), zlibCompressor{}
	newMockReadCloser := func(io.Reader) (io.ReadCloser, error) {
		return &mockReadCloser{failRead: true}, nil
	}
	uncomp, err := zc.uncompressDI(comp, newMockReadCloser)
	if uncomp != nil {
		t.Error("0xE3DA4F")
	}
	if !matchError(err, "failed mockReadCloser.Read") {
		t.Error("0xE81C62", "wrong error:", err)
	}
}

// uncompressing must fail when reader.Close() fails in uncompressDI():
func Test_zlibCompressor_6(t *testing.T) {
	comp, zc := zCompress(t), zlibCompressor{}
	newMockReadCloser := func(io.Reader) (io.ReadCloser, error) {
		return &mockReadCloser{failClose: true}, nil
	}
	uncomp, err := zc.uncompressDI(comp, newMockReadCloser)
	if uncomp != nil {
		t.Error("0xEF3A01")
	}
	if !matchError(err, "failed mockReadCloser.Close") {
		t.Error("0xEA0F76", "wrong error:", err)
	}
}

// -----------------------------------------------------------------------------

// mockReadCloser is a mock io.ReadCloser with methods you can make fail.
type mockReadCloser struct {
	failRead  bool
	failClose bool
}

// Read is a method of mockReadCloser implementing io.ReadCloser.
//
// You can make it return an error by setting mockReadCloser.failRead.
//
func (mk *mockReadCloser) Read(p []byte) (n int, err error) {
	if mk.failRead {
		return 0, makeError(0xEF8E54, "failed mockReadCloser.Read")
	}
	return len(p), nil
}

// Close is a method of mockReadCloser implementing io.ReadCloser.
//
// You can make it return an error by setting mockReadCloser.failClose.
//
func (mk *mockReadCloser) Close() error {
	if mk.failClose {
		return makeError(0xE1FB2C, "failed mockReadCloser.Close")
	}
	return nil
}

// -----------------------------------------------------------------------------

// zCompress compresses the string from zInput(): must always succeed
func zCompress(t *testing.T) []byte {
	zc := zlibCompressor{}
	input := zInput()
	comp, err := zc.Compress(input)
	if err != nil {
		t.Error("0xE26CD5", err)
	}
	return comp
}

// zInput: provides a string to compress
func zInput() []byte {
	return []byte(
		strings.Repeat("The quick brown fox jumps over the lazy dog!", 7),
	)
}

// end

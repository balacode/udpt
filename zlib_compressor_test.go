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

// go test -run Test_ZLibCompressor_
//
func Test_ZLibCompressor_(t *testing.T) {
	input := []byte(
		strings.Repeat("The quick brown fox jumps over the lazy dog!", 7),
	)
	// compresses input: must always succeed
	compress := func() []byte {
		zc := zlibCompressor{}
		comp, err := zc.Compress(input)
		if err != nil {
			t.Error("0xE26CD5", err)
		}
		return comp
	}
	{
		// tests compressing and uncompressing data
		comp := compress()
		zc := zlibCompressor{}
		uncomp, err := zc.Uncompress(comp) // must succeed
		if err != nil {
			t.Error("0xE38FD2", err)
		}
		if !bytes.Equal(input, uncomp) {
			t.Error("0xEB0E80", "data corrupted")
		}
	}
	{
		// uncompressing less than 4 bytes must fail
		zc := zlibCompressor{}
		comp, err := zc.Uncompress([]byte{1, 2, 3})
		if comp != nil {
			t.Error("0xEE55FC")
		}
		if !matchError(err, "invalid 'compressed'") {
			t.Error("0xE4AB37")
		}
	}
	{
		// test wr.Write() failing in compressDI()
		zc := zlibCompressor{}
		cbuf := bytes.Buffer{}
		wrc := &mockWriteCloser{failWrite: true}
		comp, err := zc.compressDI(input, wrc, &cbuf)
		if comp != nil {
			t.Error("0xE27BB2")
		}
		if !matchError(err, "failed mockWriteCloser.Write") {
			t.Error("0xEE9C54", err)
		}
	}
	{
		// test wr.Close() failing in compressDI()
		zc := zlibCompressor{}
		cbuf := bytes.Buffer{}
		wrc := &mockWriteCloser{failClose: true}
		comp, err := zc.compressDI(input, wrc, &cbuf)
		if comp != nil {
			t.Error("0xE6F6A5")
		}
		if !matchError(err, "failed mockWriteCloser.Close") {
			t.Error("0xE4DF92", err)
		}
	}
	{
		// test reader.Read() failing so io.Copy() fails too in uncompressDI()
		comp, zc := compress(), zlibCompressor{}
		newMockReadCloser := func(io.Reader) (io.ReadCloser, error) {
			return &mockReadCloser{failRead: true}, nil
		}
		uncomp, err := zc.uncompressDI(comp, newMockReadCloser)
		if uncomp != nil {
			t.Error("0xE3DA4F")
		}
		if !matchError(err, "failed mockReadCloser.Read") {
			t.Error("0xE81C62")
		}
	}
	{
		// test reader.Close() failing in uncompressDI()
		comp, zc := compress(), zlibCompressor{}
		newMockReadCloser := func(io.Reader) (io.ReadCloser, error) {
			return &mockReadCloser{failClose: true}, nil
		}
		uncomp, err := zc.uncompressDI(comp, newMockReadCloser)
		if uncomp != nil {
			t.Error("0xEF3A01")
		}
		if !matchError(err, "failed mockReadCloser.Close") {
			t.Error("0xEA0F76")
		}
	}
} //                                                        Test_ZLibCompressor_

// -----------------------------------------------------------------------------

// mockReadCloser is a mock io.ReadCloser with methods you can make fail.
type mockReadCloser struct {
	failRead  bool
	failClose bool
} //                                                              mockReadCloser

// Read is a method of mockReadCloser implementing io.ReadCloser.
//
// You can make it return an error by setting mockReadCloser.failRead.
//
func (ob *mockReadCloser) Read(p []byte) (n int, err error) {
	if ob.failRead {
		return 0, makeError(0xEF8E54, "failed mockReadCloser.Read")
	}
	return len(p), nil
} //                                                                        Read

// Close is a method of mockReadCloser implementing io.ReadCloser.
//
// You can make it return an error by setting mockReadCloser.failClose.
//
func (ob *mockReadCloser) Close() error {
	if ob.failClose {
		return makeError(0xE1FB2C, "failed mockReadCloser.Close")
	}
	return nil
} //                                                                       Close

// end

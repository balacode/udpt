// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                    /[compression.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

// Compression implements functions to compress and uncompress byte slices.
type Compression interface {

	// Compress compresses 'uncompressed' bytes and returns the compressed
	// bytes. If there was an error, returns nil and the error value.
	Compress(uncompressed []byte) ([]byte, error)

	// Uncompress uncompresses 'compressed' bytes and returns the uncompressed
	// bytes. If there was an error, returns nil and the error value.
	Uncompress(compressed []byte) ([]byte, error)
} //                                                                 Compression

// end

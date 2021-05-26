// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                    /[compression.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

// Compression implements functions to compress and uncompress byte slices.
type Compression interface {

	// Compress compresses 'data' and returns the compressed bytes.
	// If there was an error, returns nil and the error instance.
	Compress(data []byte) ([]byte, error)

	// Uncompress uncompresses bytes and returns the uncompressed bytes.
	// If there was an error, returns nil and the error instance.
	Uncompress(comp []byte) ([]byte, error)
} //                                                                 Compression

// end

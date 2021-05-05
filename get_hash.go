// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                       /[get_hash.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"crypto/sha256"
)

// getHash returns the SHA-256 hash of data as a slice of 32 bytes.
func getHash(data []byte) ([]byte, error) {
	sha := sha256.New()
	n, err := sha.Write(data)
	if n != len(data) || err != nil {
		return nil, makeError(0xE3A9B8, err)
	}
	ret := sha.Sum(nil)
	if len(ret) != 32 {
		return nil, makeError(0xE4D3E1, "invalid hash size")
	}
	return ret, nil
} //                                                                     getHash

// end

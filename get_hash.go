// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                       /[get_hash.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"crypto/sha256"
)

// getHash returns the SHA-256 hash of data as a slice of 32 bytes.
func getHash(data []byte) []byte {
	sha := sha256.New()
	n, err := sha.Write(data)
	if n != len(data) || err != nil {
		_ = logError(0xE3A9B8, "(Write):", err)
		logInfo("n:", n, "len(data):", len(data))
		return nil
	}
	ret := sha.Sum(nil)
	if len(ret) != 32 {
		_ = logError(0xE4D3E1, "(sha.Sum)")
		return nil
	}
	return ret
} //                                                                     getHash

// end

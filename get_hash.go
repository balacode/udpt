// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                       /[get_hash.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"crypto/sha256"
	"hash"
)

// getHash returns the SHA-256 hash of data as a slice of 32 bytes.
func getHash(data []byte) []byte {
	return getHashDI(data, sha256.New())
} //                                                                     getHash

// getHashDI is only used by getHash() and provides parameters
// for dependency injection, to enable mocking during testing.
func getHashDI(data []byte, hs hash.Hash) []byte {
	n, err := hs.Write(data)
	if n != len(data) || err != nil {
		// this should never happen (see hash.Hash.Write in Go docs)
		panic(makeError(0xE51EC0, err).Error())
	}
	ret := hs.Sum(nil)
	if len(ret) != 32 {
		// this should never happen
		panic(makeError(0xE4D3E1, "invalid hash size").Error())
	}
	return ret
} //                                                                   getHashDI

// end

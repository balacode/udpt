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
	return getHashD(data, sha256.New())
} //                                                                     getHash

// getHashD carries out the work for getHash(), and accepts
// a dependency 'hs' that can be mocked during testing.
func getHashD(data []byte, hs hash.Hash) []byte {
	n, err := hs.Write(data)
	if n != len(data) || err != nil {
		// this should never happen
		panic(makeError(0xE51EC0, err).Error())
	}
	ret := hs.Sum(nil)
	if len(ret) != 32 {
		// this should never happen
		panic(makeError(0xE4D3E1, "invalid hash size").Error())
	}
	return ret
} //                                                                    getHashD

// end

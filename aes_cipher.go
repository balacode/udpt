// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                     /[aes_cipher.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

const errKeySize = "AES-256 key must be 32 bytes long"

// AESCipher implements the SymmetricCipher interface that encrypts/decrypts
// plaintext using the AES-256 symmetric cipher algorithm.
type AESCipher struct {
	aesKey []byte
	gcm    cipher.AEAD
} //                                                                   AESCipher

// ValidateKey checks if 'key' is acceptable for use with the cipher.
// For example it must be of the right size.
//
// For AES-256, the cipher must be exactly 32 bytes long.
//
func (ob *AESCipher) ValidateKey(key []byte) error {
	if len(key) != 32 {
		return logError(0xE42FDB,
			"AES-256 key must be 32, but it is", len(key), "bytes long")
	}
	return nil
} //                                                                 ValidateKey

// InitCipher initializes a cipher with the specified secret key.
func (ob *AESCipher) InitCipher(key []byte) error {
	if len(key) != 32 {
		return logError(0xE32BD3, errKeySize)
	}
	cphr, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	ob.gcm, err = cipher.NewGCM(cphr)
	if err != nil {
		return err
	}
	ob.aesKey = key
	return nil
} //                                                                  InitCipher

// Encrypt encrypts plaintext using the key given to InitCipher and
// returns the encrypted ciphertext, using AES-256 symmetric cipher.
func (ob *AESCipher) Encrypt(plaintext []byte) (ciphertext []byte, err error) {
	if len(ob.aesKey) != 32 {
		return nil, logError(0xE64A2E, errKeySize)
	}
	// nonce is a byte array filled with cryptographically secure random bytes
	n := ob.gcm.NonceSize()
	nonce := make([]byte, n)
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}
	ciphertext = ob.gcm.Seal(
		nonce,     // dst []byte,
		nonce,     // nonce []byte,
		plaintext, // plaintext []byte,
		nil,       // additionalData []byte) []byte
	)
	return ciphertext, nil
} //                                                                     Encrypt

// Decrypt decrypts ciphertext using the key given to InitCipher and
// returns the decrypted plaintext, using AES-256 symmetric cipher.
func (ob *AESCipher) Decrypt(ciphertext []byte) (plaintext []byte, err error) {
	if len(ob.aesKey) != 32 {
		return nil, logError(0xE35A87, errKeySize)
	}
	n := ob.gcm.NonceSize()
	if len(ciphertext) < n {
		return nil, err
	}
	nonce := ciphertext[:n]
	ciphertext = ciphertext[n:]
	plaintext, err = ob.gcm.Open(
		nil,        // dst []byte
		nonce,      // nonce []byte
		ciphertext, // ciphertext []byte
		nil,        // additionalData []byte
	)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
} //                                                                     Decrypt

// end

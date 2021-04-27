// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                    /[aes_encrypt.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

// aesEncrypt encrypts plaintext using secretKey and returns
// the encrypted cipherthext, using AES-256 symmetric cipher.
func aesEncrypt(plaintext, secretKey []byte) (ciphertext []byte, err error) {
	if len(secretKey) != 32 {
		return nil, logError(0xE9A91B, "AES secretKey must be 32 bytes long")
	}
	var gcm cipher.AEAD
	{
		cip, err := aes.NewCipher(secretKey)
		if err != nil {
			return nil, err
		}
		gcm, err = cipher.NewGCM(cip)
		if err != nil {
			return nil, err
		}
	}
	// nonce is byte array filled with cryptographically secure random bytes
	var nonce []byte
	{
		n := gcm.NonceSize()
		nonce = make([]byte, n)
		_, err := io.ReadFull(rand.Reader, nonce)
		if err != nil {
			return nil, err
		}
	}
	// encrypt
	ciphertext = gcm.Seal(
		nonce,     // dst []byte,
		nonce,     // nonce []byte,
		plaintext, // plaintext []byte,
		nil,       // additionalData []byte) []byte
	)
	return ciphertext, nil
} //                                                                  aesEncrypt

// end

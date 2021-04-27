// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                    /[aes_decrypt.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"crypto/aes"
	"crypto/cipher"
)

// aesDecrypt decrypts cipherthext using secretKey and returns
// the decrypted plaintext, using AES-256 symmetric cipher.
func aesDecrypt(ciphertext, secretKey []byte) (plaintext []byte, err error) {
	if len(secretKey) != 32 {
		return nil, logError(0xE4A45F, "AES secretKey must be 32 bytes long")
	}
	chp, err := aes.NewCipher(secretKey)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(chp)
	if err != nil {
		return nil, err
	}
	n := gcm.NonceSize()
	if len(ciphertext) < n {
		return nil, err
	}
	nonce := ciphertext[:n]
	ciphertext = ciphertext[n:]
	plaintext, err = gcm.Open(
		nil,        // dst []byte
		nonce,      // nonce []byte
		ciphertext, // ciphertext []byte
		nil,        // additionalData []byte
	)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
} //                                                                  aesDecrypt

// end

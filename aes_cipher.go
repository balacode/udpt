// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                     /[aes_cipher.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

// aesCipher implements the SymmetricCipher interface that encrypts and
// decrypts plaintext using the AES-256 symmetric cipher algorithm.
type aesCipher struct {
	cryptoKey []byte
	gcm       cipher.AEAD
} //                                                                   aesCipher

// ValidateKey checks if an encryption key is suitable for use with the cipher.
// For example it must be of the right size.
//
// For AES-256, the encryption key must be exactly 32 bytes long.
//
func (ac *aesCipher) ValidateKey(cryptoKey []byte) error {
	if len(cryptoKey) != 32 {
		return makeError(0xE42FDB, "AES-256 key must be 32 bytes long")
	}
	return nil
} //                                                                 ValidateKey

// SetKey initializes the cipher with the specified encryption key.
//
// If the cipher is already initialized with the given key, does nothing.
// The same key is used for encryption and decryption.
//
func (ac *aesCipher) SetKey(cryptoKey []byte) error {
	return ac.setKeyDI(cryptoKey, aes.NewCipher, cipher.NewGCM)
} //                                                                      SetKey

// setKeyDI is only used by SetKey() and provides parameters
// for dependency injection, to enable mocking during testing.
func (ac *aesCipher) setKeyDI(
	cryptoKey []byte,
	aesNewCipher func([]byte) (cipher.Block, error),
	cipherNewGCM func(cipher.Block) (cipher.AEAD, error),
) error {
	err := ac.ValidateKey(cryptoKey)
	if err != nil {
		return makeError(0xE32BD3, err)
	}
	if bytes.Equal(ac.cryptoKey, cryptoKey) {
		return nil
	}
	cphr, err := aesNewCipher(cryptoKey)
	if err != nil {
		return err
	}
	gcm, err := cipherNewGCM(cphr)
	if err != nil {
		return err
	}
	ac.gcm = gcm
	ac.cryptoKey = cryptoKey
	return nil
} //                                                                    setKeyDI

// Encrypt encrypts plaintext using the encryption key given to SetKey
// and returns the encrypted ciphertext, using AES-256 symmetric cipher.
//
// You need to call SetKey at least once before you call Encrypt.
//
func (ac *aesCipher) Encrypt(plaintext []byte) (ciphertext []byte, err error) {
	return ac.encryptDI(plaintext, io.ReadFull)
} //                                                                     Encrypt

// encryptDI is only used by Encrypt() and provides parameters
// for dependency injection, to enable mocking during testing.
func (ac *aesCipher) encryptDI(
	plaintext []byte,
	ioReadFull func(io.Reader, []byte) (int, error),
) (ciphertext []byte, err error) {
	//
	err = ac.ValidateKey(ac.cryptoKey)
	if err != nil {
		return nil, makeError(0xE64A2E, err)
	}
	// nonce is a byte array filled with cryptographically secure random bytes
	n := ac.gcm.NonceSize() // = gcmStandardNonceSize = 12 bytes
	nonce := make([]byte, n)
	_, err = ioReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}
	ciphertext = ac.gcm.Seal(
		nonce,     // dst
		nonce,     // nonce
		plaintext, // plaintext
		nil,       // additionalData
	)
	return ciphertext, nil
} //                                                                   encryptDI

// Decrypt decrypts ciphertext using the encryption key given to SetKey
// and returns the decrypted plaintext, using AES-256 symmetric cipher.
//
// You need to call SetKey at least once before you call Decrypt.
//
func (ac *aesCipher) Decrypt(ciphertext []byte) (plaintext []byte, err error) {
	err = ac.ValidateKey(ac.cryptoKey)
	if err != nil {
		return nil, makeError(0xE35A87, err)
	}
	n := ac.gcm.NonceSize()
	if len(ciphertext) < n {
		return nil, makeError(0xE5F7E2, "invalid ciphertext")
	}
	nonce := ciphertext[:n]
	ciphertext = ciphertext[n:]
	plaintext, err = ac.gcm.Open(
		nil,        // dst
		nonce,      // nonce
		ciphertext, // ciphertext
		nil,        // additionalData
	)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
} //                                                                     Decrypt

// end

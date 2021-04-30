// -----------------------------------------------------------------------------
// github.com/balacode/udpt                               /[symmetric_cipher.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

// SymmetricCipher interface provides methods to initialize a symmetric-key
// cipher and use it to encrypt and decrypt plaintext.
// A symmetric-key cipher uses the same key to encrypt and decrypt.
type SymmetricCipher interface {

	// ValidateKey checks if 'key' is acceptable for use with the cipher.
	// For example it must be of the right size.
	ValidateKey(key []byte) error

	// InitCipher initializes a cipher with the specified secret key.
	InitCipher(key []byte) error

	// Encrypt encrypts plaintext using the key given to InitCipher and
	// returns the encrypted ciphertext, using a symmetric-key cipher.
	Encrypt(plaintext []byte) (ciphertext []byte, err error)

	// Decrypt decrypts ciphertext using the key given to InitCipher and
	// returns the decrypted plaintext, using a symmetric-key cipher.
	Decrypt(ciphertext []byte) (plaintext []byte, err error)
} //                                                             SymmetricCipher

// end

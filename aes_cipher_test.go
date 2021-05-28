// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                /[aes_cipher_test.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"crypto/aes"
	"crypto/cipher"
	"io"
	"testing"
)

const testAESKey = "A0CFDD4FA7B545088826A73C9A93AB8A"

// (ac *aesCipher) ValidateKey(cryptoKey []byte) error
//
// go test -run Test_aesCipher_ValidateKey_
//
func Test_aesCipher_ValidateKey_(t *testing.T) {
	ac := &aesCipher{}
	err := ac.ValidateKey([]byte("12345678901234567890123456789012")) // 32b
	if err != nil {
		t.Error("0xEE98F7", err)
	}
	err = ac.ValidateKey(nil)
	if !matchError(err, "AES-256 key must be 32 bytes long") {
		t.Error("0xEC1B42")
	}
	err = ac.ValidateKey([]byte("1234567890123456789012345678901")) // 31b
	if !matchError(err, "AES-256 key must be 32 bytes long") {
		t.Error("0xE1D2B9")
	}
	err = ac.ValidateKey([]byte("123456789012345678901234567890123")) // 33b
	if !matchError(err, "AES-256 key must be 32 bytes long") {
		t.Error("0xE85BC4")
	}
}

// (ac *aesCipher) SetKey(cryptoKey []byte) error
//
// go test -run Test_aesCipher_SetKey_
//
func Test_aesCipher_SetKey_(t *testing.T) {
	//
	// must succeed; encryption key is 32 bytes long:
	var cphr aesCipher
	const goodKey = "BE30FB257682466ABA9071755E780344"
	err := cphr.SetKey([]byte(goodKey))
	if err != nil {
		t.Error("0xE6AD9A", err)
	}
	// must succeed, but since the key is the same, don't create cipher again
	createdCipher := false
	aesNewCipher := func(cryptoKey []byte) (cipher.Block, error) {
		createdCipher = true
		return aes.NewCipher(cryptoKey)
	}
	cipherNewGCM := func(c cipher.Block) (cipher.AEAD, error) {
		createdCipher = true
		return cipher.NewGCM(c)
	}
	err = cphr.setKeyDI([]byte(goodKey), aesNewCipher, cipherNewGCM)
	if err != nil {
		t.Error("0xE5DB3F", err)
	}
	if createdCipher {
		t.Error("0xEF12FB")
	}
	// must fail; encryption key is not 32 bytes long:
	err = cphr.SetKey(nil)
	if !matchError(err, "AES-256 key must be 32 bytes long") {
		t.Error("0xE9DF0F", "wrong error:", err)
	}
	// 33 bytes: too long
	err = cphr.SetKey([]byte("123456789012345678901234567890123"))
	if !matchError(err, "AES-256 key must be 32 bytes long") {
		t.Error("0xEC7FF8", "wrong error:", err)
	}
	// 31 bytes: too short
	err = cphr.SetKey([]byte("1234567890123456789012345678901"))
	if !matchError(err, "AES-256 key must be 32 bytes long") {
		t.Error("0xE2C07B", "wrong error:", err)
	}
	// must fail when aes.NewCipher() errors
	cphr.cryptoKey = nil
	aesNewCipher = func([]byte) (cipher.Block, error) {
		return nil, makeError(0xEC56A6, "failed aesNewCipher")
	}
	err = cphr.setKeyDI([]byte(goodKey), aesNewCipher, cipher.NewGCM)
	if !matchError(err, "failed aesNewCipher") {
		t.Error("0xE7F54F", "wrong error:", err)
	}
	// must fail when cipher.NewGCM() errors
	cphr.cryptoKey = nil
	cipherNewGCM = func(cipher.Block) (cipher.AEAD, error) {
		return nil, makeError(0xEC4CF2, "failed cipherNewGCM")
	}
	err = cphr.setKeyDI([]byte(goodKey), aes.NewCipher, cipherNewGCM)
	if !matchError(err, "failed cipherNewGCM") {
		t.Error("0xE4C78A", "wrong error:", err)
	}
}

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
// (ac *aesCipher) Encrypt(plaintext []byte) (ciphertext []byte, err error)
//
// go test -run Test_aesCipher_Encrypt_*

// encryption must succeed:
func Test_aesCipher_Encrypt_1(t *testing.T) {
	cphr := newTestAESCipher(t)
	ciphertext, err := cphr.Encrypt([]byte("abc"))
	if err != nil {
		t.Error("0xE5B98A", err)
	}
	// decrypting ciphertext must get back plaintext
	n := cphr.gcm.NonceSize()
	nonce := ciphertext[:n]
	ciphertext = ciphertext[n:]
	plaintext, err := cphr.gcm.Open(nil, nonce, ciphertext, nil)
	if string(plaintext) != "abc" {
		t.Error("0xE48E19")
	}
	if err != nil {
		t.Error("0xE9A2DE", err)
	}
}

// must fail encrypting because there is no encryption key specified:
func Test_aesCipher_Encrypt_2(t *testing.T) {
	var cphr aesCipher
	ciphertext, err := cphr.Encrypt([]byte("abc"))
	if ciphertext != nil {
		t.Error("0xEC8FF6")
	}
	if !matchError(err, "AES-256 key must be 32 bytes long") {
		t.Error("0xEF26E8", "wrong error:", err)
	}
}

// must fail encrypting when io.ReadFull() fails:
func Test_aesCipher_Encrypt_3(t *testing.T) {
	cphr := newTestAESCipher(t)
	ioReadFull := func(io.Reader, []byte) (int, error) {
		return 0, makeError(0xED5D20, "failed ioReadFull")
	}
	ciphertext, err := cphr.encryptDI([]byte("abc"), ioReadFull)
	if ciphertext != nil {
		t.Error("0xE8D36B")
	}
	if !matchError(err, "failed ioReadFull") {
		t.Error("0xE9D47D", "wrong error:", err)
	}
}

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
// (ac *aesCipher) Decrypt(ciphertext []byte) (plaintext []byte, err error)
//
// go test -run Test_aesCipher_Decrypt_*

// decryption must succeed:
func Test_aesCipher_Decrypt_1(t *testing.T) {
	var (
		cphr           = newTestAESCipher(t)
		ciphertext     = newTestAESCiphertext()
		plaintext, err = cphr.Decrypt(ciphertext)
	)
	if string(plaintext) != "abc" {
		t.Error("0xEB8E21")
	}
	if err != nil {
		t.Error("0xE64B43", err)
	}
}

// must fail decrypting because there is no encryption key specified:
func Test_aesCipher_Decrypt_2(t *testing.T) {
	var cphr aesCipher
	plaintext, err := cphr.Decrypt(nil)
	if plaintext != nil {
		t.Error("0xED83C7")
	}
	if !matchError(err, "AES-256 key must be 32 bytes long") {
		t.Error("0xE43B12", "wrong error:", err)
	}
}

// must fail decrypting because ciphertext is too short:
func Test_aesCipher_Decrypt_3(t *testing.T) {
	cphr := newTestAESCipher(t)
	plaintext, err := cphr.Decrypt([]byte("12345678901"))
	if plaintext != nil {
		t.Error("0xE5BE69")
	}
	if !matchError(err, "invalid ciphertext") {
		t.Error("0xEC9F1F", "wrong error:", err)
	}
}

// must fail decrypting because ciphertext is garbage:
func Test_aesCipher_Decrypt_4(t *testing.T) {
	cphr := newTestAESCipher(t)
	plaintext, err := cphr.Decrypt([]byte("5DBB4C78125C442591A9293C9D5A5CE6"))
	if plaintext != nil {
		t.Error("0xE6CA73")
	}
	if !matchError(err, "cipher: message authentication failed") {
		t.Error("0xEE3FD8", "wrong error:", err)
	}
}

// must fail decrypting because ciphertext is tampered:
func Test_aesCipher_Decrypt_5(t *testing.T) {
	cphr := newTestAESCipher(t)
	ciphertext := newTestAESCiphertext()
	ciphertext[0]++ // <- tamper
	plaintext, err := cphr.Decrypt(ciphertext)
	if string(plaintext) != "" {
		t.Error("0xE2BE96")
	}
	if !matchError(err, "cipher: message authentication failed") {
		t.Error("0xEC44BA", "wrong error:", err)
	}
}

// -----------------------------------------------------------------------------

// newTestAESCipher creates an AES cipher for testing (uses testAESKey)
func newTestAESCipher(t *testing.T) *aesCipher {
	cphr, err := aes.NewCipher([]byte(testAESKey))
	if err != nil {
		t.Error("0xEF9AC2", err)
	}
	gcm, err := cipher.NewGCM(cphr)
	if err != nil {
		t.Error("0xE15BF6", err)
	}
	return &aesCipher{gcm: gcm, cryptoKey: []byte(testAESKey)}
}

// newTestAESCiphertext creates ciphertext that can
// be decrypted into "abc" with testAESKey.
func newTestAESCiphertext() []byte {
	return []byte{
		0x85, 0x8D, 0xAB, 0x9E, 0x4A, 0x89, 0x8D, 0x3B,
		0x46, 0x7C, 0xD9, 0x40, 0xEA, 0x37, 0xD5, 0x08,
		0x07, 0x0A, 0xC1, 0xBC, 0x0F, 0xAC, 0xF0, 0xC3,
		0x10, 0xF3, 0x09, 0x53, 0x3D, 0xBB, 0xD8,
	}
	// you can generate ciphertext with:
	//
	// cphr := newTestAESCipher()
	// ctext, _ := ac.Encrypt([]byte("abc"))
	// fmt.Printf("%#v", ctext)
	//
	// note that each time the ciphertext will differ
}

// end

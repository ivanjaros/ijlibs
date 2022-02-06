package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"github.com/ivanjaros/ijlibs/random"
)

type Encryptor struct {
	block cipher.AEAD
}

// Key must be exactly 12, 24 or 32 bytes long.
func New(key []byte) (*Encryptor, error) {
	keyLen := len(key)
	if key == nil || keyLen == 0 {
		return nil, errors.New("no encryption key provided")
	}

	if keyLen != 16 && keyLen != 24 && keyLen != 32 {
		return nil, errors.New("encryption key has to be 16, 24 or 32 bytes long to select AES-128, AES-192, or AES-256")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return &Encryptor{block: aesgcm}, nil
}

// Creates new encryptor instance, or panics upon error.
// Key must be exactly 12, 24 or 32 bytes long.
func MustNew(key []byte) *Encryptor {
	e, err := New(key)
	if err != nil {
		panic(err)
	}
	return e
}

// nonce must be exactly 12 bytes long
func (e *Encryptor) Encrypt(data, nonce []byte) ([]byte, error) {
	return e.block.Seal(nil, nonce, data, nil), nil
}

// nonce must be exactly 12 bytes long
func (e *Encryptor) Decrypt(data, nonce []byte) ([]byte, error) {
	return e.block.Open(nil, nonce, data, nil)
}

// encrypts data with random nonce
func (e *Encryptor) EncryptRandom(data []byte) (result []byte, nonce []byte, err error) {
	nonce = random.CryptoBytes(e.block.NonceSize())
	result, err = e.Encrypt(data, nonce)
	return
}

// Package encrypt provides value-level encryption and decryption for .env entries.
// It supports AES-GCM symmetric encryption using a passphrase-derived key,
// allowing sensitive values to be stored safely in version control.
package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

const (
	encryptedPrefix = "enc:"
	pbkdf2Iter      = 100_000
	keyLen          = 32
	saltLen         = 16
)

// ErrNotEncrypted is returned when decrypting a value that lacks the encrypted prefix.
var ErrNotEncrypted = errors.New("value is not encrypted")

// Encrypter encrypts and decrypts values using a passphrase.
type Encrypter struct {
	passphrase string
}

// New returns a new Encrypter using the provided passphrase.
func New(passphrase string) *Encrypter {
	return &Encrypter{passphrase: passphrase}
}

// IsEncrypted reports whether the value carries the encrypted prefix.
func IsEncrypted(value string) bool {
	return strings.HasPrefix(value, encryptedPrefix)
}

// Encrypt encrypts plaintext and returns a prefixed base64 blob.
func (e *Encrypter) Encrypt(plaintext string) (string, error) {
	salt := make([]byte, saltLen)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", fmt.Errorf("encrypt: generate salt: %w", err)
	}

	key := e.deriveKey(salt)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("encrypt: new cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("encrypt: new gcm: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("encrypt: generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	blob := append(salt, ciphertext...)
	return encryptedPrefix + base64.StdEncoding.EncodeToString(blob), nil
}

// Decrypt decrypts a value produced by Encrypt.
func (e *Encrypter) Decrypt(value string) (string, error) {
	if !IsEncrypted(value) {
		return "", ErrNotEncrypted
	}

	blob, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(value, encryptedPrefix))
	if err != nil {
		return "", fmt.Errorf("encrypt: decode base64: %w", err)
	}
	if len(blob) < saltLen {
		return "", fmt.Errorf("encrypt: blob too short")
	}

	salt, data := blob[:saltLen], blob[saltLen:]
	key := e.deriveKey(salt)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("encrypt: new cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("encrypt: new gcm: %w", err)
	}

	ns := gcm.NonceSize()
	if len(data) < ns {
		return "", fmt.Errorf("encrypt: ciphertext too short")
	}

	plaintext, err := gcm.Open(nil, data[:ns], data[ns:], nil)
	if err != nil {
		return "", fmt.Errorf("encrypt: decrypt: %w", err)
	}
	return string(plaintext), nil
}

func (e *Encrypter) deriveKey(salt []byte) []byte {
	return pbkdf2.Key([]byte(e.passphrase), salt, pbkdf2Iter, keyLen, sha256.New)
}

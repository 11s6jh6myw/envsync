package encrypt_test

import (
	"strings"
	"testing"

	"github.com/your-org/envsync/internal/encrypt"
)

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	e := encrypt.New("supersecret")
	plaintext := "my-database-password"

	cipher, err := e.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt: unexpected error: %v", err)
	}

	if !encrypt.IsEncrypted(cipher) {
		t.Errorf("expected encrypted prefix, got %q", cipher)
	}

	got, err := e.Decrypt(cipher)
	if err != nil {
		t.Fatalf("Decrypt: unexpected error: %v", err)
	}
	if got != plaintext {
		t.Errorf("got %q, want %q", got, plaintext)
	}
}

func TestEncrypt_ProducesUniqueOutputs(t *testing.T) {
	e := encrypt.New("passphrase")
	a, _ := e.Encrypt("value")
	b, _ := e.Encrypt("value")
	if a == b {
		t.Error("expected different ciphertexts due to random nonce/salt")
	}
}

func TestDecrypt_WrongPassphrase(t *testing.T) {
	e1 := encrypt.New("correct-passphrase")
	e2 := encrypt.New("wrong-passphrase")

	cipher, err := e1.Encrypt("secret")
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}

	_, err = e2.Decrypt(cipher)
	if err == nil {
		t.Error("expected error when decrypting with wrong passphrase")
	}
}

func TestDecrypt_NotEncrypted(t *testing.T) {
	e := encrypt.New("passphrase")
	_, err := e.Decrypt("plaintext-value")
	if err != encrypt.ErrNotEncrypted {
		t.Errorf("expected ErrNotEncrypted, got %v", err)
	}
}

func TestIsEncrypted(t *testing.T) {
	tests := []struct {
		value string
		want  bool
	}{
		{"enc:abc123==", true},
		{"plaintext", false},
		{"", false},
		{"enc:", true},
		{"ENC:uppercase", false},
	}
	for _, tt := range tests {
		got := encrypt.IsEncrypted(tt.value)
		if got != tt.want {
			t.Errorf("IsEncrypted(%q) = %v, want %v", tt.value, got, tt.want)
		}
	}
}

func TestEncrypt_PrefixPresent(t *testing.T) {
	e := encrypt.New("key")
	out, err := e.Encrypt("hello")
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	if !strings.HasPrefix(out, "enc:") {
		t.Errorf("expected enc: prefix, got %q", out)
	}
}

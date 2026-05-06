// Package encrypt provides AES-GCM encryption and decryption for individual
// .env values, enabling secrets to be stored safely in version-controlled
// environment files.
//
// # Overview
//
// Values encrypted by this package are stored with a distinguishing "enc:"
// prefix followed by a base64-encoded blob that contains a random salt,
// a random nonce, and the AES-GCM ciphertext. The encryption key is derived
// from a caller-supplied passphrase using PBKDF2-SHA256.
//
// # Usage
//
//	e := encrypt.New(os.Getenv("ENVSYNC_PASSPHRASE"))
//
//	// Encrypt a sensitive value before writing to disk.
//	encVal, err := e.Encrypt("s3cr3t")
//
//	// Decrypt when the value is needed at runtime.
//	plain, err := e.Decrypt(encVal)
//
// # Integration
//
// Use encrypt.IsEncrypted to detect whether a parsed entry's value should be
// decrypted before use, making it straightforward to layer encryption on top
// of the existing parser and sync pipeline.
package encrypt

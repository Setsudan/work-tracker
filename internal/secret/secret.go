package secret

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
)

// Secret provides authenticated encryption and keyed hashing utilities.
// It uses AES-256-GCM for encryption and HMAC-SHA256 for keyed hashing.
// The same 32-byte key is used for both unless you decide to split keys externally.
// Nonces are 12 bytes and are prepended to the ciphertext, then base64-encoded.
// All ciphertext representations are base64 strings for portability across Redis types.

type Secret struct {
	aead   cipher.AEAD
	macKey []byte
}

// New creates a new Secret from a base64-encoded 32-byte key.
func New(base64Key string) (*Secret, error) {
	key, err := base64.StdEncoding.DecodeString(base64Key)
	if err != nil {
		return nil, err
	}
	if len(key) != 32 {
		return nil, errors.New("DATA_ENCRYPTION_KEY must be base64-encoded 32 bytes (AES-256)")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return &Secret{aead: aead, macKey: key}, nil
}

// Encrypt returns base64 of nonce||ciphertext.
func (s *Secret) Encrypt(plaintext []byte) (string, error) {
	nonce := make([]byte, 12)
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	ciphertext := s.aead.Seal(nil, nonce, plaintext, nil)
	buf := append(nonce, ciphertext...)
	return base64.StdEncoding.EncodeToString(buf), nil
}

// DecryptString takes base64 of nonce||ciphertext and returns plaintext.
func (s *Secret) DecryptString(b64 string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, err
	}
	if len(data) < 12 {
		return nil, errors.New("ciphertext too short")
	}
	nonce := data[:12]
	ct := data[12:]
	pt, err := s.aead.Open(nil, nonce, ct, nil)
	if err != nil {
		return nil, err
	}
	return pt, nil
}

// HMACString returns a hex-encoded HMAC-SHA256 of the input using the secret key.
func (s *Secret) HMACString(input string) string {
	h := hmac.New(sha256.New, s.macKey)
	h.Write([]byte(input))
	return hex.EncodeToString(h.Sum(nil))
}

// MaybeDecryptBytes attempts to treat the input as a base64-encoded ciphertext; if that fails,
// it returns the original bytes and false. On success, returns plaintext and true.
func (s *Secret) MaybeDecryptBytes(in []byte) ([]byte, bool) {
	// quick reject: base64 encoding only uses a limited charset; however, we attempt decode regardless
	decoded, err := base64.StdEncoding.DecodeString(string(in))
	if err != nil || len(decoded) < 12 {
		return in, false
	}
	nonce := decoded[:12]
	ct := decoded[12:]
	pt, err := s.aead.Open(nil, nonce, ct, nil)
	if err != nil {
		return in, false
	}
	return pt, true
}
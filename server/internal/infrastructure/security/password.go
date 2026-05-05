package security

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/argon2"
)

const (
	argon2idVersion = argon2.Version
	saltLength      = 16
	keyLength       = 32
)

type Argon2idConfig struct {
	MemoryKiB   int
	Iterations  int
	Parallelism int
}

type Argon2idPasswordHasher struct {
	cfg Argon2idConfig
}

func NewArgon2idPasswordHasher(cfg Argon2idConfig) *Argon2idPasswordHasher {
	return &Argon2idPasswordHasher{cfg: cfg}
}

func (h *Argon2idPasswordHasher) Hash(password string) (string, error) {
	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		uint32(h.cfg.Iterations),
		uint32(h.cfg.MemoryKiB),
		uint8(h.cfg.Parallelism),
		keyLength,
	)

	encodedSalt := base64.RawStdEncoding.EncodeToString(salt)
	encodedHash := base64.RawStdEncoding.EncodeToString(hash)
	return fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2idVersion,
		h.cfg.MemoryKiB,
		h.cfg.Iterations,
		h.cfg.Parallelism,
		encodedSalt,
		encodedHash,
	), nil
}

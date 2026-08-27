package password

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	memory      = 64 * 1024
	iterations  = 3
	parallelism = 4
	saltLength  = 16
	keyLength   = 32
)

type Argon2Hasher struct{}

func NewArgon2Hasher() *Argon2Hasher {
	return &Argon2Hasher{}
}

func (h *Argon2Hasher) Hash(password string) (string, error) {
	if password == "" {
		return "", errors.New("password cannot be empty")
	}

	salt := make([]byte, saltLength)

	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		iterations,
		memory,
		parallelism,
		keyLength,
	)

	return fmt.Sprintf(
		"$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		memory,
		iterations,
		parallelism,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	), nil
}

func (h *Argon2Hasher) Compare(
	passwordHash string,
	password string,
) error {
	if passwordHash == "" {
		return errors.New("password hash cannot be empty")
	}

	if password == "" {
		return errors.New("password cannot be empty")
	}

	parts := strings.Split(passwordHash, "$")
	if len(parts) != 6 {
		return errors.New("invalid password hash format")
	}

	if parts[1] != "argon2id" {
		return errors.New("unsupported password hash algorithm")
	}

	if parts[2] != "v=19" {
		return errors.New("unsupported argon2 version")
	}

	params := strings.Split(parts[3], ",")
	if len(params) != 3 {
		return errors.New("invalid argon2 parameters")
	}

	memory, err := parseArgon2Parameter(params[0], "m")
	if err != nil {
		return err
	}

	iterations, err := parseArgon2Parameter(params[1], "t")
	if err != nil {
		return err
	}

	parallelism, err := parseArgon2Parameter(params[2], "p")
	if err != nil {
		return err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return fmt.Errorf("decode salt: %w", err)
	}

	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return fmt.Errorf("decode password hash: %w", err)
	}

	actualHash := argon2.IDKey(
		[]byte(password),
		salt,
		uint32(iterations),
		uint32(memory),
		uint8(parallelism),
		uint32(len(expectedHash)),
	)

	if subtle.ConstantTimeCompare(actualHash, expectedHash) != 1 {
		return errors.New("invalid password")
	}

	return nil
}

func parseArgon2Parameter(value string, name string) (int, error) {
	parts := strings.SplitN(value, "=", 2)

	if len(parts) != 2 || parts[0] != name {
		return 0, fmt.Errorf("invalid argon2 parameter: %s", value)
	}

	result, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, fmt.Errorf("invalid argon2 parameter value: %s", value)
	}

	if result <= 0 {
		return 0, fmt.Errorf("argon2 parameter must be positive: %s", value)
	}

	return result, nil
}

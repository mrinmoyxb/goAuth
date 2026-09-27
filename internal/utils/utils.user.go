package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"goAuth/internal/apperrors"
	"strings"

	"golang.org/x/crypto/argon2"
)

type argon2Params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

var defaultParams = &argon2Params{
	memory:      64 * 1024,
	iterations:  3,
	parallelism: 2,
	saltLength:  16,
	keyLength:   32,
}

func UtilHashPassword(password string) (string, error) {
	salt := make([]byte, defaultParams.saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("generate salt: %w", err)
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		defaultParams.iterations,
		defaultParams.memory,
		defaultParams.parallelism,
		defaultParams.keyLength,
	)

	encodedSalt := base64.RawStdEncoding.EncodeToString(salt)
	encodedHash := base64.RawStdEncoding.EncodeToString(hash)

	full := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		defaultParams.iterations,
		defaultParams.parallelism,
		encodedSalt,
		encodedHash,
	)

	return full, nil
}

func utilDecodeHash(encodedHash string) (*argon2Params, []byte, []byte, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return nil, nil, nil, apperrors.ErrInvalidHashFormat
	}

	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
		return nil, nil, nil, apperrors.ErrInvalidHashFormat
	}
	if version != argon2.Version {
		return nil, nil, nil, apperrors.ErrIncompatibleVersion
	}

	params := &argon2Params{}
	if _, err := fmt.Sscanf(
		parts[3],
		"m=%d,t=%d,p=%d",
		&params.memory,
		&params.iterations,
		&params.parallelism,
	); err != nil {
		return nil, nil, nil, apperrors.ErrInvalidHashFormat
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%w: %v", apperrors.ErrInvalidHashFormat, err)
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%w: %v", apperrors.ErrInvalidHashFormat, err)
	}

	return params, salt, hash, nil
}

func UtilVerifyPassword(password, encodeHash string) (bool, error) {
	params, salt, hash, err := utilDecodeHash(encodeHash)
	if err != nil {
		return false, err
	}

	candidateHash := argon2.IDKey(
		[]byte(password),
		salt,
		params.iterations,
		params.memory,
		params.parallelism,
		uint32(len(hash)),
	)

	if subtle.ConstantTimeCompare(hash, candidateHash) == 1 {
		return true, nil
	}

	return false, fmt.Errorf("%w %v", apperrors.ErrInValidPassword, err)
}

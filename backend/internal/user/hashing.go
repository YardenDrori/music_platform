package user

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"golang.org/x/crypto/argon2"
)

const (
	argonTime    = 1
	argonMemory  = 64 * 1024
	argonThreads = 4
	argonKeyLen  = 32
	argonSaltLen = 16
)

type passwordHasher interface {
	hashPassword(password string) string
	verifyPassword(password string, hashedPassword string) (bool, error)
}

type tokenHasher interface {
	hashToken(token string) string
	verifyToken(token string, hashedToken string) bool
}

type sha256TokenHasher struct{}

type argon2idPasswordHasher struct{}

func (h *argon2idPasswordHasher) hashPassword(password string) string {
	salt := rand.Text() //returns 26 runes
	for utf8.RuneCount([]byte(salt)) < argonSaltLen {
		addedSalt := rand.Text()
		salt += addedSalt
	}

	rawHash := argon2.IDKey(
		[]byte(password),
		[]byte(salt)[:argonSaltLen],
		argonTime,
		argonMemory,
		argonThreads,
		argonKeyLen,
	)
	base64Hash := base64.RawStdEncoding.EncodeToString(rawHash)

	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		argonMemory, argonTime, argonThreads,
		salt[:argonSaltLen],
		base64Hash,
	)
}

func (h *argon2idPasswordHasher) verifyPassword(
	password string,
	hashedPassword string,
) (bool, error) {

	//extracting settings
	//settings appear as follow
	//$argon2id$v=19$m=65536,t=1,p=4$<salt>$<hash>
	parts := strings.Split(hashedPassword, "$")
	if len(parts) != 6 {
		return false, fmt.Errorf("verifying password hash: invalid format")
	}

	algo := parts[1]
	if algo != "argon2id" {
		return false, fmt.Errorf("verifying password hash: method incompatible with hash")
	}
	version := parts[2]
	if version != "v=19" {
		return false, fmt.Errorf(
			"verifying password hash: argon2id version mismatch expected v=19, found %v",
			version,
		)
	}

	settings := strings.Split(parts[3], ",")
	if len(settings) != 3 {
		return false, fmt.Errorf("verifying password hash: invalid settings format")
	}
	memory, errMemory := strconv.ParseUint(settings[0][2:], 10, 32)
	if errMemory != nil {
		return false, fmt.Errorf("converting hash setting memory from strings: %w", errMemory)
	}
	time, errTime := strconv.ParseUint(settings[1][2:], 10, 32)
	if errTime != nil {
		return false, fmt.Errorf("converting hash setting time from strings: %w", errTime)
	}
	threads, errThreads := strconv.ParseUint(settings[2][2:], 10, 32)
	if errThreads != nil {
		return false, fmt.Errorf("converting hash setting threads from strings: %w", errThreads)
	}

	salt := parts[4]

	hash := parts[5]
	hashLen := base64.RawStdEncoding.DecodedLen(len(hash))

	rawHash := argon2.IDKey(
		[]byte(password),
		[]byte(salt),
		uint32(time),
		uint32(memory),
		uint8(threads),
		uint32(hashLen),
	)
	base64Hash := base64.RawStdEncoding.EncodeToString(rawHash)

	if subtle.ConstantTimeCompare([]byte(base64Hash), []byte(hash)) == 1 {
		return true, nil
	}
	return false, nil
}

func (h *sha256TokenHasher) hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])
	return tokenHash
}

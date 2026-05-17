package user

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
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

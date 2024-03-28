package cipher

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"time"
)

type Auth interface {
	Validate(string) (bool, error)
}

type TOTPAuth struct {
	user   string
	secret string
}

func (ta *TOTPAuth) Validate(data string) (bool, error) {
	code, err := GetTOTPCode(ta.secret)
	if err != nil {
		return false, err
	}

	if data != code {
		return false, errors.New("wrong one time passwd")
	}

	return true, nil
}

func NewTOTPAuth(user, sec string) (Auth, error) {
	// Secret should be a base32 string
	_, err := base32.StdEncoding.DecodeString(sec)
	if err != nil {
		return nil, err
	}

	return &TOTPAuth{user: user, secret: sec}, nil
}

func GetTOTPCode(secret string) (string, error) {
	// TOTP algorithm refer to https://datatracker.ietf.org/doc/html/rfc6238
	secretBytes, err := base32.StdEncoding.DecodeString(secret)
	if err != nil {
		return "", err
	}

	// Get hash of current unix time with window 30
	t := time.Now()
	tc := t.Unix()
	counter := tc / 30
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(counter))

	// Create a sha1 hash, use write to add data, and use sum to count hash
	hs := hmac.New(sha1.New, secretBytes)
	hs.Write(buf)
	sum := hs.Sum(nil)

	// Get the offset from the last 4 bits
	offset := int(sum[len(sum)-1] & 0x0f)

	// Get the dynamic data of totp
	code := (uint32(sum[offset]&0x7F) << 24) +
		(uint32(sum[offset+1]) << 16) +
		(uint32(sum[offset+2]) << 8) +
		(uint32(sum[offset+3]))

	// Use the last 6 numbers as the totp
	totp := int(int(code) % int(math.Pow10(6)))
	return fmt.Sprintf("%06d", totp), nil
}

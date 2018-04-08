package gokits

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"io"
)

const (
	sha256saltLen   = 8
	sha256itercount = 1000
)

// generate a salt string by special length
func getSalt(len int) string {
	salt := make([]byte, len)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(salt)
}

// generate a hmac sha256 string using salt string and iterate times
func hmacSha256(plaintext, salt string) string {
	bs := []byte(salt)
	mac := hmac.New(sha256.New, bs)
	toencrypt := []byte(plaintext)
	for i := 0; i < sha256itercount; i++ {
		mac.Reset()
		mac.Write(toencrypt)
		mac.Write(bs)
		toencrypt = mac.Sum(nil)
	}
	return hex.EncodeToString(toencrypt)
}

// generate a hmac sha256 string
func GenHmacSha256(plaintext string, saltlen int) string {
	salt := getSalt(sha256saltLen)
	encrypted := hmacSha256(plaintext, salt)
	return encrypted
}

// generate a hmac sha256 string
func GenPasswd(password string, saltlen int) (string, string) {
	salt := getSalt(sha256saltLen)
	encrypted := hmacSha256(password, salt)
	return encrypted, salt
}

// compare the password
func CmpPasswd(passwd, salt, encrypted string) bool {
	nc := hmacSha256(passwd, salt)
	if nc == encrypted {
		return true
	}
	return false
}

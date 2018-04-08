package gokits

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
)

const (
	itercount = 5000
)

var keysalt = []byte{
	0x45, 0xEF, 0x2F, 0x62,
	0xAB, 0x9C, 0xE8, 0x03,
	0x28, 0xB1, 0xF2, 0x61,
	0xDE, 0xF1, 0xD2, 0x58,
}

var (
	ErrAESTextSize = errors.New("ciphertext is not a multiple of the block size")
	ErrAESPadding  = errors.New("cipher padding size error")
	ErrCrc         = errors.New("factor crc is invalid")
)

// the crypto object
type Crypto struct {
	block cipher.Block
	key   []byte
}

func calKey(factor string) []byte {
	fbs := []byte(factor)
	mac := hmac.New(sha256.New, keysalt)
	for i := 0; i < itercount; i++ {
		mac.Reset()
		mac.Write(fbs)
		mac.Write([]byte{byte(i >> 24 & 0xFF), byte(i >> 16 & 0xFF), byte(i >> 8 & 0xFF), byte(i & 0xFF)})
		fbs = mac.Sum(nil)
	}

	blen := aes.BlockSize
	if len(fbs) >= blen {
		return fbs[:blen]
	} else {
		panic("export key fatal failed.")
	}
}

func NewCrypto(factor, crc string) (*Crypto, error) {
	c := new(Crypto)
	return c, c.Init(factor, crc)
}

func (c *Crypto) Init(factor, crc string) error {
	if KermitStr(factor) != crc {
		return ErrCrc
	}

	c.key = calKey(factor)
	block, err := aes.NewCipher([]byte(c.key))
	if err != nil {
		return err
	}
	c.block = block
	return nil
}

func (c *Crypto) Decrypt(src []byte) ([]byte, error) {
	blen := aes.BlockSize

	// check the length
	if len(src) < blen*2 || len(src)%blen != 0 {
		return nil, ErrAESTextSize
	}

	// IV
	iv := src[:blen]
	// encrypt(text)
	srcLen := len(src) - blen
	decryptText := make([]byte, srcLen)

	mode := cipher.NewCBCDecrypter(c.block, iv)
	mode.CryptBlocks(decryptText, src[blen:])

	// unpadding
	paddingLen := int(decryptText[srcLen-1])
	if paddingLen > 16 {
		return nil, ErrAESPadding
	}

	return decryptText[:srcLen-paddingLen], nil
}

// AES解密
func (c *Crypto) DecryptStr(scuritytext string) (string, error) {
	src, err := base64.StdEncoding.DecodeString(scuritytext)
	if err != nil {
		return "", err
	}
	if d, err := c.Decrypt(src); err != nil {
		return "", err
	} else {
		return string(d), err
	}
}

// AES加密
func (c *Crypto) Encrypt(src []byte) ([]byte, error) {
	blen := aes.BlockSize

	// padding
	padLen := blen - (len(src) % blen)
	for i := 0; i < padLen; i++ {
		src = append(src, byte(padLen))
	}

	// iv || encrypt(text)
	srcLen := len(src)
	encryptText := make([]byte, blen+srcLen)

	// iv
	iv := encryptText[:blen]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	mode := cipher.NewCBCEncrypter(c.block, iv)
	mode.CryptBlocks(encryptText[blen:], src)

	return encryptText, nil

}

// AES加密
func (c *Crypto) EncryptStr(plaintext string) (string, error) {
	src := []byte(plaintext)
	if encrypted, err := c.Encrypt(src); err != nil {
		return "", err
	} else {
		return base64.StdEncoding.EncodeToString(encrypted), nil
	}

}

package utils

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"fmt"
)

// DesECBEncrypt
func DesECBEncrypt(src, key []byte, padding string) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return ECBEncrypt(block, src, padding)
}

// DesECBDecrypt
func DesECBDecrypt(src, key []byte, padding string) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	return ECBDecrypt(block, src, padding)
}

//ecb.go

func ECBEncrypt(block cipher.Block, src []byte, padding string) ([]byte, error) {
	blockSize := block.BlockSize()
	src = Padding(padding, src, blockSize)

	encryptData := make([]byte, len(src))

	ecb := NewECBEncrypter(block)
	ecb.CryptBlocks(encryptData, src)

	return encryptData, nil
}

func ECBDecrypt(block cipher.Block, src []byte, padding string) ([]byte, error) {
	dst := make([]byte, len(src))

	mode := NewECBDecrypter(block)
	mode.CryptBlocks(dst, src)

	dst = UnPadding(padding, dst)
	if len(dst) == 0 {
		return []byte{}, fmt.Errorf("err: cannot decrypt, wrong src or key")
	}

	return dst, nil
}

type ecb struct {
	b         cipher.Block
	blockSize int
}

func newECB(b cipher.Block) *ecb {
	return &ecb{
		b:         b,
		blockSize: b.BlockSize(),
	}
}

type ecbEncrypter ecb

// NewECBEncrypter returns a BlockMode which encrypts in electronic code book
// mode, using the given Block.
func NewECBEncrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbEncrypter)(newECB(b))
}

func (x *ecbEncrypter) BlockSize() int { return x.blockSize }

func (x *ecbEncrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Encrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

type ecbDecrypter ecb

// NewECBDecrypter returns a BlockMode which decrypts in electronic code book
// mode, using the given Block.
func NewECBDecrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbDecrypter)(newECB(b))
}

func (x *ecbDecrypter) BlockSize() int { return x.blockSize }

func (x *ecbDecrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Decrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

//padding.go

const PKCS5_PADDING = "PKCS5"
const PKCS7_PADDING = "PKCS7"

func Padding(padding string, src []byte, blockSize int) []byte {
	switch padding {
	case PKCS5_PADDING:
		src = PKCS5Padding(src, blockSize)
	case PKCS7_PADDING:
		src = PKCS7Padding(src, blockSize)
	}
	return src
}

func UnPadding(padding string, src []byte) []byte {
	switch padding {
	case PKCS5_PADDING:
		src = PKCS5Unpadding(src)
	case PKCS7_PADDING:
		src = PKCS7UnPadding(src)
	}
	return src
}

func PKCS5Padding(src []byte, blockSize int) []byte {
	return PKCS7Padding(src, blockSize)
}

func PKCS5Unpadding(src []byte) []byte {
	return PKCS7UnPadding(src)
}

func PKCS7Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func PKCS7UnPadding(src []byte) []byte {
	length := len(src)
	if length == 0 {
		return []byte{}
	}
	unpadding := int(src[length-1])
	end := length - unpadding
	if end < 0 || end > length {
		return []byte{}
	}
	return src[:(length - unpadding)]
}

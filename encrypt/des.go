package encrypt

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"errors"
)

func DesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	origData = PKCS5Padding(origData, block.BlockSize())
	//origData = ZeroPadding(origData, block.BlockSize())
	var blockMode cipher.BlockMode
	if len(key) == 0 {
		blockMode = cipher.NewCBCEncrypter(block, []byte("12345679"))
	} else {
		blockMode = cipher.NewCBCEncrypter(block, key)
	}
	crypted := make([]byte, len(origData))
	// 根据CryptBlocks方法的说明，如下方式初始化crypted也可以
	// crypted := origData
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func DesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	var blockMode cipher.BlockMode
	if len(key) == 0 {
		blockMode = cipher.NewCBCDecrypter(block, []byte("12345679"))
	} else {
		blockMode = cipher.NewCBCDecrypter(block, key)
	}
	origData := make([]byte, len(crypted))
	// origData := crypted
	blockMode.CryptBlocks(origData, crypted)
	origData, err = SafePKCS5UnPadding(origData)
	if err != nil {
		return nil, err
	}
	//origData = ZeroUnPadding(origData)
	return origData, nil
}

// 3DES加密
func TripleDesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	origData = PKCS5Padding(origData, block.BlockSize())
	// origData = ZeroPadding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key[:8])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

// 3DES解密
func TripleDesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, key[:8])
	origData := make([]byte, len(crypted))
	// origData := crypted
	blockMode.CryptBlocks(origData, crypted)
	origData, err = SafePKCS5UnPadding(origData)
	if err != nil {
		return nil, err
	}
	// origData = ZeroUnPadding(origData)
	return origData, nil
}

func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	pad_len := blockSize - len(ciphertext)%blockSize
	if pad_len == 8 {
		return ciphertext
	}

	padtext := bytes.Repeat([]byte{0}, pad_len)
	return append(ciphertext, padtext...)
}

func ZeroUnPadding(origData []byte) []byte {
	return bytes.TrimRightFunc(origData, func(r rune) bool {
		return r == rune(0)
	})
}

//Use PKCS7 ??
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func SafePKCS5UnPadding(origData []byte) ([]byte, error) {
	length := len(origData)
	if length == 0 {
		return nil, errors.New("pkcs5 unpadding error input: empty")
	}
	if length%8 != 0 { //pkcs5 padding with 8 bytes
		return nil, errors.New("pkcs5 unpadding error input: malformed")
	}
	pad_len := int(origData[length-1])
	if pad_len < 0x1 || pad_len > 0x8 {
		return nil, errors.New("pkcs5 unpadding error input: malformed")
	}
	if length < pad_len {
		return nil, errors.New("pkcs5 unpadding error input: malformed")
	}
	return origData[:(length - pad_len)], nil
}

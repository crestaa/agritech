package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

func Encrypt(p string, key string) (string, error) {
	// Key hashing and conversion
	keyBytes, err := hex.DecodeString(MD5(key))
	if err != nil {
		fmt.Println("Errore durante la conversione della chiave:", err)
		return "", err
	}

	// AES init
	iv := keyBytes[:aes.BlockSize]
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		fmt.Println("Errore durante la creazione del blocco AES:", err)
		return "", err
	}

	// Add padding
	padding := aes.BlockSize - len(p)%aes.BlockSize
	for i := 0; i < padding; i++ {
		p += string(padding)
	}

	// Encrypt and return
	ciphertext := make([]byte, len(p))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, []byte(p))
	return string(ciphertext), nil
}

func Decrypt(c string, key string) (string, error) {
	// Key hashing and conversion
	keyBytes, err := hex.DecodeString(MD5(key))
	if err != nil {
		fmt.Println("Errore durante la conversione della chiave:", err)
		return "", err
	}

	// AES init
	iv := keyBytes[:aes.BlockSize]
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		fmt.Println("Errore durante la creazione del blocco AES:", err)
		return "", err
	}

	// Decrypt
	decryptedtext := make([]byte, len(c))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(decryptedtext, []byte(c))

	// Remove padding and return
	unpadding := int(decryptedtext[len(decryptedtext)-1])
	decryptedtext = decryptedtext[:len(decryptedtext)-unpadding]
	return string(decryptedtext), nil
}

func MD5(input string) string {
	hasher := md5.New()
	hasher.Write([]byte(input))
	hashBytes := hasher.Sum(nil)
	return hex.EncodeToString(hashBytes)
}

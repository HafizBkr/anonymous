package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"os"
)

// EmailEncryption gère le chiffrement et déchiffrement des emails
type EmailEncryption struct {
	key []byte
}

// NewEmailEncryption crée une nouvelle instance avec la clé de chiffrement
func NewEmailEncryption() (*EmailEncryption, error) {
	encKey := os.Getenv("EMAIL_ENCRYPTION_KEY")
	if encKey == "" {
		return nil, errors.New("la clé de chiffrement EMAIL_ENCRYPTION_KEY n'est pas définie")
	}
	
	// Hash la clé pour garantir une taille fixe (32 bytes pour AES-256)
	hasher := sha256.New()
	hasher.Write([]byte(encKey))
	key := hasher.Sum(nil)
	
	return &EmailEncryption{
		key: key,
	}, nil
}

// EncryptEmail chiffre l'email et retourne une version encodée en base64
func (ee *EmailEncryption) EncryptEmail(email string) (string, error) {
	block, err := aes.NewCipher(ee.key)
	if err != nil {
		return "", err
	}
	
	// Utilisation d'un IV déterministe (dérivé de la clé)
	iv := ee.key[:aes.BlockSize]
	
	ciphertext := make([]byte, aes.BlockSize+len(email))
	copy(ciphertext[:aes.BlockSize], iv)
	
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(email))
	
	// Encoder en base64 pour le stockage
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptEmail déchiffre un email chiffré en base64
func (ee *EmailEncryption) DecryptEmail(encryptedEmail string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedEmail)
	if err != nil {
		return "", err
	}
	
	block, err := aes.NewCipher(ee.key)
	if err != nil {
		return "", err
	}
	
	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext trop court")
	}
	
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)
	
	return string(ciphertext), nil
}

// IsEncrypted vérifie si un email est déjà chiffré (pour éviter double chiffrement)
func (ee *EmailEncryption) IsEncrypted(email string) bool {
	_, err := base64.StdEncoding.DecodeString(email)
	return err == nil
}
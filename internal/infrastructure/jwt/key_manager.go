package jwt

import (
	"crypto/rsa"
	"fmt"
	"os"
	"sync"

	"github.com/Vanak/Backend/bootstrap"

	"github.com/golang-jwt/jwt/v5"
)

type JWTKeyManager struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	mutex      sync.RWMutex
	isLoaded   bool
	cfg        *bootstrap.Config
}

func NewJwtKeyManager(cfg *bootstrap.Config) *JWTKeyManager {
	jkm := &JWTKeyManager{cfg: cfg}

	err := jkm.LoadKeys(jkm.cfg.Constant.JWTKeysPath.PrivateKey, jkm.cfg.Constant.JWTKeysPath.PublicKey)

	if err != nil {
		panic(err)
	}

	return jkm
}

func (jkm *JWTKeyManager) LoadKeys(privateKeyPath, publicKeyPath string) error {
	jkm.mutex.Lock()
	defer jkm.mutex.Unlock()

	if privateKeyPath == "" || publicKeyPath == "" {
		return fmt.Errorf("jwt private or public key path is empty")
	}

	privateKeyBytes, err := os.ReadFile(privateKeyPath)

	if err != nil {
		return fmt.Errorf("failed to read private key file from %s: %w", privateKeyPath, err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)

	if err != nil {
		return fmt.Errorf("failed to parse private key: %w", err)
	}

	publicKeyBytes, err := os.ReadFile(publicKeyPath)

	if err != nil {
		return fmt.Errorf("failed to read public key file from %s: %w", publicKeyPath, err)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)

	if err != nil {
		return fmt.Errorf("failed to parse public key: %w", err)
	}

	jkm.privateKey = privateKey
	jkm.publicKey = publicKey
	jkm.isLoaded = true

	return nil
}

func (jkm *JWTKeyManager) GetPrivateKey() *rsa.PrivateKey {
	jkm.mutex.RLock()
	defer jkm.mutex.RUnlock()

	if !jkm.isLoaded {
		return nil
	}

	return jkm.privateKey
}

func (jkm *JWTKeyManager) GetPublicKey() *rsa.PublicKey {
	jkm.mutex.RLock()
	defer jkm.mutex.RUnlock()

	if !jkm.isLoaded {
		return nil
	}

	return jkm.publicKey
}

package keystore

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path"
	"strings"
)

type key struct {
	privatePEM string
	publicPEM  string
}

type KeyStore struct {
	store map[string]key
}

func (ks *KeyStore) PrivateKey(kid string) (string, error) {
	key, found := ks.store[kid]
	if !found {
		return "", errors.New("kid lookup failed")
	}

	return key.privatePEM, nil
}

func (ks *KeyStore) PublicKey(kid string) (string, error) {
	key, found := ks.store[kid]
	if !found {
		return "", errors.New("kid lookup failed")
	}

	return key.publicPEM, nil
}

func New() *KeyStore {
	return &KeyStore{
		store: make(map[string]key),
	}
}

func (ks *KeyStore) LoadRSAKeys(fsys fs.FS) error {
	fn := func(fileName string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walkdir failure: %w", err)
		}

		if dirEntry.IsDir() {
			return nil
		}

		if path.Ext(fileName) != ".pem" {
			return nil
		}

		file, err := fsys.Open(fileName)
		if err != nil {
			return fmt.Errorf("opening key file: %w", err)
		}
		defer file.Close()

		pemKey, err := io.ReadAll(io.LimitReader(file, 1024*1024))
		if err != nil {
			return fmt.Errorf("reading auth private key: %w", err)
		}

		privatePEM := string(pemKey)
		publicPEM, err := toPublicPEM(privatePEM)
		if err != nil {
			return fmt.Errorf("converting private PEM to public: %w", err)
		}

		key := key{
			privatePEM: privatePEM,
			publicPEM:  publicPEM,
		}

		ks.store[strings.TrimSuffix(dirEntry.Name(), ".pem")] = key

		return nil
	}

	if err := fs.WalkDir(fsys, ".", fn); err != nil {
		return fmt.Errorf("walking directory: %w", err)
	}

	return nil
}

func toPublicPEM(privatePEM string) (string, error) {
	block, _ := pem.Decode([]byte(privatePEM))
	if block == nil {
		return "", errors.New("invalid key: Key must be a PEM encoded PKCS1 or PKCS8 key")
	}

	var parsedKey any
	parsedKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		parsedKey, err = x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return "", err
		}
	}

	pk, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		return "", errors.New("key is not a valid RSA private key")
	}

	asn1Bytes, err := x509.MarshalPKIXPublicKey(&pk.PublicKey)
	if err != nil {
		return "", fmt.Errorf("marshaling public key: %w", err)
	}

	publicBlock := pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	var buf bytes.Buffer
	if err := pem.Encode(&buf, &publicBlock); err != nil {
		return "", fmt.Errorf("encoding to public PEM: %w", err)
	}

	return buf.String(), nil
}

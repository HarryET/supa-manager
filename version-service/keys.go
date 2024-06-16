package main

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/asn1"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"math/big"
	"net/http"
	"strings"
)

func VerifySignature(pubKey crypto.PublicKey, message, signature []byte) error {
	hashed := sha256.Sum256(message)

	switch key := pubKey.(type) {
	case *rsa.PublicKey:
		return rsa.VerifyPKCS1v15(key, crypto.SHA256, hashed[:], signature)
	case *ecdsa.PublicKey:
		var ecdsaSig struct {
			R, S *big.Int
		}
		if _, err := asn1.Unmarshal(signature, &ecdsaSig); err != nil {
			return err
		}
		if !ecdsa.Verify(key, hashed[:], ecdsaSig.R, ecdsaSig.S) {
			return fmt.Errorf("ecdsa: verification failed")
		}
		return nil
	default:
		return fmt.Errorf("unsupported key type %T", pubKey)
	}
}

type AccountKey struct {
	Key      crypto.PublicKey
	Username string
}

// For each of the accounts load the SSH keys from github i.e. https://github.com/<account>.keys
func LoadSSHKeys(config *Config) ([]AccountKey, error) {
	var sshKeys []AccountKey
	for _, account := range config.PushingAccounts {
		keys, err := loadKeys(account)
		if err != nil {
			return nil, err
		}
		sshKeys = append(sshKeys, keys...)
	}
	return sshKeys, nil
}

func loadKeys(account string) ([]AccountKey, error) {
	// make a http request to github.com/account.keys
	res, err := http.Get(fmt.Sprintf("https://github.com/%s.keys", account))
	if err != nil {
		return nil, err
	}

	// read the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// close the response body
	err = res.Body.Close()
	if err != nil {
		return nil, err
	}

	// split the body by new line
	keys := strings.Split(string(body), "\n")

	// create a slice to store the keys
	var sshKeys []AccountKey
	for _, key := range keys {
		// parse the key
		parsedKey, err := parseKey(key)
		if err != nil {
			return nil, err
		}
		// append the key to the slice
		sshKeys = append(sshKeys, AccountKey{
			Key:      parsedKey,
			Username: strings.ToLower(account),
		})
	}

	// return the keys
	return sshKeys, nil
}

func parseKey(sshPubKey string) (crypto.PublicKey, error) {
	pubKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(sshPubKey))
	if err != nil {
		return nil, err
	}

	switch key := pubKey.(type) {
	case *ssh.Certificate:
		return key.Key.(crypto.PublicKey), nil
	default:
		return key.(crypto.PublicKey), nil
	}
}

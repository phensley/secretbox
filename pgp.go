package main

import (
	"bytes"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
)

const (
	pgpMessageType = "PGP MESSAGE"
)

// Encrypt plaintext using the given key
func encryptPGP(plaintext []byte, key *openpgp.Entity) ([]byte, error) {
	crypted := bytes.Buffer{}
	out, err := openpgp.Encrypt(&crypted, []*openpgp.Entity{key}, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	out.Write([]byte(plaintext))
	out.Close()
	return crypted.Bytes(), nil
}

// Wrap bytes in PGP ASCII armor
func armorPGP(crypted []byte) (string, error) {
	message := bytes.Buffer{}
	out, err := armor.Encode(&message, pgpMessageType, nil)
	if err != nil {
		return "", err
	}
	out.Write(crypted)
	out.Close()
	return message.String(), nil
}

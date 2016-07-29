package main

import (
	crand "crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.com/hashicorp/vault/shamir"
)

var (
	once sync.Once
)

// ShamirParams represents parameters needed to generate a new
// encryption key
type ShamirParams struct {
	Parts     int // total parts to split key into
	Threshold int // number of parts required to decrypt
}

// ShamirKey represents a key that has been split into N parts
// where any M are required to encrypt / decrypt.
type ShamirKey struct {
	Key   [56]byte
	Parts [][]byte
}

func generateShamirKey(params *ShamirParams) *ShamirKey {
	key := generateNaclKey()
	parts, err := shamir.Split(key[:], params.Parts, params.Threshold)
	if err != nil {
		log.Fatalln("Unable to split key", err)
	}
	return &ShamirKey{key, parts}
}

func displayShamirKey(key *ShamirKey, encoding string) {
	for i := range key.Parts {
		displayShamirKeyPart(i, key.Parts[i], encoding)
	}
}

func displayShamirKeyPart(index int, part []byte, encoding string) {
	fmt.Printf(" [%d] %s\n", index+1, encode(part, encoding))
}

func generateNaclKey() [56]byte {
	var key [56]byte
	fillRandomBytes(key[:])
	return key
}

func splitNaclKey(naclKey [56]byte) ([24]byte, [32]byte) {
	var nonce [24]byte
	var key [32]byte
	copy(nonce[:], naclKey[:24])
	copy(key[:], naclKey[24:])
	return nonce, key
}

func fillRandomBytes(raw []byte) {
	if _, err := io.ReadFull(crand.Reader, raw); err != nil {
		log.Panicln(err)
	}
}

func encode(part []byte, enc string) string {
	switch enc {
	case encodingBASE64:
		return base64.StdEncoding.EncodeToString(part)
	default:
		return hex.EncodeToString(part)
	}
}

func decode(part, enc string) ([]byte, error) {
	switch enc {
	case encodingBASE64:
		return base64.StdEncoding.DecodeString(part)
	default:
		return hex.DecodeString(part)
	}
}

// Check if the given path can be read
func canReadFile(path string) bool {
	file, err := os.Open(path)
	defer file.Close()
	return err == nil
}

// Ensure the path exists and is accessible
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// Exit if the value is empty string
func exitOnEmpty(msg, value string) {
	if value == "" {
		exit(1, msg)
	}
}

// Exit if the error value is non-nil
func exitOnError(msg string, err error) {
	if err != nil {
		exit(1, msg, err)
	}
}

// Exit with the given code. Prints ERROR header if code != 0
func exit(code int, msg ...interface{}) {
	if code != 0 {
		fmt.Printf("ERROR: ")
	}
	fmt.Println(msg...)
	os.Exit(code)
}

package main

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/hashicorp/vault/shamir"
)

func canReadFile(fileName string) bool {
	file, err := os.Open(fileName)
	defer file.Close()
	return err == nil
}

func generateEncryptionKey() []byte {
	encryptionKey := make([]byte, 56) // Nonce + Encryption Key
	if _, err := io.ReadFull(rand.Reader, encryptionKey); err != nil {
		log.Panicln(err)
	}
	return encryptionKey
}

func printShamirValues(key []byte) {
	secretBytes, err := shamir.Split(key, 5, 3)
	if err != nil {
		log.Fatalln("Unable to split key ", err)
	}

	for x := range secretBytes {
		fmt.Println(x+1, hex.EncodeToString(secretBytes[x]))
	}
}

func obtainShamirValues() []byte {
	shamirValues := make([][]byte, 3)

	stdinReader := bufio.NewReader(os.Stdin)
	fmt.Println("Please enter at least 3 Shamir secret values")

	for i := 0; i < 3; i++ {
		text, err := stdinReader.ReadString('\n')

		if err != nil {
			log.Fatalln("Unable to read shamir secret", err)
		}

		theBytes, err := hex.DecodeString(text[:len(text)-1]) // Remove newline

		if err != nil {
			log.Fatalln("Unable to decode hex", err, text)
		}
		shamirValues[i] = theBytes
	}

	key, err := shamir.Combine(shamirValues)
	if err != nil {
		log.Fatalln("Unable to obtain key", err)
	}

	return key

}

func makeNaclKey(key []byte) (nonce [24]byte, eKey [32]byte) {
	copy(nonce[:], key[:24])
	copy(eKey[:], key[24:])
	return
}

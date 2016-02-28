package main

import (
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/crypto/nacl/secretbox"
)

func decryptFile(key []byte) {
	if !canReadFile(*fileName) {
		log.Fatalln("Unable to read file ", *fileName)
	}

	fileBytes, err := ioutil.ReadFile(*fileName)
	if err != nil {
		log.Fatalln("Unable to open file", err)
	}

	nonce, eKey := makeNaclKey(key)

	plainBytes, ok := secretbox.Open(nil, fileBytes, &nonce, &eKey)
	if !ok {

		log.Fatalln("Decryption failed")
	}

	ioutil.WriteFile(*fileName+"decrypted", plainBytes, os.FileMode(400))
}

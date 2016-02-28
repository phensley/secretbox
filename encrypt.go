package main

import (
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/crypto/nacl/secretbox"
)

func encryptFile(key []byte) {
	nonce, eKey := makeNaclKey(key)

	fileBytes, err := ioutil.ReadFile(*fileName)
	if err != nil {
		log.Fatal(err)
	}

	outBytes := secretbox.Seal(nil, fileBytes, &nonce, &eKey)

	ioutil.WriteFile(*fileName+"encrypted", outBytes, os.FileMode(400))

}

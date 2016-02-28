package main

import (
	"flag"
	"os"
)

var (
	encrypt  = flag.Bool("encrypt", false, "Specify this flag to encrypt the file")
	decrypt  = flag.Bool("decrypt", false, "Specify this flag to decrypt the file")
	fileName = flag.String("file", "", "The file you wish to preform the operation on")
)

func main() {
	flag.Parse()

	switch {
	case *encrypt && *decrypt:
		flag.PrintDefaults()
		os.Exit(-1)
	case !*encrypt && !*decrypt:
		flag.PrintDefaults()
		os.Exit(-1)
	case *fileName == "":
		flag.PrintDefaults()
		os.Exit(-1)
	}

	if *encrypt {
		encryptionKey := generateEncryptionKey()
		printShamirValues(encryptionKey)
		encryptFile(encryptionKey)
		return
	}

	if *decrypt {
		key := obtainShamirValues()
		decryptFile(key)
		return
	}

}

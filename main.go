package main

import (
	"fmt"
	"io/ioutil"
	"os"
	filepath "path"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/nacl/secretbox"
)

const (
	encodingHEX    = "hex"
	encodingBASE64 = "base64"
)

var (
	cmd = &cobra.Command{
		Use:   os.Args[0],
		Short: "secretbox",
	}

	cmdEncrypt = &cobra.Command{
		Use:   "encrypt a file",
		Short: "encrypts a file",
		Run:   runEncrypt,
	}

	cmdDecrypt = &cobra.Command{
		Use:   "decrypt a file",
		Short: "decrypts a file",
		Run:   runDecrypt,
	}

	parts     int
	threshold int
	input     string
	output    string
	encoding  string
)

func init() {
	cmdEncrypt.Flags().IntVarP(&parts, "parts", "p", 5, "total parts to split key into")
	cmdEncrypt.Flags().IntVarP(&threshold, "threshold", "t", 3, "minimum parts needed to decrypt")

	// Flags common to both encrypt / decrypt
	addFlags := func(cmd *cobra.Command, op string) {
		cmd.Flags().StringVarP(&input, "input", "i", "", fmt.Sprintf("file to %s", op))
		cmd.Flags().StringVarP(&output, "output", "o", "", fmt.Sprintf("destination for %sed file", op))
		cmd.Flags().StringVarP(&encoding, "encoding", "e", encodingHEX, "key encoding to use (hex or base64)")
	}

	addFlags(cmdEncrypt, "encrypt")
	addFlags(cmdDecrypt, "decrypt")

	cmd.AddCommand(cmdEncrypt, cmdDecrypt)
}

func main() {
	cmd.Execute()
}

// Encrypt the input file
func runEncrypt(cmd *cobra.Command, args []string) {
	checkEncoding()
	checkFiles()

	plain, err := ioutil.ReadFile(input)
	exitOnError("Reading input file", err)

	params := &ShamirParams{
		Parts:     parts,
		Threshold: threshold,
	}
	shamirKey := generateShamirKey(params)

	fmt.Printf("Generated secret key in %d parts with thresdhold %d:\n\n", parts, threshold)
	displayShamirKey(shamirKey, encoding)

	fmt.Printf("\nEncrypting to '%s'\n", output)

	nonce, key := splitNaclKey(shamirKey.Key)
	crypted := secretbox.Seal(nil, plain, &nonce, &key)
	err = ioutil.WriteFile(output, crypted, os.FileMode(400))
	exitOnError("Failed to write output file:", err)
	fmt.Println("\nSuccess!")
}

// Decrypt the input file
func runDecrypt(cmd *cobra.Command, args []string) {
	checkEncoding()
	checkFiles()

	crypted, err := ioutil.ReadFile(input)
	exitOnError("Reading input file", err)

	naclKey := obtainShamirKey()
	nonce, key := splitNaclKey(naclKey)
	plain, ok := secretbox.Open(nil, crypted, &nonce, &key)
	if !ok {
		exit(1, "Decryption failed!")
	}

	err = ioutil.WriteFile(output, plain, os.FileMode(400))
	exitOnError("Failed to write output file:", err)
	fmt.Println("\nSuccess!")
}

// Ensure the encoding is supported
func checkEncoding() {
	switch encoding {
	case encodingHEX, encodingBASE64:
	default:
		exit(1, fmt.Sprintf("unknown encoding '%s'", encoding))
	}
}

// Check the input / output files are correct
func checkFiles() {
	// Check if the variables are set and not equivalent
	exitOnEmpty("you must provide the 'input' filename", input)
	exitOnEmpty("you must provide the 'output' filename", output)

	// Must never overwrite input
	if strings.Compare(input, output) == 0 {
		exit(1, "input and output paths must not be the same!")
	}

	// Ensure the input file exists
	if !fileExists(input) {
		exit(1, fmt.Sprintf("input file '%s' must exist", input))
	}

	// Ensure output file doesn't already exist.
	if fileExists(output) {
		exit(1, fmt.Sprintf("output file '%s' already exists! not overwriting.", output))
	}

	// Ensure output file's parent directory exists
	parent := filepath.Dir(output)
	if !fileExists(parent) {
		exit(1, fmt.Sprintf("output directory '%s' does not exist", parent))
	}
}

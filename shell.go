package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/google/shlex"
	"gopkg.in/readline.v1"

	"github.com/hashicorp/vault/shamir"
)

const (
	shellHelp = `
This interactive shell will allow you to enter the key parts.

Commands:

 add <part>    - adds a key part
 list          - view the parts that have been entered
 del <num>     - deletes the part by #
 done          - indicate you have entered parts and are ready to decrypt

 exit          - exit immediately without decrypting
 help          - display this message
`
)

func obtainShamirKey() [56]byte {
	parts := [][]byte{}

	cli := setupReadline()
	defer cli.Close()
	showShellHelp()

	for {
		command, args, err := parseCommand(cli.Readline())
		if err != nil {
			if err != io.EOF && err.Error() != "Interrupt" {
				fmt.Println("error: ", err)
			}
			continue
		}

		switch command {
		case "add":
			if len(args) == 0 {
				fmt.Println("expected an encoded key")
				continue
			}
			fmt.Println("add:", args[0])
			part, err := decode(args[0], encoding)
			if err != nil {
				fmt.Printf("failed to decode key part using %s: %s\n", encoding, err)
				continue
			}
			parts = append(parts, part)

		case "del":
			if len(args) == 0 {
				fmt.Println("expected a key number")
				continue
			}
			index, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil || index < 1 || int(index) > len(parts) {
				fmt.Printf("key index is invalid")
				continue
			}
			// Convert 0-based to 1-based
			fmt.Printf("deleting key at index %d\n", index)
			parts = append(parts[:index-1], parts[index:]...)

		case "done":
			// Check key for validity. If invalid, provide another chance to
			// correct one or more key parts.
			fmt.Printf("%d keys entered. validating.\n", len(parts))

			key, err := shamir.Combine(parts)
			if err != nil {
				fmt.Println("error combining key parts:", err)
				continue
			}
			if len(key) != 56 {
				fmt.Printf("expected key of length 56. got length: %d\n", len(key))
				continue
			}

			// Key looks good. Correct the type and return
			var result [56]byte
			copy(result[:], key)
			return result

		case "exit":
			fmt.Println("quitting without decrypting")
			os.Exit(0)

		case "help":
			showShellHelp()

		case "list":
			for i := range parts {
				displayShamirKeyPart(i, parts[i], encoding)
			}
		}
	}
}

func setupReadline() *readline.Instance {
	completer := readline.NewPrefixCompleter(
		readline.PcItem("add"),
		readline.PcItem("del"),
		readline.PcItem("done"),
		readline.PcItem("exit"),
		readline.PcItem("help"),
		readline.PcItem("list"),
	)

	cli, err := readline.NewEx(&readline.Config{
		Prompt:       ">> ",
		AutoComplete: completer,
	})

	exitOnError("failed to configure readline", err)
	return cli
}

func parseCommand(line string, err error) (string, []string, error) {
	if err != nil {
		return "", nil, err
	}
	parts, err := shlex.Split(line)
	if err != nil {
		return "", nil, err
	}
	if len(parts) == 0 {
		return "", []string{}, nil
	}
	return strings.ToLower(parts[0]), parts[1:], nil
}

func showShellHelp() {
	fmt.Println(shellHelp)
}

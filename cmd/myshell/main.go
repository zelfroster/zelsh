package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type BuiltinCommands struct {
	Exit, Echo, Type string
}

func (b *BuiltinCommands) IsValid(cmd string) bool {
	return cmd == b.Exit || cmd == b.Echo || cmd == b.Type
}

var Builtins = &BuiltinCommands{
	Exit: "exit",
	Echo: "echo",
	Type: "type",
}

// @TODO: Make type command work exactly as actual shell
// func EvalCommand(cmd string) string {}

func checkIfFileInPaths(fp string) (bool, string) {
	paths := strings.Split(os.Getenv("PATH"), ":")
	for _, path := range paths {
		filePath := filepath.Join(path, fp)
		if _, err := os.Stat(filePath); err == nil {
			return true, path
		}
	}
	return false, ""

	// ------------------ EASIER IMPLEMENTATION ------------------
	// if path, err := exec.LookPath(fullCommand[1]); err != nil {
	// 	retmsg = fmt.Sprintf("%s: not found", fullCommand[1])
	// } else {
	// 	retmsg = fmt.Sprintf("%s is %s", fullCommand[1], path)
	// }
	// -----------------------------------------------------------
}

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		// Wait for user input
		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			log.Fatalln(err)
		}

		fullCommand := strings.Split(input[:len(input)-1], " ")

		var retmsg string

		switch fullCommand[0] {
		case Builtins.Exit:
			os.Exit(0)
		case Builtins.Echo:
			retmsg = strings.Join(fullCommand[1:], " ")
		case Builtins.Type:
			if Builtins.IsValid(fullCommand[1]) {
				retmsg = fmt.Sprintf("%s is a shell builtin", fullCommand[1])
			} else {
				exists, path := checkIfFileInPaths(fullCommand[1])
				if exists {
					retmsg = fmt.Sprintf("%s is %s/%s", fullCommand[1], path, fullCommand[1])
				} else {
					retmsg = fmt.Sprintf("%s: not found", fullCommand[1])
				}
			}
		default:
			// Execute command if found in provided PATH else print not found
			if exists, _ := checkIfFileInPaths(fullCommand[0]); exists {
				shellCmd := exec.Command(fullCommand[0], fullCommand[1:]...)
				stdout, err := shellCmd.Output()
				if err != nil {
					log.Fatalln(err)
				}
				retmsg = strings.Trim(string(stdout), "\r\n")
			} else {
				retmsg = fmt.Sprintf("%s: command not found", fullCommand[0])
			}
		}

		fmt.Fprint(os.Stdout, retmsg+"\n")
	}
}

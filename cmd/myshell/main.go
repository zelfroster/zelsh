package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type BuiltinCommands struct {
	Exit, Echo, Type, Pwd, Cd string
}

func (b *BuiltinCommands) IsValid(cmd string) bool {
	return cmd == b.Exit || cmd == b.Echo || cmd == b.Type || cmd == b.Pwd || cmd == b.Cd
}

var Builtins = &BuiltinCommands{
	Exit: "exit",
	Echo: "echo",
	Type: "type",
	Pwd:  "pwd",
	Cd:   "cd",
}

//  TODO: Make type command work exactly as actual shell
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

// TODO: Test if nested quotes work exactly as in zsh/bash
func parseInput(inputString string) []string {
	var fullCommand []string
	current := ""
	inQuotes := false
	inDoubleQuotes := false
	for i := 0; i < len(inputString); i++ {
		char := inputString[i]

		if char == '\'' && !inDoubleQuotes {
			inQuotes = !inQuotes
			continue
		}

		if char == '"' && !inQuotes {
			inDoubleQuotes = !inDoubleQuotes
			continue
		}

		if char == ' ' {
			if !inQuotes && !inDoubleQuotes && current != "" {
				fullCommand = append(fullCommand, current)
				current = ""
			} else if inQuotes || inDoubleQuotes {
				current += string(char)
			}
		} else {
			current += string(char)
		}
	}

	if current != "" {
		fullCommand = append(fullCommand, current)
	}

	return fullCommand
}

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		// Wait for user input
		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			log.Fatalln(err)
		}

		fullCommand := parseInput(input[:len(input)-1])

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
		case Builtins.Pwd:
			pwd, err := os.Getwd()
			if err != nil {
				log.Fatalln(err)
			}
			retmsg = fmt.Sprintf("%s", pwd)
		case Builtins.Cd:
			path := fullCommand[1]
			if path == "~" {
				path = os.Getenv("HOME")
			}
			err := os.Chdir(path)
			if err != nil {
				var pathError *os.PathError
				if errors.As(err, &pathError) {
					retmsg = fmt.Sprintf("cd: %s: No such file or directory", path)
				} else {
					log.Fatalln(err)
				}
			} else {
				continue
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

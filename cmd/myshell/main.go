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
	// 	retMsg = fmt.Sprintf("%s: not found", fullCommand[1])
	// } else {
	// 	retMsg = fmt.Sprintf("%s is %s", fullCommand[1], path)
	// }
	// -----------------------------------------------------------
}

// TODO: Test if nested quotes work exactly as in zsh/bash
func parseInput(inputString string) (string, []string) {
	var fullCommand []string
	current := ""
	inQuotes := false
	inDoubleQuotes := false
	backslash := false
	for i := 0; i < len(inputString); i++ {
		char := inputString[i]

		switch char {
		case '\'':
			if backslash && inDoubleQuotes {
				current += string('\\')
			}
			if backslash || inDoubleQuotes {
				current += string(char)
			} else {
				inQuotes = !inQuotes
			}
			backslash = false

		case '"':
			if backslash || inQuotes {
				current += string(char)
			} else {
				inDoubleQuotes = !inDoubleQuotes
			}
			backslash = false

		case '\\':
			if backslash || inQuotes {
				current += string(char)
				backslash = false
			} else {
				backslash = true
			}

		case ' ':
			if backslash && inDoubleQuotes {
				current += string('\\')
			}
			if backslash || inQuotes || inDoubleQuotes {
				current += string(char)
			} else if current != "" {
				fullCommand = append(fullCommand, current)
				current = ""
			}
			backslash = false

		default:
			if backslash && inDoubleQuotes {
				current += string('\\')
			}
			current += string(char)
			backslash = false
		}
	}

	if current != "" {
		fullCommand = append(fullCommand, current)
	}

	return fullCommand[0], fullCommand[1:]
}

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		// Wait for user input
		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			log.Fatalln(err)
		}

		cmd, args := parseInput(input[:len(input)-1])

		// check for redirect operator, if found then store them separately
		redirectOperatorIndex := -1
		var operator string
		var remArgs []string

		for i, val := range args {
			if val == ">" || val == "1>" || val == "2>" || val == ">>" || val == "1>>" || val == "2>>" {
				redirectOperatorIndex = i
				operator = val
				break
			}
		}

		if redirectOperatorIndex != -1 {
			remArgs = args[redirectOperatorIndex+1:]
			args = args[:redirectOperatorIndex]
		}

		var retMsg string
		var errMsg string

		switch cmd {
		case Builtins.Exit:
			os.Exit(0)

		case Builtins.Echo:
			retMsg = strings.Join(args, " ")

		case Builtins.Type:
			if Builtins.IsValid(args[0]) {
				retMsg = fmt.Sprintf("%s is a shell builtin", args[0])
			} else {
				exists, path := checkIfFileInPaths(args[0])
				if exists {
					retMsg = fmt.Sprintf("%s is %s/%s", args[0], path, args[0])
				} else {
					retMsg = fmt.Sprintf("%s: not found", args[0])
				}
			}

		case Builtins.Pwd:
			pwd, err := os.Getwd()
			if err != nil {
				log.Fatalln(err)
			}
			retMsg = fmt.Sprintf("%s", pwd)

		case Builtins.Cd:
			path := args[0]
			if path == "~" {
				path = os.Getenv("HOME")
			}
			err := os.Chdir(path)
			if err != nil {
				var pathError *os.PathError
				if errors.As(err, &pathError) {
					retMsg = fmt.Sprintf("cd: %s: No such file or directory", path)
				} else {
					log.Fatalln(err)
				}
			} else {
				continue
			}

		case "cat", "ls":
			for _, arg := range args {
				shellCmd := exec.Command(cmd, arg)
				stdout, err := shellCmd.Output()
				if err != nil {
					if operator == "2>" || operator == "2>>" {
						errMsg += fmt.Sprintf("%s: %s: No such file or directory\n", cmd, arg)
					} else {
						fmt.Fprintf(os.Stdout, "%s: %s: No such file or directory\n", cmd, arg)
					}
				} else {
					if operator == ">" || operator == "1>" {
						retMsg += strings.Trim(string(stdout), "\r\n")
					} else if operator == ">>" || operator == "1>>" {
						retMsg += fmt.Sprintf("%s\n", strings.Trim(string(stdout), "\r\n"))
					} else if operator == "2>" || operator == "2>>" {
						fmt.Fprintf(os.Stdout, "%s\n", strings.Trim(string(stdout), "\r\n"))
					} else {
						fmt.Fprintf(os.Stdout, "%s", strings.Trim(string(stdout), "\r\n"))
					}
				}
			}

		default:
			// Execute command if found in provided PATH else print not found
			if exists, _ := checkIfFileInPaths(cmd); exists {
				shellCmd := exec.Command(cmd, args...)
				// fmt.Println(shellCmd.String())
				stdout, err := shellCmd.Output()
				if err != nil {
					log.Fatalln(err)
				}
				retMsg = strings.Trim(string(stdout), "\r\n")
			} else {
				retMsg = fmt.Sprintf("%s: command not found", cmd)
			}

		}

		if len(operator) != 0 && len(remArgs) > 0 {
			if len(remArgs) > 1 {
				retMsg += " " + strings.Join(remArgs[1:], " ")
			}

			// for creating path/to/dir/file: we need to check if path/to/dir exists
			dirArr := strings.Split(remArgs[0], "/")
			if len(dirArr) > 1 {
				dir := strings.Join(dirArr[:len(dirArr)-1], "/")
				if _, err := os.Stat(dir); err != nil {
					retMsg = fmt.Sprintf("no such file or directory: %s", remArgs[0])
				}
			}

			var err error
			switch operator {
			case ">", "1>":
				err = os.WriteFile(remArgs[0], []byte(retMsg), 0644)
			case "2>":
				err = os.WriteFile(remArgs[0], []byte(errMsg), 0644)
			case ">>", "1>>":
				f, err := os.OpenFile(remArgs[0], os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
				if err != nil {
					log.Fatalln(err)
				}
				defer f.Close()
				_, err = f.WriteString("\n" + retMsg)
			case "2>>":
				f, err := os.OpenFile(remArgs[0], os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
				if err != nil {
					log.Fatalln(err)
				}
				defer f.Close()
				_, err = f.WriteString("\n" + errMsg)
			}

			if err != nil {
				log.Fatalln(err)
			}

			if cmd == Builtins.Echo && operator == "2>" {
				fmt.Fprint(os.Stdout, retMsg+"\n")
			}

			continue
		}

		fmt.Fprint(os.Stdout, retMsg+"\n")
	}
}

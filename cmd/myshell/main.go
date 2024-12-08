package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
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

func CommandNotFound(cmd string) string {
	return fmt.Sprintf("%s: command not found", cmd)
}

// @TODO: Make type command work exactly as actual shell
// func EvalCommand(cmd string) string {}

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
				retmsg = fmt.Sprintf("%s: not found", fullCommand[1])
			}
		default:
			retmsg = CommandNotFound(fullCommand[0])
		}

		fmt.Fprint(os.Stdout, retmsg+"\n")
	}
}

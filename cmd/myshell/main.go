package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

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
		case "exit":
			os.Exit(0)
		case "echo":
			retmsg = strings.Join(fullCommand[1:], " ")
		default:
			retmsg = fmt.Sprintf("%s: command not found", fullCommand[0])
		}

		fmt.Fprint(os.Stdout, retmsg+"\n")
	}
}

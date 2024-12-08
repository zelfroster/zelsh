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
		// Uncomment this block to pass the first stage
		fmt.Fprint(os.Stdout, "$ ")

		// Wait for user input
		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			log.Fatalln(err)
		}

		fullCommand := strings.Split(input[:len(input)-1], " ")

		if fullCommand[0] == "exit" {
			if len(fullCommand) > 1 && fullCommand[1] == "0" {
				os.Exit(0)
			} else {
				os.Exit(1)
			}
		}

		retmsg := fmt.Sprintf("%s: command not found\n", fullCommand[0])
		fmt.Fprint(os.Stdout, retmsg)
	}
}

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	for {
		// Uncomment this block to pass the first stage
		fmt.Fprint(os.Stdout, "$ ")

		// Wait for user input
		str, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			log.Fatalln(err)
		}

		retmsg := fmt.Sprintf("%s: command not found\n", str[:len(str)-1])
		fmt.Fprint(os.Stdout, retmsg)
	}
}

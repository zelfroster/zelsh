package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	// Uncomment this block to pass the first stage
	fmt.Fprint(os.Stdout, "$ ")

	// Wait for user input
	reader, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		log.Fatalln(err)
	}

	bytes := []byte(reader)
	str := string(bytes[:len(bytes)-1])

	retmsg := fmt.Sprintf("%s: command not found\n", str)
	fmt.Fprint(os.Stdout, retmsg)
}

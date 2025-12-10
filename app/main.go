package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

func main() {
	// TODO: Uncomment the code below to pass the first stage
	scanner := bufio.NewScanner(os.Stdin)
	for {

		fmt.Print("$ ")
		var cmd string
		scanner.Scan()
		cmd = scanner.Text()
		tokens := strings.Fields(cmd)

		switch strings.ToLower(tokens[0]) {
		case "echo":
			echo(tokens)

		case "exit":
			break

		default:
			fmt.Print(tokens[0])
			fmt.Println(": command not found")

		}

		if cmd == "exit" {
			break
		}

	}

}

func echo(tokens []string) {
	println(strings.Join(tokens[1:], " "))
}

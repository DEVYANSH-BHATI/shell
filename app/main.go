package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
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
		if scanner.Err() != nil {
			continue
		}
		cmd = scanner.Text()

		tokens := strings.Fields(cmd)

		if len(tokens) == 0 {
			continue
		}

		switch strings.ToLower(tokens[0]) {
		case "echo":
			echo(tokens)

		case "type":
			typee(strings.ToLower(tokens[1]))

		case "exit":

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

func typee(token string) {
	cmds := []string{"echo", "exit"}
	if slices.Contains(cmds, token) {
		fmt.Println(token, "is a shell builtin")
	} else {
		fmt.Println(token + ": not found")
	}

}

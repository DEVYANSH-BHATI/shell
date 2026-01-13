package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
)

// go run main.go executableHandling.go parser.go
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

		// tokens := strings.Fields(cmd)
		tokens := tokenize(cmd)
		// for _, token := range tokens {
		// 	fmt.Println(token)
		// }

		if len(tokens) == 0 {
			continue
		}
		switch strings.ToLower(tokens[0]) {
		case "echo":
			echo(tokens)

		case "type":
			typee(strings.ToLower(tokens[1]))

		case "pwd":
			pwd, _ := os.Getwd()
			fmt.Println(pwd)

		case "cd":
			var err error
			if tokens[1] == "~" {
				homepath := os.Getenv("HOME")
				err = os.Chdir(homepath)
			} else {
				err = os.Chdir(tokens[1])
			}
			if err != nil {
				fmt.Println("cd: " + tokens[1] + ": No such file or directory")
			}

		case "ls":
			dir, _ := os.Getwd()
			directories, err := os.ReadDir(dir)
			// err := cmd.Run()
			if err != nil {
				fmt.Println(err)
			} else {
				for _, dir := range directories {
					fmt.Println(dir.Name())
				}
			}

		case "exit":

		default:
			foundExecutable, _ := findExecutable(tokens[0])
			if foundExecutable {
				args := tokens[1:]
				cmd := exec.Command(tokens[0], args...)
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr

				err := cmd.Run()
				if err != nil {
					fmt.Println(err)
				}
			} else {
				fmt.Print(tokens[0])
				fmt.Println(": command not found")
			}

			// break

		}

		if strings.ToLower(cmd) == "exit" {
			break
		}

	}
}

func echo(tokens []string) {

	println(string(strings.Join(tokens[1:], " ")))
}

func typee(token string) {
	cmds := []string{"echo", "exit", "type", "pwd", "cd", "ls"}
	if slices.Contains(cmds, token) {
		fmt.Println(token, "is a shell builtin")
		return
	}
	foundExecutable, pathOfExecutable := findExecutable(token)
	if foundExecutable {
		fmt.Println(token, "is", pathOfExecutable)
	} else {
		fmt.Println(token + ": not found")

	}

}

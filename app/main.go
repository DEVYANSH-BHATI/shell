package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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
			// findExecutable(strings.ToLower((tokens[1])))

		case "exit":

		default:
			// fmt.Print(tokens[0])
			// fmt.Println(": command not found")
			args := tokens[1:]
			exec.Command(tokens[0], args...)

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
	cmds := []string{"echo", "exit", "type"}
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

func findExecutable(token string) (bool, string) {
	cur_os := runtime.GOOS
	pathEnv := os.Getenv("PATH")
	if pathEnv == "" {
		fmt.Println("PATH environment variable is not set.")
		return false, "emptypath"
	}
	if cur_os == "windows" {
		directories := strings.Split(pathEnv, ";")

		// // break pathenv into different directories
		// iterate on each dir and check if tokenpath exists in that dir
		// check if its regular file and is executable
		for _, dir := range directories {

			// fmt.Println(dir)
			tokenpath := filepath.Join(dir + "/" + token)
			info, err := os.Stat(tokenpath)
			if err == nil {
				if info.Mode().IsRegular() {
					if info.Mode()&0111 != 0 {
						return true, tokenpath
					}
				}
			}
		}
		return false, ""

		// break

		// }
	} else {
		directories := strings.Split(pathEnv, ":")
		for _, dir := range directories {

			tokenpath := filepath.Join(dir + "/" + token)
			info, err := os.Stat(tokenpath)
			if err == nil {
				if info.Mode().IsRegular() {
					if info.Mode()&0111 != 0 {
						return true, tokenpath
					}
				}
			}
		}

		// break
		return false, ""
	}
}

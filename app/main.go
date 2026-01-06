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

		if cmd == "exit" {
			break
		}

	}
}

func tokenize(cmd string) []string {
	//tokenize without strings.fields
	tokens := []string{}
	inSingleQuote := false
	var currentToken strings.Builder
	for _, char := range cmd {
		if char == '\'' && inSingleQuote == false {
			inSingleQuote = true
			continue
		}

		if inSingleQuote {
			if char == '\'' {
				inSingleQuote = false
				continue
			}
			currentToken.WriteRune(char)
			continue
		}
		if char != ' ' || char != '\t' {
			currentToken.WriteRune(char)
			continue
		}
		if currentToken.Len() > 0 {
			tokens = append(tokens, currentToken.String())
			currentToken.Reset()
		}
	}
	if currentToken.Len() > 0 {
		tokens = append(tokens, currentToken.String())
		currentToken.Reset()
	}
	return tokens
}

// func tokenize(cmd string) []string {
// 	var tokens []string
// 	var current strings.Builder
// 	inSingleQuote := false

// 	for _, ch := range cmd {
// 		switch ch {
// 		case '\'':
// 			inSingleQuote = !inSingleQuote

// 		case ' ', '\t':
// 			if inSingleQuote {
// 				current.WriteRune(ch)
// 			} else if current.Len() > 0 {
// 				tokens = append(tokens, current.String())
// 				current.Reset()
// 			}

// 		default:
// 			current.WriteRune(ch)
// 		}
// 	}

// 	if current.Len() > 0 {
// 		tokens = append(tokens, current.String())
// 	}

// 	return tokens
// }

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

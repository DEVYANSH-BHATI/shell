package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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

		// Redirection logic :
		// if command has 1> or > means out needs redirection
		// so we will define new io.writer (output) by default same as os.stdout that is out cli
		// but if > is detected we take the filename from command and open the file for edit 2
		var output io.Writer = os.Stdout
		var erroutput io.Writer = os.Stderr
		var err_openfile *os.File
		var openfile *os.File
		tokens_copy := tokens
		var filepath_redirection string
		redirectionexists := slices.Index(tokens, ">")
		stderr_redirection_exists := slices.Index(tokens, "2>")
		var stderr_filepath_redirection string
		if stderr_redirection_exists != -1 {
			tokens_copy = tokens[:stderr_redirection_exists]
			stderr_filepath_redirection = tokens[stderr_redirection_exists+1]
			err_openfile, err := os.OpenFile(stderr_filepath_redirection,
				os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
				0644)
			if err != nil {
				fmt.Println(err)
				continue
			}
			erroutput = err_openfile

		}
		if redirectionexists == -1 {
			redirectionexists = slices.Index(tokens, "1>")
		}
		if redirectionexists != -1 {
			tokens_copy = tokens[:redirectionexists]
			if redirectionexists+1 >= len(tokens) {
				fmt.Println("syntax error")
			}
			filepath_redirection = tokens[redirectionexists+1]
			openfile, err := os.OpenFile(filepath_redirection,
				os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
				0644)
			if err != nil {
				fmt.Println(err)
				continue
			}
			output = openfile

		}

		if len(tokens_copy) == 0 {
			continue
		}
		switch strings.ToLower(tokens_copy[0]) {
		case "echo":
			echo(tokens_copy, output)

		case "type":
			typee(strings.ToLower(tokens_copy[1]), output, erroutput)

		case "pwd":
			pwd, _ := os.Getwd()
			fmt.Fprintln(output, pwd)

		case "cd":
			var err error
			if tokens_copy[1] == "~" {
				homepath := os.Getenv("HOME")
				err = os.Chdir(homepath)
			} else {
				err = os.Chdir(tokens_copy[1])
			}
			if err != nil {
				fmt.Fprintln(output, "cd: "+tokens_copy[1]+": No such file or directory")
			}

		case "exit":

		default:
			foundExecutable, pathOfExecutable := findExecutable(tokens_copy[0])
			if foundExecutable {
				args := tokens_copy[1:]
				cmd := exec.Command(pathOfExecutable, args...)
				cmd.Stdin = os.Stdin
				cmd.Stdout = output
				cmd.Stderr = erroutput
				cmd.Args[0] = tokens_copy[0]
				cmd.Run()

			} else {
				fmt.Print(tokens_copy[0])
				fmt.Println(": command not found")
			}

			// break

		}
		if openfile != nil {
			openfile.Close()
		}
		if err_openfile != nil {
			err_openfile.Close()
		}

		if strings.ToLower(cmd) == "exit" {
			break
		}

	}
}

func echo(tokens []string, output io.Writer) {

	fmt.Fprintln(output, string(strings.Join(tokens[1:], " ")))
}

func typee(token string, output io.Writer, erroutput io.Writer) {
	cmds := []string{"echo", "exit", "type", "pwd", "cd", "ls"}
	if slices.Contains(cmds, token) {
		fmt.Fprintln(output, token, "is a shell builtin")
		return
	}
	foundExecutable, pathOfExecutable := findExecutable(token)
	if foundExecutable {
		fmt.Fprintln(output, token, "is", pathOfExecutable)
	} else {
		fmt.Fprintln(erroutput, token+": not found")

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

// custom parser
func tokenize(cmd string) []string {
	//tokenize without strings.fields
	tokens := []string{}
	inSingleQuote := false
	inDoubleQuote := false
	escapedCharacter := false
	var currentToken strings.Builder
	for _, char := range cmd {

		// single quote overpower any escapes and everything is treated as literalchar
		if inSingleQuote {
			if char == '\'' {
				inSingleQuote = false
				continue
			}
			currentToken.WriteRune(char)
			continue
		}

		// fmt.Println(currentToken.String())
		// fmt.Println(currentToken.String(), string(char), escapedCharacter, inDoubleQuote, inSingleQuote)
		if escapedCharacter {
			// ", \, $, `, and newline
			if !inSingleQuote && !inDoubleQuote {
				currentToken.WriteRune(char)
				escapedCharacter = false
				continue
			}
			if char == '"' || char == '\\' || char == '$' || char == '\n' || char == ' ' {
				// fmt.Println("spl", char)
				currentToken.WriteRune(char)
				escapedCharacter = false
				continue
			}
		}

		if char == '\'' {
			if inDoubleQuote {
				currentToken.WriteRune(char)
				continue
			}
			inSingleQuote = !inSingleQuote
			continue
		}

		if char == '"' {
			inDoubleQuote = !inDoubleQuote
			continue
		}

		if char == '\\' {
			if inSingleQuote {
				currentToken.WriteRune(char)
				continue
			}
			escapedCharacter = true
			continue
		}

		if char == ' ' || char == '\t' {
			if inSingleQuote || inDoubleQuote {
				currentToken.WriteRune(char)
				continue
			} else {
				if currentToken.Len() > 0 {
					tokens = append(tokens, currentToken.String())
					currentToken.Reset()
				}
				continue
			}
		}

		currentToken.WriteRune(char)
	}
	if currentToken.Len() > 0 {
		tokens = append(tokens, currentToken.String())
		currentToken.Reset()
	}
	return tokens
}

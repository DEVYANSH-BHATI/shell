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

type Command struct {
	Name    string
	Args    []string
	Stdout  io.Writer
	Stderr  io.Writer
	Cleanup func()
}

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

		command := parsecommand(tokens)

		if len(command.Args) == 0 {
			continue
		}
		switch strings.ToLower(command.Args[0]) {
		case "echo":
			echo(command.Args, command.Stdout)

		case "type":
			typee(strings.ToLower(command.Args[1]), command.Stdout, command.Stderr)

		case "pwd":
			pwd, _ := os.Getwd()
			fmt.Fprintln(command.Stdout, pwd)

		case "cd":
			var err error
			if command.Args[1] == "~" {
				homepath := os.Getenv("HOME")
				err = os.Chdir(homepath)
			} else {
				err = os.Chdir(command.Args[1])
			}
			if err != nil {
				fmt.Fprintln(command.Stdout, "cd: "+command.Args[1]+": No such file or directory")
			}

		case "exit":

		default:
			foundExecutable, pathOfExecutable := findExecutable(command.Args[0])
			if foundExecutable {
				args := command.Args[1:]
				cmd := exec.Command(pathOfExecutable, args...)
				cmd.Stdin = os.Stdin
				cmd.Stdout = command.Stdout
				cmd.Stderr = command.Stderr
				cmd.Args[0] = command.Args[0]
				cmd.Run()

			} else {
				fmt.Print(command.Args[0])
				fmt.Println(": command not found")
			}

			// break

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

func parsecommand(tokens []string) *Command {
	// Redirection logic :
	// if command has 1> or > means out needs redirection
	// so we will define new io.writer (output) by default same as os.stdout that is out cli
	// but if > is detected we take the filename from command and open the file for edit 2

	cmd := &Command{
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
		Cleanup: func() {},
	}

	for pos := 0; pos < len(tokens); pos++ {
		// for pos, cur_token := range tokens {
		cur_token := tokens[pos]
		if cur_token == "1>" || cur_token == ">" {
			if pos+1 >= len(tokens) {
				fmt.Println("syntax error")
				return cmd
			}
			var openfile *os.File
			openfile, err := os.OpenFile(tokens[pos+1],
				os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
				0644)
			if err != nil {
				fmt.Println(err)
				return cmd
			}
			cmd.Stdout = openfile
			output_file_cleanup := cmd.Cleanup
			cmd.Cleanup = func() {
				output_file_cleanup()
				openfile.Close()
			}
			pos++
			continue
		}

		if cur_token == "2>" {
			if pos+1 >= len(tokens) {
				fmt.Println("syntax error")
				return cmd
			}
			var err_openfile *os.File
			err_openfile, err := os.OpenFile(tokens[pos+1],
				os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
				0644)
			if err != nil {
				fmt.Println(err)
				return cmd
			}
			cmd.Stderr = err_openfile
			err_file_cleanup := cmd.Cleanup
			cmd.Cleanup = func() {
				err_file_cleanup()
				err_openfile.Close()
			}
			pos++
			continue
		}
		cmd.Args = append(cmd.Args, cur_token)
	}
	if len(cmd.Args) > 0 {
		cmd.Name = cmd.Args[0]
	}

	return cmd
}

// // var output io.Writer = os.Stdout
// // var erroutput io.Writer = os.Stderr
// // var err_openfile *os.File
// // var openfile *os.File
// cmd_object.Args := tokens
// var filepath_redirection string
// // redirectionexists := slices.Index(tokens, ">")
// stderr_redirection_exists := slices.Index(tokens, "2>")
// var stderr_filepath_redirection string
// if stderr_redirection_exists != -1 {
// 	cmd_object.Args = tokens[:stderr_redirection_exists]
// 	stderr_filepath_redirection = tokens[stderr_redirection_exists+1]
// 	err_openfile, err := os.OpenFile(stderr_filepath_redirection,
// 		os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
// 		0644)
// 	if err != nil {
// 		fmt.Println(err)
// 		continue
// 	}
// 	erroutput = err_openfile

// }
// if redirectionexists == -1 {
// 	redirectionexists = slices.Index(tokens, "1>")
// }
// if redirectionexists != -1 {
// 	cmd_object.Args = tokens[:redirectionexists]
// 	if redirectionexists+1 >= len(tokens) {
// 		fmt.Println("syntax error")
// 	}
// 	filepath_redirection = tokens[redirectionexists+1]
// 	openfile, err := os.OpenFile(filepath_redirection,
// 		os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
// 		0644)
// 	if err != nil {
// 		fmt.Println(err)
// 		continue
// 	}
// 	output = openfile

// }

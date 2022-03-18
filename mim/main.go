package main

import (
	"fmt"
	"os"
	"strings"
)

type Variable struct {
	Name  string
	Value uint8
}

type StringVariable struct {
	Name  string
	Value string
}

type Program struct {
	code       []string
	variables  []Variable
	stringVars []StringVariable
}

var ourProgram Program

func parseLine(line string) (string, string) {
	// make line lowercase
	line = strings.ToLower(line)
	// first word is the instruction
	instruction := line[:strings.Index(line, " ")]
	// rest is the arguments
	args := line[len(instruction)+1:]
	return instruction, args
}

func handleInstruction(instruction string, arguments string) (string, error) {
	switch instruction {
	case "print":
		return handlePrint(arguments)
	case "let":
		return handleLet(arguments)
	default:
		// return error
		return "", fmt.Errorf("invalid instruction")
	}
}

func main() {
	// first argument should be infile
	// second argument should be outfile (optional)
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go infile [outfile]")
		os.Exit(1)
	}

	infile := os.Args[1]
	outfile := os.Args[1] + ".mips" // compile to mips assembly code
	if len(os.Args) > 2 {
		outfile = os.Args[2]
	}

}

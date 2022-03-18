package main

import (
	"fmt"
	"strings"
)

func handleLet(args string) error {
	// make sure there are two arguments
	argsArray := strings.Split(args, " ")
	if len(argsArray) != 2 {
		return fmt.Errorf("let: invalid number of arguments")
	}
	// first argument should be the variable name
	varName := argsArray[0]
}

func handlePrint(s string) (string, error) {
	string := ""
	// check if followed by a quote
	if s[0] == '"' && s[len(s)-1] == '"' {
		// remove quotes and add to string
		string = s[1 : len(s)-1]
	} else if s[0] >= '0' && s[0] <= '9' {
		// are all digits?
		if strings.Index(s, " ") == -1 {
			// yes, convert to int and add to string
			string = fmt.Sprintf("%d", s)
		} else {
			// no, return error
			return "", fmt.Errorf("print: invalid argument")
		}
	} else if len(s) > 0 {
		// check if variable exists
		for _, variable := range ourProgram.variables {
			if variable.Name == s {
				// yes, add to string
				string = fmt.Sprintf("%d", variable.Value)
			}
		}
		// check if string variable exists
		for _, variable := range ourProgram.stringVars {
			if variable.Name == s {
				// yes, add to string
				string = variable.Value
			}
		}
	} else {
		// no, return error
		return "", fmt.Errorf("print: invalid argument")
	}
	// return string

	// todo: convert string to bytecode

	return string, nil
}

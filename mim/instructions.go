package main

import (
	"fmt"
	"strconv"
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
	// second argument should be the value
	varValue := argsArray[1]

	// for debugging
	fmt.Printf("varName: %s\n", varName)
	fmt.Printf("varValue: %s\n", varValue)

	// if the value is future, then it will be assigned later
	if varValue == "future" {
		// add the variable to the future map
		ourProgram.futureVars = append(ourProgram.futureVars, FutureVariable{
			Name: varName,
		})
		return nil
	} else {
		// check if variable is a number
		if _, err := strconv.Atoi(varValue); err == nil {
			fmt.Println("is number")
			// if it is a number, then assign it to the variable
			i, err := strconv.Atoi(varValue)
			if err != nil {
				return fmt.Errorf("let: %s", err)
			}
			ourProgram.variables = append(ourProgram.variables, Variable{
				Name:  varName,
				Value: uint8(i),
			})
			return nil
		} else {
			// if it is not a number, then check for quotes
			if varValue[0] == '"' && varValue[len(varValue)-1] == '"' {
				// if it is a string, then assign it to the variable
				ourProgram.stringVars = append(ourProgram.stringVars, StringVariable{
					Name:  varName,
					Value: varValue[1 : len(varValue)-1],
				})
				return nil
			} else {
				// throw error
				return fmt.Errorf("let: invalid value: %s", varValue)
			}
		}
	}
}

func (c *Context) handleAddi(args string) error {
	// args should be var1, var2, number
	// or var1, 0, number
	argsArray := strings.Split(args, " ")
	if len(argsArray) != 3 {
		return fmt.Errorf("addi: invalid number of arguments")
	}
	// first argument should be a variable name
	varAname := argsArray[0]
	// second argument should be a variable name
	varBname := argsArray[1]
	// third argument should be a number
	varValue, err := strconv.Atoi(argsArray[2])
	if err != nil {
		return fmt.Errorf("addi: invalid number")
	}
	tmpResult := 0
	// check if variable B is 0
	if varBname == "0" {
		// todo
	}
	// check if variable B is a number (it cannot be a future variable)
	for _, variable := range ourProgram.variables {
		if variable.Name == varBname {
			// find variable A
			for _, variableA := range ourProgram.variables {
				if variableA.Name == varAname {
					// add the two variables
					tmpResult = varValue + int(variable.Value)
					// check if the result is a number
					if tmpResult > 255 {
						return fmt.Errorf("addi: result is too large")
					}
					// assign the result to the variable
					variableA.Value = uint8(tmpResult)
					return nil
				}
			}
			// if we're still here, variable A might be a future variable
			// if so, remove it from futures and add it to variables
			for i, futureVariable := range ourProgram.futureVars {
				if futureVariable.Name == varAname {
					ourProgram.variables = append(ourProgram.variables, Variable{
						Name:  varAname,
						Value: uint8(variable.Value + uint8(varValue)),
					})
					ourProgram.futureVars = append(ourProgram.futureVars[:i], ourProgram.futureVars[i+1:]...)
					return nil
				}
			}
			// variable A doesn't exist
			return fmt.Errorf("addi: variable does not exist: %s", varAname)
		}
	}
	return fmt.Errorf("addi: variable does not exist: %s", varBname)
}

func (c *Context) handlePrint(s string) error {
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
			return fmt.Errorf("print: invalid argument")
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
		return fmt.Errorf("print: invalid argument")
	}

	// now time to convert to assembly
	// for each character, get the ascii value
	var asciiCodes []int
	for _, char := range string {
		asciiCodes = append(asciiCodes, int(char))
	}
	/*
		we need to generate
		addi (unused register), $0, (ascii code)
		sw (previous register), 0(mem address of output)
	*/
	// addi (unused register), $0, (ascii code)
	ioReg := c.FindUnusedTemporaryRegister()

	// setup the register
	c.AddInstruction(&Instruction{
		Opcode:        "addi",
		Args:          fmt.Sprintf("$t%d, $0, 0x30", ioReg),
		RegistersUsed: []uint8{ioReg},
	}, false)
	c.AddInstruction(&Instruction{
		Opcode:        "sll",
		Args:          fmt.Sprintf("$t%d, $t%d, 24", ioReg, ioReg),
		RegistersUsed: []uint8{ioReg},
	}, false)
	c.AddInstruction(&Instruction{
		Opcode:        "addi",
		Args:          fmt.Sprintf("$t%d, $t%d, 4", ioReg, ioReg),
		RegistersUsed: []uint8{ioReg},
	}, false)

	tmpCharReg := c.FindUnusedTemporaryRegister()
	// for each ascii code, generate assembly
	for _, asciiCode := range asciiCodes {
		c.AddInstruction(&Instruction{
			Opcode:        "addi",
			Args:          fmt.Sprintf("$t%d, $0, %d", tmpCharReg, asciiCode),
			RegistersUsed: []uint8{tmpCharReg},
		}, false)
		c.AddInstruction(&Instruction{
			Opcode:        "sw",
			Args:          fmt.Sprintf("$t%d, 0($t%d)", tmpCharReg, ioReg),
			RegistersUsed: []uint8{ioReg, tmpCharReg},
		}, false)
	}

	// print a newline
	c.AddInstruction(&Instruction{
		Opcode:        "addi",
		Args:          fmt.Sprintf("$t%d, $0, 10", tmpCharReg),
		RegistersUsed: []uint8{ioReg},
	}, false)
	c.AddInstruction(&Instruction{
		Opcode:        "sw",
		Args:          fmt.Sprintf("$t%d, 0($t%d)", tmpCharReg, ioReg),
		RegistersUsed: []uint8{ioReg, ioReg},
	}, false)

	// free the registers
	c.ReleaseTemporaryRegister(ioReg)
	c.ReleaseTemporaryRegister(tmpCharReg)

	return nil
}

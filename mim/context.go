package main

func CreateMainContext() *Context {
	// make map of temporary registers in use
	tmpRegisters := make(map[uint8]bool)
	// set each to false
	var i uint8
	for i = 0; i <= 9; i++ {
		tmpRegisters[i] = false
	}
	return &Context{
		TemporaryRegistersInUse: tmpRegisters,
		Instructions:            []*Instruction{},
	}
}

func (c *Context) FindUnusedTemporaryRegister() uint8 {
	var i uint8
	for i = 0; i <= 9; i++ {
		if !c.TemporaryRegistersInUse[i] {
			c.TemporaryRegistersInUse[i] = true
			return i
		}
	}
	return 0
}

func (c *Context) ReleaseTemporaryRegister(reg uint8) {
	c.TemporaryRegistersInUse[reg] = false
}

func (c *Context) AddInstruction(instruction *Instruction, freeTmp bool) {
	// if freeTmp is true, then free all registers used by the instruction
	if freeTmp {
		for _, reg := range instruction.RegistersUsed {
			c.ReleaseTemporaryRegister(reg)
		}
	}
	c.Instructions = append(c.Instructions, instruction)
}

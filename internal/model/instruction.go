package model

type Instruction interface {
	Execute()
}

type printInstruction struct {
	message string
}

func (pi *printInstruction) Execute() {
	println(pi.message)
}

func NewPrintInstruction(message string) Instruction {
	return &printInstruction{message: message}
}

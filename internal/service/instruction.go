package service

import (
	"artisan-backend/internal/geo"
	"fmt"
	"math"
)

type Instruction interface {
	Execute()
	String() string
}

type printInstruction struct {
	message string
}

type forwardInstruction struct {
	px int
}

type leftInstruction struct {
	deg int
}

type rightInstruction struct {
	deg int
}

type reverseInstruction struct {
}

type inversionInstruction struct {
}

func DegToRad(deg float64) float64 {
	return deg * math.Pi / 180
}

func NewPrintInstruction(message string) Instruction {
	return &printInstruction{message: message}
}

func (pi *printInstruction) Execute() {
	println(pi.message)
}

func (pi *printInstruction) String() string {
	return fmt.Sprintf("print %s", pi.message)
}

func NewForwardInstruction(px int) Instruction {
	return &forwardInstruction{px: px}
}

func (fi *forwardInstruction) Execute() {
	line := geo.NewLine(currentPosition.Origin, currentPosition.Origin.Add(currentPosition.Direction.MulScalar(float64(fi.px))))
	geometries = append(geometries, line)
	currentPosition = geo.NewRay(line.Target, currentPosition.Direction)
}

func (fi *forwardInstruction) String() string {
	return fmt.Sprintf("forward %d", fi.px)
}

func NewLeftInstruction(deg int) Instruction {
	return &leftInstruction{deg: deg}
}

func (li *leftInstruction) Execute() {
	currentPosition = geo.NewRay(currentPosition.Origin, geo.Rotate(currentPosition.Direction, DegToRad(float64(li.deg))))
}

func (li *leftInstruction) String() string {
	return fmt.Sprintf("left %d", li.deg)
}

func NewRightInstruction(deg int) Instruction {
	return &rightInstruction{deg: deg}
}
func (ri *rightInstruction) Execute() {
	currentPosition = geo.NewRay(currentPosition.Origin, geo.Rotate(currentPosition.Direction, -DegToRad(float64(ri.deg))))
}

func (ri *rightInstruction) String() string {
	return fmt.Sprintf("right %d", ri.deg)
}

func NewReverseInstruction() Instruction {
	return nil
}

func (ri *reverseInstruction) Execute() {
}

func (ri *reverseInstruction) String() string {
	return "reverse"
}

func NewInversionInstruction() Instruction {
	return &inversionInstruction{}
}

func (ii *inversionInstruction) String() string {
	return "inversion"
}

func (ii *inversionInstruction) Execute() {
	inverted = !inverted
}

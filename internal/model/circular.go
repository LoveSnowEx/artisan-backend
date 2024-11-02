package model

import (
	"artisan-backend/pkg/tool"

	"github.com/golang-ds/linkedlist/circularly"
)

type Instruction interface{}

type Circular = circularly.LinkedList[Instruction]

type Circulars []*Circular

func NewCirculars() Circulars {
	return make(Circulars, 0)
}

func NewCircular() *Circular {
	return tool.Ptr(circularly.New[Instruction]())
}

func (c *Circulars) Enqueue(circular *Circular) {
	*c = append(*c, circular)
}

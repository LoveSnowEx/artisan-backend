package model

import (
	"github.com/golang-ds/linkedlist/circularly"
	queue "github.com/golang-ds/queue/linkedqueue"
)

type Instruction string

type Circular = circularly.LinkedList[Instruction]

type Circulars = queue.LinkedQueue[Circular]

func NewCirculars() Circulars {
	return queue.New[Circular]()
}

func NewCircular() Circular {
	return circularly.New[Instruction]()
}

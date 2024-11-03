package service

type Circular struct {
	current      int
	instructions []Instruction
}

func NewCircular() *Circular {
	return &Circular{
		instructions: make([]Instruction, 0),
	}
}

type Circulars []*Circular

func NewCirculars() Circulars {
	return make(Circulars, 0)
}

func (c *Circular) Size() int {
	return len(c.instructions)
}

func (c *Circular) First() (Instruction, bool) {
	if len(c.instructions) == 0 {
		return nil, false
	}
	return c.instructions[c.current], true
}

func (c *Circular) AddLast(instruction Instruction) {
	c.instructions = append(c.instructions, instruction)
}

func (c *Circular) Rotate() {
	if len(c.instructions) == 0 {
		return
	}
	c.current = (c.current + 1) % len(c.instructions)
}

func (c *Circular) Values() []string {
	size := c.Size()
	values := make([]string, 0, size)
	for i := 0; i < size; i++ {
		values = append(values, c.instructions[i].String())
	}
	return values
}

func (cs *Circulars) Enqueue(circular *Circular) {
	*cs = append(*cs, circular)
}

func (cs *Circulars) Values() [][]string {
	size := len(*cs)
	values := make([][]string, 0, size)
	for i := 0; i < size; i++ {
		values = append(values, (*cs)[i].Values())
	}
	return values
}

package geo

import "github.com/fogleman/gg"

type Canvas struct {
	*gg.Context
}

func NewCanvas(width, height int) *Canvas {
	c := &Canvas{
		gg.NewContext(width, height),
	}
	c.SetLineWidth(2)
	return c
}

func (c *Canvas) Draw(g Geometry) {
	if g == nil {
		return
	}
	g.Draw(c.Context)
}

func (c *Canvas) DrawMany(gs []Geometry) {
	for _, g := range gs {
		c.Draw(g)
	}
}

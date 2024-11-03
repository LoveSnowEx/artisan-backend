package geo

import "github.com/fogleman/gg"

type canvas struct {
	*gg.Context
}

func NewCanvas(width, height int) *canvas {
	c := &canvas{
		gg.NewContext(width, height),
	}
	c.SetLineWidth(2)
	return c
}

func (c *canvas) Draw(g Geometry) {
	if g == nil {
		return
	}
	g.Draw(c.Context)
}

func (c *canvas) DrawMany(gs []Geometry) {
	for _, g := range gs {
		c.Draw(g)
	}
}

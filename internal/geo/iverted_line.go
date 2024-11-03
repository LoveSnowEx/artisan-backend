package geo

import (
	"github.com/deeean/go-vector/vector2"
	"github.com/fogleman/gg"
)

type InvertedLine struct {
	Line
}

func (l *InvertedLine) Draw(dst *gg.Context) {
	direction := l.Target.Sub(l.Origin)
	ray1 := NewRay(l.Origin, vector2.New(-direction.Y, direction.X))
	ray2 := NewRay(l.Origin, vector2.New(direction.Y, -direction.X))
	ray1.Draw(dst)
	ray2.Draw(dst)
}

func (l *InvertedLine) Invert(c *Circle) Geometry {
	return l
}

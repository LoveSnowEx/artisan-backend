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
	invertedOrigin := InvertPoint(l.Origin, c)
	invertedTarget := InvertPoint(l.Target, c)

	angle1 := Angle(invertedOrigin.Sub(c.Origin))
	angle2 := Angle(invertedTarget.Sub(c.Origin))
	angleCircle := Angle(c.Origin.Sub(invertedOrigin.Lerp(invertedTarget, 0.5)))

	if IsAngleBetween(angleCircle, angle1, angle2) {
		return &Arc{
			Origin: c.Origin,
			Radius: c.Radius,
			Angle1: angle1,
			Angle2: angle2,
		}
	}

	// not implemented
	return l
}

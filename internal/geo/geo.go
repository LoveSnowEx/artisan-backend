package geo

import (
	"math"

	"github.com/deeean/go-vector/vector2"
	"github.com/fogleman/gg"
)

const threshold = 1e-9

type InvertedLine struct {
	Ray1 *Ray
	Ray2 *Ray
}

type Geometry interface {
	Draw(dst *gg.Context)
	Invert(*Circle) Geometry
}

type Drawable interface {
	Draw(g Geometry)
	DrawMany(gs []Geometry)
}

type canvas struct {
	*gg.Context
}

func Rotate(v *vector2.Vector2, angle float64) *vector2.Vector2 {
	return vector2.New(
		v.X*math.Cos(angle)-v.Y*math.Sin(angle),
		v.X*math.Sin(angle)+v.Y*math.Cos(angle),
	)
}

func Center(p1, p2 *vector2.Vector2) *vector2.Vector2 {
	return p1.Add(p2).DivScalar(2)
}

func InvertPoint(p *vector2.Vector2, c *Circle) *vector2.Vector2 {
	// Calculate the distance from the circle origin to the point
	d := c.Origin.Distance(p)

	// If point is at the origin of the inversion circle, return nil to avoid division by zero
	if d < threshold {
		return nil
	}

	// Inverted distance (R^2 / d)
	invertedDistance := (c.Radius * c.Radius) / d

	// Direction vector from circle origin to point p
	direction := p.Sub(c.Origin).Normalize()

	// Scaled inverted point
	invertedPoint := c.Origin.Add(direction.MulScalar(invertedDistance))

	return invertedPoint
}

func (c *canvas) Draw(g Geometry) {
	if g == nil {
		return
	}
	g.Draw(c.Context)
}

func NewCanvas(width, height int) *canvas {
	c := &canvas{
		gg.NewContext(width, height),
	}
	c.SetLineWidth(2)
	return c
}

func (l *InvertedLine) Draw(dst *gg.Context) {
	l.Ray1.Draw(dst)
	l.Ray2.Draw(dst)
}

func (l *InvertedLine) Invert(c *Circle) Geometry {
	return l
}

package geo

import (
	"image/color"
	"math"

	"github.com/deeean/go-vector/vector2"
	"github.com/fogleman/gg"
)

type Circle struct {
	Origin *vector2.Vector2
	Radius float64
}

func (c *Circle) Draw(dst *gg.Context) {
	// Draw the circle
	clr := color.RGBA{0, 255, 0, 255} // Green color for the circle

	dst.SetColor(clr)
	dst.DrawCircle(c.Origin.X, float64(dst.Height())-c.Origin.Y, c.Radius)
	dst.Stroke()
}

func (c *Circle) Invert(oc *Circle) Geometry {
	if math.Abs(c.Origin.Distance(oc.Origin)-c.Radius) < threshold {
		invertedOrigin := InvertPoint(c.Origin, oc)
		direction := c.Origin.Sub(oc.Origin)
		perpendicular := vector2.New(-direction.Y, direction.X)
		return &InvertedLine{
			Ray1: &Ray{
				Origin:    invertedOrigin,
				Direction: perpendicular,
			},
			Ray2: &Ray{
				Origin:    invertedOrigin,
				Direction: perpendicular.MulScalar(-1),
			},
		}
	}
	return c
}

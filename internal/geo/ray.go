package geo

import (
	"math"

	"github.com/deeean/go-vector/vector2"
	"github.com/fogleman/gg"
)

type Ray struct {
	Origin    *vector2.Vector2
	Direction *vector2.Vector2
}

func NewRay(origin, direction *vector2.Vector2) *Ray {
	return &Ray{
		Origin:    origin,
		Direction: direction.Normalize(),
	}
}

func (r *Ray) Draw(dst *gg.Context) {
	line := &Line{
		Origin: r.Origin,
		Target: r.Origin.Add(r.Direction.MulScalar(math.Inf(1))),
	}
	line.Draw(dst)
}

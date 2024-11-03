package geo

import (
	"math"

	"github.com/deeean/go-vector/vector2"
)

func Rotate(v *vector2.Vector2, angle float64) *vector2.Vector2 {
	return vector2.New(
		v.X*math.Cos(angle)-v.Y*math.Sin(angle),
		v.X*math.Sin(angle)+v.Y*math.Cos(angle),
	)
}

func Center(p1, p2 *vector2.Vector2) *vector2.Vector2 {
	return p1.Add(p2).DivScalar(2)
}

func Angle(v *vector2.Vector2) float64 {
	return math.Atan2(v.Y, v.X)
}

func IsAngleBetween(angle, a, b float64) bool {
	var lowerBound, upperBound float64
	if a < b {
		lowerBound = math.Ceil((a - angle) / (2 * math.Pi))
		upperBound = math.Floor((b - angle) / (2 * math.Pi))
	} else {
		lowerBound = math.Ceil((b - angle) / (2 * math.Pi))
		upperBound = math.Floor((a - angle) / (2 * math.Pi))
	}
	return lowerBound < upperBound
}

func InvertPoint(p *vector2.Vector2, c *Circle) *vector2.Vector2 {
	// Calculate the distance from the circle origin to the point
	distance := c.Origin.Distance(p)

	// If point is at the origin of the inversion circle, return nil to avoid division by zero
	if distance < threshold {
		return nil
	}

	// Direction vector from circle origin to point p
	direction := p.Sub(c.Origin)

	// Scalar to scale the direction vector
	scalar := math.Pow(c.Radius, 2) / math.Pow(distance, 2)

	// Scaled inverted point
	if scalar < threshold {
		return c.Origin
	}
	invertedPoint := c.Origin.Add(direction.MulScalar(scalar))

	return invertedPoint
}

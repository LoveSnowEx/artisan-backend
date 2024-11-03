package geo

import (
	"image/color"
	"math"

	"github.com/deeean/go-vector/vector2"
	"github.com/fogleman/gg"
)

type Arc struct {
	Origin *vector2.Vector2
	Radius float64
	Angle1 float64
	Angle2 float64
}

func (a *Arc) Draw(dst *gg.Context) {
	// Draw the curve
	clr := color.RGBA{255, 0, 0, 255} // Red color for the curve

	dst.SetColor(clr)
	dst.DrawArc(a.Origin.X, float64(dst.Height())-a.Origin.Y, a.Radius, -a.Angle1, -a.Angle2)
	dst.Stroke()
}

func (a *Arc) Invert(c *Circle) Geometry {
	direction := a.Origin.Sub(c.Origin)
	distance := direction.Distance(vector2.New(0, 0))
	// Invert radius
	// The radius of the arc is not inverted directly; we need to find the inversion of two points on the arc boundary
	angle1 := a.Angle1
	angle2 := a.Angle2
	point1 := a.Origin.Add(Rotate(vector2.New(a.Radius, 0), angle1))
	point2 := a.Origin.Add(Rotate(vector2.New(a.Radius, 0), angle2))

	// Invert the boundary points of the arc
	invertedPoint1 := InvertPoint(point1, c)
	invertedPoint2 := InvertPoint(point2, c)

	// If one of the points is at the center of the inversion circle (nil case), return nil
	if invertedPoint1 == nil && invertedPoint2 == nil {
		return nil // This would be a degenerate case
	}
	if invertedPoint1 == nil {
		angle1 = (angle1 + angle2) / 2
		point1 = a.Origin.Add(Rotate(vector2.New(a.Radius, 0), angle1))
		invertedPoint1 = InvertPoint(point1, c)
	}
	if invertedPoint2 == nil {
		angle2 = (angle1 + angle2) / 2
		point2 = a.Origin.Add(Rotate(vector2.New(a.Radius, 0), angle2))
		invertedPoint2 = InvertPoint(point2, c)
	}

	if math.Abs(distance-a.Radius) < threshold {
		angleCircle := Angle(c.Origin.Sub(a.Origin))
		if IsAngleBetween(angleCircle, angle1, angle2) {
			return &InvertedLine{
				Line: Line{
					Origin: invertedPoint1,
					Target: invertedPoint2,
				},
			}
		}
		return &Line{
			Origin: invertedPoint1,
			Target: invertedPoint2,
		}
	}

	centerPoint := Center(
		InvertPoint(a.Origin.Add(direction.Normalize().MulScalar(a.Radius)), c),
		InvertPoint(a.Origin.Add(direction.Normalize().MulScalar(-a.Radius)), c),
	)
	angle1 = math.Atan2(invertedPoint1.Y-centerPoint.Y, invertedPoint1.X-centerPoint.X)
	angle2 = math.Atan2(invertedPoint2.Y-centerPoint.Y, invertedPoint2.X-centerPoint.X)

	return &Arc{
		Origin: centerPoint,
		Radius: centerPoint.Distance(invertedPoint1),
		Angle1: angle1,
		Angle2: angle2,
	}
}

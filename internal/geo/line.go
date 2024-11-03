package geo

import (
	"image/color"
	"math"

	"github.com/deeean/go-vector/vector2"
	"github.com/fogleman/gg"
)

type Line struct {
	Origin *vector2.Vector2
	Target *vector2.Vector2
}

func (l *Line) Draw(dst *gg.Context) {
	// Draw the line
	clr := color.RGBA{0, 0, 255, 255} // Blue color for the line
	dst.SetColor(clr)
	origin := l.Origin.Copy()
	target := l.Target.Copy()
	direction := target.Sub(origin)
	// clamp the line to the canvas
	if direction.X == 0 {
		if origin.Y < 0 {
			origin.Y = 0
		} else if origin.Y > float64(dst.Height()) {
			origin.Y = float64(dst.Height())
		}
		if target.Y < 0 {
			target.Y = 0
		} else if target.Y > float64(dst.Height()) {
			target.Y = float64(dst.Height())
		}
	} else {
		m := direction.Y / direction.X
		b := origin.Y - m*origin.X
		if origin.X < 0 {
			origin.X = 0
			origin.Y = b
		} else if origin.X > float64(dst.Width()) {
			origin.X = float64(dst.Width())
			origin.Y = m*origin.X + b
		}
		if target.X < 0 {
			target.X = 0
			target.Y = b
		} else if target.X > float64(dst.Width()) {
			target.X = float64(dst.Width())
			target.Y = m*target.X + b
		}
	}
	dst.DrawLine(origin.X, float64(dst.Height())-origin.Y, target.X, float64(dst.Height())-target.Y)
	dst.Stroke()
}

func (l *Line) Invert(c *Circle) Geometry {
	// Invert the endpoints of the line
	invertedOrigin := InvertPoint(l.Origin, c)
	invertedTarget := InvertPoint(l.Target, c)

	// If one of the points is at the center of the inversion circle (nil case), return nil
	if invertedOrigin == nil || invertedTarget == nil {
		return nil // This would be a degenerate case
	}

	// If the original line passes through the inversion circle's center, it remains a line
	originToCenter := l.Origin.Sub(c.Origin)
	targetToCenter := l.Target.Sub(c.Origin)
	crossProduct := (originToCenter.X * targetToCenter.Y) - (originToCenter.Y * targetToCenter.X)

	if math.Abs(crossProduct) < threshold {
		// The points lie on a line through the origin, so inversion gives another line
		return &Line{
			Origin: invertedOrigin,
			Target: invertedTarget,
		}
	}

	direction := l.Target.Sub(l.Origin).Normalize()
	projection := direction.MulScalar(direction.Dot(originToCenter))
	perpendicular := originToCenter.Sub(projection)

	invertedPerpendicular := InvertPoint(perpendicular.Add(c.Origin), c)
	angle1 := Angle(invertedOrigin.Sub(invertedPerpendicular))
	angle2 := Angle(invertedTarget.Sub(invertedPerpendicular))
	return &Arc{
		Origin: Center(invertedPerpendicular, c.Origin),
		Radius: invertedPerpendicular.Distance(c.Origin) / 2,
		Angle1: angle1,
		Angle2: angle2,
	}
}

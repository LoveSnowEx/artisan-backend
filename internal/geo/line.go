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
	orgin := l.Origin
	target := l.Target
	direction := target.Sub(orgin)
	// clamp the line to the canvas
	if direction.X == 0 {
		if orgin.Y < 0 {
			orgin.Y = 0
		} else if orgin.Y > float64(dst.Height()) {
			orgin.Y = float64(dst.Height())
		}
		if target.Y < 0 {
			target.Y = 0
		} else if target.Y > float64(dst.Height()) {
			target.Y = float64(dst.Height())
		}
	} else {
		m := direction.Y / direction.X
		b := orgin.Y - m*orgin.X
		if orgin.X < 0 {
			orgin.X = 0
			orgin.Y = b
		} else if orgin.X > float64(dst.Width()) {
			orgin.X = float64(dst.Width())
			orgin.Y = m*orgin.X + b
		}
		if target.X < 0 {
			target.X = 0
			target.Y = b
		} else if target.X > float64(dst.Width()) {
			target.X = float64(dst.Width())
			target.Y = m*target.X + b
		}
	}
	dst.DrawLine(l.Origin.X, float64(dst.Height())-l.Origin.Y, l.Target.X, float64(dst.Height())-l.Target.Y)
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
	angle1 := math.Atan2(invertedOrigin.Y-invertedPerpendicular.Y, invertedOrigin.X-invertedPerpendicular.X)
	angle2 := math.Atan2(invertedTarget.Y-invertedPerpendicular.Y, invertedTarget.X-invertedPerpendicular.X)
	return &Arc{
		Origin: Center(invertedPerpendicular, c.Origin),
		Radius: invertedPerpendicular.Distance(c.Origin) / 2,
		Angle1: angle1,
		Angle2: angle2,
	}
}

package geo

import (
	"image/color"
	"math"

	"github.com/deeean/go-vector/vector2"
	"github.com/fogleman/gg"
)

const threshold = 1e-9

type Line struct {
	Origin *vector2.Vector2
	Target *vector2.Vector2
}

type Ray struct {
	Origin    *vector2.Vector2
	Direction *vector2.Vector2
}

type InvertedLine struct {
	Ray1 *Ray
	Ray2 *Ray
}

type Arc struct {
	Origin *vector2.Vector2
	Radius float64
	Angle1 float64
	Angle2 float64
}

type Circle struct {
	Origin *vector2.Vector2
	Radius float64
}

type Geometry interface {
	Draw(dst *gg.Context)
	Invert(*Circle) Geometry
}

type Drawable interface {
	Draw(g Geometry)
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
	g.Draw(c.Context)
}

func NewCanvas(width, height int) *canvas {
	c := &canvas{
		gg.NewContext(width, height),
	}
	c.SetLineWidth(2)
	return c
}

func (l *Line) Draw(dst *gg.Context) {
	// Draw the line
	clr := color.RGBA{0, 0, 255, 255} // Blue color for the line
	dst.SetColor(clr)
	dst.DrawLine(l.Origin.X, float64(dst.Height())-l.Origin.Y, l.Target.X, float64(dst.Height())-l.Target.Y)
	dst.Stroke()
}

func (r *Ray) Draw(dst *gg.Context) {
	// Draw the ray
	clr := color.RGBA{255, 0, 255, 255} // Magenta color for the ray
	dst.SetColor(clr)
	rayEnd := r.Origin.Add(r.Direction.MulScalar(float64(dst.Height()) + float64(dst.Width())))
	dst.DrawLine(r.Origin.X, float64(dst.Height())-r.Origin.Y, rayEnd.X, float64(dst.Height())-rayEnd.Y)
	dst.Stroke()
}

func (l *InvertedLine) Draw(dst *gg.Context) {
	l.Ray1.Draw(dst)
	l.Ray2.Draw(dst)
}

func (a *Arc) Draw(dst *gg.Context) {
	// Draw the curve
	clr := color.RGBA{255, 0, 0, 255} // Red color for the curve

	dst.SetColor(clr)
	dst.DrawArc(a.Origin.X, float64(dst.Height())-a.Origin.Y, a.Radius, -a.Angle1, -a.Angle2)
	dst.Stroke()
}

func (c *Circle) Draw(dst *gg.Context) {
	// Draw the circle
	clr := color.RGBA{0, 255, 0, 255} // Green color for the circle

	dst.SetColor(clr)
	dst.DrawCircle(c.Origin.X, float64(dst.Height())-c.Origin.Y, c.Radius)
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
		angleCircle := math.Atan2(c.Origin.Y-a.Origin.Y, c.Origin.X-a.Origin.X)
		var lowerBound, upperBound float64
		if angle1 < angle2 {
			lowerBound = math.Ceil((angle1 - angleCircle) / (2 * math.Pi))
			upperBound = math.Floor((angle2 - angleCircle) / (2 * math.Pi))
		} else {
			lowerBound = math.Ceil((angle2 - angleCircle) / (2 * math.Pi))
			upperBound = math.Floor((angle1 - angleCircle) / (2 * math.Pi))
		}
		if lowerBound <= upperBound {
			return &InvertedLine{
				Ray1: &Ray{
					Origin:    invertedPoint1,
					Direction: invertedPoint1.Sub(invertedPoint2).Normalize(),
				},
				Ray2: &Ray{
					Origin:    invertedPoint2,
					Direction: invertedPoint2.Sub(invertedPoint1).Normalize(),
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

func (c *Circle) Invert(oc *Circle) Geometry {
	if math.Abs(c.Origin.Distance(oc.Origin)-c.Radius) < threshold {
		direction := c.Origin.Sub(oc.Origin)
		perpendicular := vector2.New(-direction.Y, direction.X)
		return &InvertedLine{
			Ray1: &Ray{
				Origin:    c.Origin,
				Direction: perpendicular,
			},
			Ray2: &Ray{
				Origin:    c.Origin,
				Direction: perpendicular.MulScalar(-1),
			},
		}
	}
	return c
}

func (l *InvertedLine) Invert(c *Circle) Geometry {
	return l
}

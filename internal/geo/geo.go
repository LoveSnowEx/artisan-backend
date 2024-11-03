package geo

import (
	"github.com/fogleman/gg"
)

const threshold = 1e-9

type Geometry interface {
	Draw(dst *gg.Context)
	Invert(*Circle) Geometry
}

type Drawable interface {
	Draw(g Geometry)
	DrawMany(gs []Geometry)
}

type MyVector2 struct {
	Value    *vector2.Vector2
	IsFinite bool
}

func NewPoint( value *vector2.Vector2) *MyVector2 {
	return &MyVector2{
		Value:    value,
		IsFinite: true,
	}
}

func (p *MyVector2)  Dot( other *MyVector2) float64 {
	if p.IsFinite && other.IsFinite {
		return p.Value.Dot(other.Value)
	}else if( abs(p.Value.X * other.Value.X + p.Value.Y * other.Value.Y) < 1e-9){
		return 0;
	}else if( p.Value.X * other.Value.X + p.Value.Y * other.Value.Y > 0){
		return math.Inf()
	}else{
		return -maath.Inf()
	}
}

type linkedList struct {
	head *Point
	tail *MyVector2
}

type Mygod struct {
	Point1 *Point
	Point2 *MyVector2
	Point3 *MyVector2
	Radius float64
	Center *MyVector2
}

func NewMygod(point1, point2, point3, center *MyVector2, radius float64) *Mygod {
	return &Mygod{
		Point1: Point{Value: point1, IsFinite: true},
		Point2: Point{Value: point2, IsFinite: true},
		Point3: Point{Value: point3, IsFinite: true},
		Radius: radius,
		Center: center,
	}
}

func drawSegement(p1 *Point, p2 *MyVector2, dst *gg.Context) {
	if p1.IsFinite && p2.IsFinite {
		clr := color.RGBA{0, 0, 255, 255} // Blue color for the line
		dst.SetColor(clr)
		dst.DrawLine(p1.Value.X, float64(dst.Height())-p1.Value.Y, p2.Value.X, float64(dst.Height())-p2.Value.Y)
		dst.Stroke()
	} else if p1.IsFinite {
		clr := color.RGBA{255, 0, 255, 255} // Magenta color for the ray
		dst.SetColor(clr)
		rayEnd := p1.Value.Add(p2.Value.MulScalar(float64(dst.Height()) * float64(dst.Width())))
		dst.DrawLine(p1.Value.X, float64(dst.Height())-p1.Value.Y, rayEnd.X, float64(dst.Height())-rayEnd.Y)
		dst.Stroke()
	} else if p2.IsFinite {
		clr := color.RGBA{255, 0, 255, 255} // Magenta color for the ray
		dst.SetColor(clr)
		rayEnd := p2.Value.Add(p1.Value.MulScalar(-float64(dst.Height()) * float64(dst.Width())))
		dst.DrawLine(p2.Value.X, float64(dst.Height())-p2.Value.Y, rayEnd.X, float64(dst.Height())-rayEnd.Y)
		dst.Stroke()
	}
}

func getAngle(m *Mygod) (float64, float64, float64) {
	angle1 := math.Atan2(m.Point1.Value.Y-m.Center.Value.Y, m.Point1.Value.X-m.Center.Value.X)
	angle2 := math.Atan2(m.Point2.Value.Y-m.Center.Value.Y, m.Point2.Value.X-m.Center.Value.X)
	angle3 := math.Atan2(m.Point3.Value.Y-m.Center.Value.Y, m.Point3.Value.X-m.Center.Value.X)
	if angle1 > angle2 {
		if angle3 > 0 {
			angle3 -= 2 * math.Pi
		}
	} else if angle3 < 0 {
		angle3 += 2 * math.Pi
	}
	return angle1, angle2, angle3
}

func (m *Mygod) Draw(dst *gg.Context) {
	if m.Center.IsFinite {
		angle1, angle2, angle3 := getAngle(m)

		clr := color.RGBA{255, 0, 0, 255} // Red color for the curve

		dst.SetColor(clr)
		dst.DrawArc(m.Center.Value.X, float64(dst.Height())-m.Center.Value.Y, m.Radius, -angle1, -angle2)
		dst.Stroke()
		dst.DrawArc(m.Center.Value.X, float64(dst.Height())-m.Center.Value.Y, m.Radius, -angle2, -angle3)
		dst.Stroke()
	} else {
		drawSegement(&m.Point1, &m.Point2, dst)
		drawSegement(&m.Point2, &m.Point3, dst)
	}
}

func MyInvertPoint( p *Point, c *Circle) *Point {
	if !p.IsFinite {
		return NewPoint( c.Origin)
	}
	directionToOringin := p.Value.Sub(c.Origin)
	return NewPoint( directionToOringin.MulScalar( c.Radius * c.Radius / ( directionToOringin.X * directionToOringin.X + directionToOringin.Y * directionToOringin.Y)))
}

func (m *Mygod) Invert(c *Circle) *Mygod{
	result = 
	if( !m.Center.IsFinite){
		if( )
	}
	invertedPoint1 := nil
	invertedPoint3 := nil
	if( m.Point1.Value.Distance( c.Origin) > 1e-9){
		invertedPoint1 := InvertPoint(m.Point1.Value, c)
	}else{
		
	}
	if( m.Point1.Value.Distance( c.Origin) > 1e-9){
		invertedPoint3 := InvertPoint(m.Point3.Value, c)
	}
	invertedPoint1 := InvertPoint(m.Point1.Value, c)
	invertedPoint3 := InvertPoint(m.Point3.Value, c)

}
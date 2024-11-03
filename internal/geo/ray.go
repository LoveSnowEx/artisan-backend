package geo

import (
	"image/color"

	"github.com/deeean/go-vector/vector2"
	"github.com/fogleman/gg"
)

type Ray struct {
	Origin    *vector2.Vector2
	Direction *vector2.Vector2
}

func (r *Ray) Draw(dst *gg.Context) {
	// Draw the ray
	clr := color.RGBA{255, 0, 255, 255} // Magenta color for the ray
	dst.SetColor(clr)
	rayEnd := r.Origin.Add(r.Direction.MulScalar(float64(dst.Height()) + float64(dst.Width())))
	dst.DrawLine(r.Origin.X, float64(dst.Height())-r.Origin.Y, rayEnd.X, float64(dst.Height())-rayEnd.Y)
	dst.Stroke()
}

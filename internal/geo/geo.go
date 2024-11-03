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

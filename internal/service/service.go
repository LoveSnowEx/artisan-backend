package service

import (
	"artisan-backend/internal/geo"
	"bytes"
	"image"
	"strconv"
	"sync"
	"time"

	"github.com/deeean/go-vector/vector2"
	"github.com/gofiber/fiber/v3"
)

const width, height = 1920, 1080

var (
	_               Service = (*service)(nil)
	geometries              = make([]geo.Geometry, 0, 10000)
	inverted                = false
	currentPosition *geo.Ray

	bufPool = sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(nil)
		},
	}
)

type Service interface {
}

type service struct {
	buf       *bytes.Buffer
	canvas    *geo.Canvas
	circulars Circulars
	mu        *sync.Mutex
	bufMu     *sync.Mutex
}

func New() *service {
	srv := &service{
		buf:       bytes.NewBuffer(nil),
		canvas:    geo.NewCanvas(width, height),
		circulars: NewCirculars(),
		mu:        &sync.Mutex{},
		bufMu:     &sync.Mutex{},
	}
	srv.circulars.Enqueue(NewCircular())
	currentPosition = geo.NewRay(
		vector2.New(float64(width)/2, float64(height)/2),
		vector2.New(1, 0),
	)
	go srv.Run()
	return srv
}

func (s *service) Run() {
	t := time.NewTicker(1 * time.Millisecond)
	for range t.C {
		s.mu.Lock()
		for _, c := range s.circulars {
			inst, _ := c.First()
			if inst == nil {
				continue
			}
			inst.Execute()
			c.Rotate()
		}
		s.mu.Unlock()
		s.Draw()
	}
}

func (s *service) Draw() {
	s.canvas.SetColor(image.Transparent)
	s.canvas.Clear()
	circle := &geo.Circle{
		Origin: currentPosition.Origin,
		Radius: vector2.New(width, height).Magnitude() / 10,
	}
	for _, g := range geometries {
		if inverted {
			s.canvas.Draw(g.Invert(circle))
		} else {
			s.canvas.Draw(g)
		}
	}
	s.canvas.Draw(circle)
	oldBuf := s.buf
	tmpBuf := bufPool.Get().(*bytes.Buffer)
	s.canvas.EncodePNG(tmpBuf)
	s.bufMu.Lock()
	s.buf = tmpBuf
	s.bufMu.Unlock()
	go func() {
		oldBuf.Reset()
		bufPool.Put(oldBuf)
	}()

	// go func() {
	// 	s.canvas.SavePNG("image.png")
	// }()
}

func (s *service) Forward(ctx fiber.Ctx) error {
	px := ctx.Params("px")
	val, err := strconv.Atoi(px)
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	c := s.circulars[0]
	c.AddLast(NewForwardInstruction(val))
	return ctx.SendStatus(fiber.StatusOK)
}

func (s *service) Left(ctx fiber.Ctx) error {
	deg := ctx.Params("deg")
	val, err := strconv.Atoi(deg)
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	c := s.circulars[0]
	c.AddLast(NewLeftInstruction(val))
	return ctx.SendStatus(fiber.StatusOK)
}

func (s *service) Right(ctx fiber.Ctx) error {
	deg := ctx.Params("deg")
	val, err := strconv.Atoi(deg)
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	c := s.circulars[0]
	c.AddLast(NewRightInstruction(val))
	return ctx.SendStatus(fiber.StatusOK)
}

func (s *service) Reverse(ctx fiber.Ctx) error {
	c := s.circulars[0]
	c.AddLast(NewReverseInstruction())
	return ctx.SendStatus(fiber.StatusOK)
}

func (s *service) Inversion(ctx fiber.Ctx) error {
	c := s.circulars[0]
	c.AddLast(NewInversionInstruction())
	return ctx.SendStatus(fiber.StatusOK)
}

func (s *service) Circulars(ctx fiber.Ctx) error {
	s.mu.Lock()
	values := s.circulars.Values()
	s.mu.Unlock()
	ctx.JSON(values)
	return ctx.SendStatus(fiber.StatusOK)
}

func (s *service) Img(ctx fiber.Ctx) error {
	s.bufMu.Lock()
	defer s.bufMu.Unlock()
	ctx.Set(fiber.HeaderContentType, "image/png")
	ctx.Write(s.buf.Bytes())
	return ctx.SendStatus(fiber.StatusOK)
}

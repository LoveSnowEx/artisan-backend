package service

import (
	"artisan-backend/internal/model"
)

var (
	_ Service = (*service)(nil)
)

type Instruction string

type Service interface {
}

type service struct {
	circulars model.Circulars
}

func New() *service {
	srv := &service{
		model.NewCirculars(),
	}
	srv.circulars.Enqueue(model.NewCircular())
	return srv
}

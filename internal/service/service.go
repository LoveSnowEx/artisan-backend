package service

var (
	_ Service = (*service)(nil)
)

type Service interface {
}

type service struct {
}

func New() *service {
	return &service{}
}

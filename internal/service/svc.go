package service

import (
	"cro_test/internal/domain"
)

func New(repo domain.Repositorier) domain.Servicer {
	return Service{
		Repositorier: repo,
	}
}

type Service struct {
	domain.Repositorier
}

func (svc Service) repo() domain.Repositorier {
	return svc.Repositorier
}

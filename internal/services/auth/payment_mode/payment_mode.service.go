package payment_mode

import (
	"context"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) FindAll(ctx context.Context) ([]PaymentMode, error) {
	return s.repo.FindAll(ctx)
}

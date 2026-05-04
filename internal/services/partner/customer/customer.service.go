package customer

import (
	"context"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) Create(ctx context.Context, customer *Customer) error {
	return s.repository.Create(ctx, customer)
}

func (s *Service) Update(ctx context.Context, customer *Customer) error {
	return s.repository.Update(ctx, customer)
}

func (s *Service) Delete(ctx context.Context, id uint) error {
	return s.repository.Delete(ctx, id)
}

func (s *Service) FindByID(ctx context.Context, id uint) (*Customer, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *Service) FindAll(ctx context.Context) ([]Customer, error) {
	return s.repository.FindAll(ctx)
}

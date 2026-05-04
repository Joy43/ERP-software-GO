package tax_bracket

import (
	"context"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, tb *TaxBracket) error {
	return s.repo.Create(ctx, tb)
}

func (s *Service) FindAll(ctx context.Context) ([]TaxBracket, error) {
	return s.repo.FindAll(ctx)
}

func (s *Service) FindByID(ctx context.Context, id uint) (*TaxBracket, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *Service) Update(ctx context.Context, tb *TaxBracket) error {
	return s.repo.Update(ctx, tb)
}

func (s *Service) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

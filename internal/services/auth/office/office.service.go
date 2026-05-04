package office

import (
	"context"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, o *Office) error {
	return s.repo.Create(ctx, o)
}

func (s *Service) FindAll(ctx context.Context) ([]Office, error) {
	return s.repo.FindAll(ctx)
}

func (s *Service) FindByID(ctx context.Context, id uint) (*Office, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *Service) Update(ctx context.Context, o *Office) error {
	return s.repo.Update(ctx, o)
}

func (s *Service) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

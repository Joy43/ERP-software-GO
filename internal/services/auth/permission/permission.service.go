package permission

import (
	"context"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListPermissions(ctx context.Context) ([]Permission, error) {
	return s.repo.FindAll(ctx)
}

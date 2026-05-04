package partner_sub_group

import (
	"context"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, psg *PartnerSubGroup) error {
	return s.repo.Create(ctx, psg)
}

func (s *Service) FindAll(ctx context.Context) ([]PartnerSubGroup, error) {
	return s.repo.FindAll(ctx)
}

func (s *Service) FindByID(ctx context.Context, id uint) (*PartnerSubGroup, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *Service) FindByGroupID(ctx context.Context, groupID uint) ([]PartnerSubGroup, error) {
	return s.repo.FindByGroupID(ctx, groupID)
}

func (s *Service) Update(ctx context.Context, psg *PartnerSubGroup) error {
	return s.repo.Update(ctx, psg)
}

func (s *Service) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

package responsibility_transfer

import (
	"context"
	"time"
)

type Service interface {
	CreateTransfer(ctx context.Context, req CreateTransferRequest) (*ResponsibilityTransfer, error)
	ListTransfers(ctx context.Context, filter TransferFilter) ([]ResponsibilityTransfer, error)
	ApproveTransfer(ctx context.Context, id uint, approvedByID uint, req ApproveTransferRequest) (*ResponsibilityTransfer, error)
	DeleteTransfer(ctx context.Context, id uint) error
	GetTransfer(ctx context.Context, id uint) (*ResponsibilityTransfer, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateTransfer(ctx context.Context, req CreateTransferRequest) (*ResponsibilityTransfer, error) {
	transfer := &ResponsibilityTransfer{
		FromUserID: req.FromUserID,
		ToUserID:   req.ToUserID,
		FromDate:   req.FromDate,
		ToDate:     req.ToDate,
		Remarks:    req.Remarks,
		Status:     "Pending",
	}

	if err := s.repo.Create(ctx, transfer); err != nil {
		return nil, err
	}
	return s.repo.FindByIDWithUsers(ctx, transfer.ID)
}

func (s *service) ListTransfers(ctx context.Context, filter TransferFilter) ([]ResponsibilityTransfer, error) {
	return s.repo.FindAllWithUsers(ctx, filter)
}

func (s *service) ApproveTransfer(ctx context.Context, id uint, approvedByID uint, req ApproveTransferRequest) (*ResponsibilityTransfer, error) {
	transfer, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	transfer.Status = req.Status
	transfer.ApprovedByID = &approvedByID
	transfer.ApprovedAt = &now

	if err := s.repo.Update(ctx, transfer); err != nil {
		return nil, err
	}
	return s.repo.FindByIDWithUsers(ctx, transfer.ID)
}

func (s *service) DeleteTransfer(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) GetTransfer(ctx context.Context, id uint) (*ResponsibilityTransfer, error) {
	return s.repo.FindByIDWithUsers(ctx, id)
}

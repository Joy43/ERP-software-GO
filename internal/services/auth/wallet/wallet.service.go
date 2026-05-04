package wallet

import (
	"context"
)

type Service interface {
	Create(ctx context.Context, req CreateWalletRequest) (*Wallet, error)
	List(ctx context.Context) ([]Wallet, error)
	GetByID(ctx context.Context, id uint) (*Wallet, error)
	Update(ctx context.Context, id uint, req UpdateWalletRequest) (*Wallet, error)
	Delete(ctx context.Context, id uint) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, req CreateWalletRequest) (*Wallet, error) {
	wallet := &Wallet{
		Name:              req.Name,
		CommissionPercent: req.CommissionPercent,
		BankAccount:       req.BankAccount,
		ReferenceCode:     req.ReferenceCode,
		IsRounding:        req.IsRounding,
		IsCoupon:          req.IsCoupon,
		IsWalletCharge:    req.IsWalletCharge,
	}
	if req.IsActive != nil {
		wallet.IsActive = *req.IsActive
	} else {
		wallet.IsActive = true
	}

	if err := s.repo.Create(ctx, wallet); err != nil {
		return nil, err
	}
	return wallet, nil
}

func (s *service) List(ctx context.Context) ([]Wallet, error) {
	return s.repo.FindAll(ctx)
}

func (s *service) GetByID(ctx context.Context, id uint) (*Wallet, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) Update(ctx context.Context, id uint, req UpdateWalletRequest) (*Wallet, error) {
	wallet, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		wallet.Name = req.Name
	}
	if req.CommissionPercent != nil {
		wallet.CommissionPercent = *req.CommissionPercent
	}
	if req.BankAccount != nil {
		wallet.BankAccount = *req.BankAccount
	}
	if req.ReferenceCode != nil {
		wallet.ReferenceCode = *req.ReferenceCode
	}
	if req.IsRounding != nil {
		wallet.IsRounding = *req.IsRounding
	}
	if req.IsCoupon != nil {
		wallet.IsCoupon = *req.IsCoupon
	}
	if req.IsWalletCharge != nil {
		wallet.IsWalletCharge = *req.IsWalletCharge
	}
	if req.IsActive != nil {
		wallet.IsActive = *req.IsActive
	}

	if err := s.repo.Update(ctx, wallet); err != nil {
		return nil, err
	}
	return wallet, nil
}

func (s *service) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

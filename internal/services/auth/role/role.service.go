package role

import (
	"context"
	"fmt"
	"strings"

	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/auth/permission"
)

type CreateRoleRequest struct {
	Name          string `json:"name" binding:"required"`
	Description   string `json:"description"`
	IsActive      *bool  `json:"is_active"`
	PermissionIDs []uint `json:"permission_ids"`
}

type UpdateRoleRequest struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	IsActive      *bool  `json:"is_active"`
	PermissionIDs []uint `json:"permission_ids"`
}

type Service struct {
	repo           Repository
	permissionRepo permission.Repository
}

func NewService(repo Repository, permissionRepo permission.Repository) *Service {
	return &Service{repo: repo, permissionRepo: permissionRepo}
}

func (s *Service) CreateRole(ctx context.Context, req CreateRoleRequest) (*Role, error) {
	var role *Role
	err := s.repo.Transaction(ctx, func(txRepo Repository) error {
		slug := strings.ToLower(strings.ReplaceAll(req.Name, " ", "-"))
		role = &Role{
			Name:        req.Name,
			Slug:        slug,
			Description: req.Description,
		}

		if req.IsActive != nil {
			role.IsActive = *req.IsActive
		} else {
			role.IsActive = true
		}

		if err := txRepo.Create(ctx, role); err != nil {
			return err
		}

		if len(req.PermissionIDs) > 0 {
			if err := txRepo.ReplacePermissions(ctx, role, req.PermissionIDs); err != nil {
				return fmt.Errorf("assign permissions: %w", err)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.repo.FindByID(ctx, role.ID)
}

func (s *Service) ListRoles(ctx context.Context) ([]Role, error) {
	return s.repo.FindAll(ctx)
}

func (s *Service) GetRole(ctx context.Context, id uint) (*Role, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *Service) UpdateRole(ctx context.Context, id uint, req UpdateRoleRequest) (*Role, error) {
	err := s.repo.Transaction(ctx, func(txRepo Repository) error {
		role, err := txRepo.FindByID(ctx, id)
		if err != nil {
			return err
		}

		if req.Name != "" {
			role.Name = req.Name
			role.Slug = strings.ToLower(strings.ReplaceAll(req.Name, " ", "-"))
		}
		if req.Description != "" {
			role.Description = req.Description
		}
		if req.IsActive != nil {
			role.IsActive = *req.IsActive
		}

		if err := txRepo.Update(ctx, role); err != nil {
			return err
		}

		if req.PermissionIDs != nil {
			if err := txRepo.ReplacePermissions(ctx, role, req.PermissionIDs); err != nil {
				return fmt.Errorf("update permissions: %w", err)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.repo.FindByID(ctx, id)
}

func (s *Service) DeleteRole(ctx context.Context, id uint) error {
	return s.repo.Transaction(ctx, func(txRepo Repository) error {
		return txRepo.Delete(ctx, id)
	})
}

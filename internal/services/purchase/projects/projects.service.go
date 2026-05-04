package projects

import (
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

var (
	ErrProjectCodeAlreadyExists = errors.New("project code already exists")
	ErrProjectNotFound          = errors.New("project not found")
	ErrInvalidProjectStatus     = errors.New("invalid project status")
	ErrCannotDeleteCompletedProject = errors.New("cannot delete a completed project")
	ErrCannotDeleteCancelledProject = errors.New("cannot delete a cancelled project")
	ErrInvalidDateRange         = errors.New("end date must be after start date")
)

type Service struct {
	repository Repository
}

// =============================
// SERVICE LAYER DTOs (Internal use)
// =============================

// CreateProjectCmd is the internal command for creating a project
type CreateProjectCmd struct {
	ProjectName string
	Description *string
	ProjectCode string
	StartDate   *time.Time
	EndDate     *time.Time
	Budget      *float64
	Status      *ProjectStatus
	IsActive    *bool
	ManagerID   *uint
	OfficeID    uint
}

// UpdateProjectCmd is the internal command for updating a project
type UpdateProjectCmd struct {
	ProjectName *string
	Description *string
	ProjectCode *string
	StartDate   *time.Time
	EndDate     *time.Time
	Budget      *float64
	Status      *ProjectStatus
	IsActive    *bool
	ManagerID   *uint
	OfficeID    *uint
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

// CreateProject creates a new project
func (s *Service) CreateProject(ctx context.Context, cmd CreateProjectCmd) (*Project, error) {
	// Validate project code
	projectCode := strings.TrimSpace(cmd.ProjectCode)
	existingProject, err := s.repository.GetByProjectCode(ctx, projectCode)
	if err == nil && existingProject != nil {
		return nil, ErrProjectCodeAlreadyExists
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Validate dates
	if cmd.StartDate != nil && cmd.EndDate != nil {
		if cmd.EndDate.Before(*cmd.StartDate) {
			return nil, ErrInvalidDateRange
		}
	}

	// Set default values
	status := StatusPlanning
	if cmd.Status != nil {
		status = *cmd.Status
	}

	isActive := true
	if cmd.IsActive != nil {
		isActive = *cmd.IsActive
	}

	project := &Project{
		ProjectName: cmd.ProjectName,
		Description: cmd.Description,
		ProjectCode: projectCode,
		StartDate:   cmd.StartDate,
		EndDate:     cmd.EndDate,
		Budget:      cmd.Budget,
		Status:      status,
		IsActive:    isActive,
		ManagerID:   cmd.ManagerID,
		OfficeID:    cmd.OfficeID,
	}

	if err := s.repository.Create(ctx, project); err != nil {
		return nil, err
	}

	return project, nil
}

// GetProjectByID retrieves a project by its ID
func (s *Service) GetProjectByID(ctx context.Context, id uint) (*Project, error) {
	project, err := s.repository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProjectNotFound
		}
		return nil, err
	}
	return project, nil
}

// GetAllProjects retrieves all projects
func (s *Service) GetAllProjects(ctx context.Context) ([]Project, error) {
	return s.repository.GetAll(ctx)
}

// GetProjectsPaginated retrieves projects with pagination and filters
func (s *Service) GetProjectsPaginated(ctx context.Context, page, limit int, filters map[string]interface{}) ([]Project, int64, error) {
	return s.repository.GetAllPaginated(ctx, page, limit, filters)
}

// UpdateProject updates an existing project
func (s *Service) UpdateProject(ctx context.Context, id uint, cmd UpdateProjectCmd) (*Project, error) {
	// Get existing project
	project, err := s.repository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProjectNotFound
		}
		return nil, err
	}

	// Validate project code if being updated
	if cmd.ProjectCode != nil {
		newCode := strings.TrimSpace(*cmd.ProjectCode)
		if newCode != project.ProjectCode {
			existingProject, err := s.repository.GetByProjectCode(ctx, newCode)
			if err == nil && existingProject != nil {
				return nil, ErrProjectCodeAlreadyExists
			}
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, err
			}
			project.ProjectCode = newCode
		}
	}

	// Validate dates if being updated
	startDate := project.StartDate
	endDate := project.EndDate
	if cmd.StartDate != nil {
		startDate = cmd.StartDate
	}
	if cmd.EndDate != nil {
		endDate = cmd.EndDate
	}
	if startDate != nil && endDate != nil {
		if endDate.Before(*startDate) {
			return nil, ErrInvalidDateRange
		}
	}

	// Update fields
	if cmd.ProjectName != nil {
		project.ProjectName = *cmd.ProjectName
	}
	if cmd.Description != nil {
		project.Description = cmd.Description
	}
	if cmd.StartDate != nil {
		project.StartDate = cmd.StartDate
	}
	if cmd.EndDate != nil {
		project.EndDate = cmd.EndDate
	}
	if cmd.Budget != nil {
		project.Budget = cmd.Budget
	}
	if cmd.Status != nil {
		project.Status = *cmd.Status
	}
	if cmd.IsActive != nil {
		project.IsActive = *cmd.IsActive
	}
	if cmd.ManagerID != nil {
		project.ManagerID = cmd.ManagerID
	}
	if cmd.OfficeID != nil {
		project.OfficeID = *cmd.OfficeID
	}

	if err := s.repository.Update(ctx, project); err != nil {
		return nil, err
	}

	return project, nil
}

// DeleteProject deletes a project
func (s *Service) DeleteProject(ctx context.Context, id uint) error {
	// Get project
	project, err := s.repository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrProjectNotFound
		}
		return err
	}

	// Prevent deletion of completed or cancelled projects
	if project.Status == StatusCompleted {
		return ErrCannotDeleteCompletedProject
	}
	if project.Status == StatusCancelled {
		return ErrCannotDeleteCancelledProject
	}

	return s.repository.Delete(ctx, id)
}

// GetProjectsByOffice retrieves all projects for a specific office
func (s *Service) GetProjectsByOffice(ctx context.Context, officeID uint) ([]Project, error) {
	return s.repository.GetByOfficeID(ctx, officeID)
}

// GetProjectsByManager retrieves all projects managed by a specific user
func (s *Service) GetProjectsByManager(ctx context.Context, managerID uint) ([]Project, error) {
	return s.repository.GetByManagerID(ctx, managerID)
}

// GetProjectsByStatus retrieves all projects with a specific status
func (s *Service) GetProjectsByStatus(ctx context.Context, status ProjectStatus) ([]Project, error) {
	return s.repository.GetByStatus(ctx, status)
}

// GetActiveProjects retrieves all active projects
func (s *Service) GetActiveProjects(ctx context.Context) ([]Project, error) {
	return s.repository.GetActiveProjects(ctx)
}

// ChangeProjectStatus changes the status of a project
func (s *Service) ChangeProjectStatus(ctx context.Context, id uint, status ProjectStatus) (*Project, error) {
	project, err := s.repository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProjectNotFound
		}
		return nil, err
	}

	if err := s.repository.UpdateStatus(ctx, id, status); err != nil {
		return nil, err
	}

	// Update the local project object for response
	project.Status = status
	return project, nil
}

// DeactivateProject deactivates a project
func (s *Service) DeactivateProject(ctx context.Context, id uint) (*Project, error) {
	project, err := s.repository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProjectNotFound
		}
		return nil, err
	}

	project.IsActive = false
	if err := s.repository.Update(ctx, project); err != nil {
		return nil, err
	}

	return project, nil
}

// ActivateProject activates a project
func (s *Service) ActivateProject(ctx context.Context, id uint) (*Project, error) {
	project, err := s.repository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProjectNotFound
		}
		return nil, err
	}

	project.IsActive = true
	if err := s.repository.Update(ctx, project); err != nil {
		return nil, err
	}

	return project, nil
}
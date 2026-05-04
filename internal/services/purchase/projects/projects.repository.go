package projects

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, project *Project) error
	GetByID(ctx context.Context, id uint) (*Project, error)
	GetByProjectCode(ctx context.Context, code string) (*Project, error)
	GetAll(ctx context.Context) ([]Project, error)
	GetAllPaginated(ctx context.Context, page, limit int, filters map[string]interface{}) ([]Project, int64, error)
	Update(ctx context.Context, project *Project) error
	UpdateStatus(ctx context.Context, id uint, status ProjectStatus) error
	Delete(ctx context.Context, id uint) error
	GetByOfficeID(ctx context.Context, officeID uint) ([]Project, error)
	GetByManagerID(ctx context.Context, managerID uint) ([]Project, error)
	GetByStatus(ctx context.Context, status ProjectStatus) ([]Project, error)
	GetActiveProjects(ctx context.Context) ([]Project, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

//---------- Create creates a new project----------
func (r *repository) Create(ctx context.Context, project *Project) error {
	if err := r.db.WithContext(ctx).Create(project).Error; err != nil {
		return fmt.Errorf("create project: %w", err)
	}
	return nil
}

//------------ GetByID retrieves a project by its ID------------
func (r *repository) GetByID(ctx context.Context, id uint) (*Project, error) {
	var project Project
	if err := r.db.WithContext(ctx).
		Preload("Manager").
		Preload("Office").
		First(&project, id).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

// GetByProjectCode retrieves a project by its project code
func (r *repository) GetByProjectCode(ctx context.Context, code string) (*Project, error) {
	var project Project
	if err := r.db.WithContext(ctx).
		Where("project_code = ?", code).
		Preload("Manager").
		Preload("Office").
		First(&project).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

// GetAll retrieves all projects
func (r *repository) GetAll(ctx context.Context) ([]Project, error) {
	var projects []Project
	if err := r.db.WithContext(ctx).
		Preload("Manager").
		Preload("Office").
		Where("is_active = ?", true).
		Order("created_at DESC").
		Find(&projects).Error; err != nil {
		return nil, fmt.Errorf("get all projects: %w", err)
	}
	return projects, nil
}

// GetAllPaginated retrieves all projects with pagination and filters
func (r *repository) GetAllPaginated(ctx context.Context, page, limit int, filters map[string]interface{}) ([]Project, int64, error) {

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	// Build query
	query := r.db.WithContext(ctx)

	// Apply filters
	if status, ok := filters["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}
	if officeID, ok := filters["office_id"].(uint); ok && officeID != 0 {
		query = query.Where("office_id = ?", officeID)
	}
	if managerID, ok := filters["manager_id"].(uint); ok && managerID != 0 {
		query = query.Where("manager_id = ?", managerID)
	}
	if isActive, ok := filters["is_active"].(bool); ok {
		query = query.Where("is_active = ?", isActive)
	}

	// Count total records
	var total int64
	if err := query.Model(&Project{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count projects: %w", err)
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Apply pagination
	var projects []Project
	if err := query.
		Preload("Manager").
		Preload("Office").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&projects).Error; err != nil {
		return nil, 0, fmt.Errorf("get paginated projects: %w", err)
	}

	return projects, total, nil
}

// --------- Update updates an existing project----------
func (r *repository) Update(ctx context.Context, project *Project) error {
	// Use Updates with map to only update specific fields, avoiding relation issues
	updates := map[string]interface{}{
		"project_name": project.ProjectName,
		"description":  project.Description,
		"project_code": project.ProjectCode,
		"start_date":   project.StartDate,
		"end_date":     project.EndDate,
		"budget":       project.Budget,
		"status":       project.Status,
		"is_active":    project.IsActive,
		"manager_id":   project.ManagerID,
		"office_id":    project.OfficeID,
	}
	
	if err := r.db.WithContext(ctx).Model(&Project{}).Where("id = ?", project.ID).Updates(updates).Error; err != nil {
		return fmt.Errorf("update project: %w", err)
	}
	return nil
}

// UpdateStatus updates only the status field of a project
func (r *repository) UpdateStatus(ctx context.Context, id uint, status ProjectStatus) error {
	if err := r.db.WithContext(ctx).Model(&Project{}).Where("id = ?", id).Update("status", status).Error; err != nil {
		return fmt.Errorf("update project status: %w", err)
	}
	return nil
}

//--------------  Delete deletes a project (soft delete)-------------
func (r *repository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&Project{}, id).Error; err != nil {
		return fmt.Errorf("delete project: %w", err)
	}
	return nil
}

//----------- GetByOfficeID retrieves all projects by office ID-------------
func (r *repository) GetByOfficeID(ctx context.Context, officeID uint) ([]Project, error) {
	var projects []Project
	if err := r.db.WithContext(ctx).
		Where("office_id = ? AND is_active = ?", officeID, true).
		Preload("Manager").
		Preload("Office").
		Order("created_at DESC").
		Find(&projects).Error; err != nil {
		return nil, fmt.Errorf("get projects by office: %w", err)
	}
	return projects, nil
}

//--------- GetByManagerID retrieves all projects managed by a specific user-----------
func (r *repository) GetByManagerID(ctx context.Context, managerID uint) ([]Project, error) {
	var projects []Project
	if err := r.db.WithContext(ctx).
		Where("manager_id = ? AND is_active = ?", managerID, true).
		Preload("Manager").
		Preload("Office").
		Order("created_at DESC").
		Find(&projects).Error; err != nil {
		return nil, fmt.Errorf("get projects by manager: %w", err)
	}
	return projects, nil
}

//--------- GetByStatus retrieves all projects with a specific status-----------
func (r *repository) GetByStatus(ctx context.Context, status ProjectStatus) ([]Project, error) {
	var projects []Project
	if err := r.db.WithContext(ctx).
		Where("status = ? AND is_active = ?", status, true).
		Preload("Manager").
		Preload("Office").
		Order("created_at DESC").
		Find(&projects).Error; err != nil {
		return nil, fmt.Errorf("get projects by status: %w", err)
	}
	return projects, nil
}

// GetActiveProjects retrieves all active projects
func (r *repository) GetActiveProjects(ctx context.Context) ([]Project, error) {
	var projects []Project
	if err := r.db.WithContext(ctx).
		Where("is_active = ? AND status IN ?", true, []ProjectStatus{StatusPlanning, StatusActive, StatusOnHold}).
		Preload("Manager").
		Preload("Office").
		Order("created_at DESC").
		Find(&projects).Error; err != nil {
		return nil, fmt.Errorf("get active projects: %w", err)
	}
	return projects, nil
}
package projects

// =============================
// API REQUEST DTOs
// =============================


type CreateProjectRequest struct {
	ProjectName string     `json:"project_name" binding:"required,max=255" example:"New Office Project"`
	Description *string    `json:"description,omitempty" example:"Description of the project"`
	ProjectCode string     `json:"project_code" binding:"required,max=50" example:"PROJ-001"`
	StartDate   *string    `json:"start_date,omitempty" example:"2024-04-28"`
	EndDate     *string    `json:"end_date,omitempty" example:"2024-12-31"`
	Budget      *float64   `json:"budget,omitempty" binding:"omitempty,gt=0" example:"50000.00"`
	Status      *ProjectStatus `json:"status,omitempty" binding:"omitempty,oneof=PLANNING ACTIVE ON_HOLD COMPLETED CANCELLED" example:"PLANNING"`
	IsActive    *bool      `json:"is_active,omitempty" example:"true"`
	ManagerID   *uint      `json:"manager_id,omitempty" example:"1"`
	OfficeID    uint       `json:"office_id" binding:"required,gt=0" example:"1"`
}


type UpdateProjectRequest struct {
	ProjectName *string    `json:"project_name,omitempty" binding:"omitempty,max=255" example:"Updated Project Name"`
	Description *string    `json:"description,omitempty" binding:"omitempty,max=1000" example:"Updated description"`
	ProjectCode *string    `json:"project_code,omitempty" binding:"omitempty,max=50"`
	StartDate   *string    `json:"start_date,omitempty" example:"2024-05-01"`
	EndDate     *string    `json:"end_date,omitempty" example:"2024-11-30"`
	Budget      *float64   `json:"budget,omitempty" binding:"omitempty,gt=0" example:"60000.00"`
	Status      *ProjectStatus `json:"status,omitempty" binding:"omitempty,oneof=PLANNING ACTIVE ON_HOLD COMPLETED CANCELLED"`
	IsActive    *bool      `json:"is_active,omitempty" example:"true"`
	
}


type ChangeStatusRequest struct {
	Status ProjectStatus `json:"status" binding:"required,oneof=PLANNING ACTIVE ON_HOLD COMPLETED CANCELLED"`
}

// =============================
// API PAGINATION DTOs
// =============================

// ListProjectsRequest represents pagination and filter parameters
type ListProjectsRequest struct {
	Page      int    `form:"page,default=1" binding:"min=1"`
	Limit     int    `form:"limit,default=20" binding:"min=1,max=100"`
	Status    *string `form:"status,omitempty"`
	OfficeID  *uint  `form:"office_id,omitempty"`
	ManagerID *uint  `form:"manager_id,omitempty"`
}

// =============================
// API RESPONSE DTOs
// =============================

// ManagerResponse represents manager details in project response
type ManagerResponse struct {
	ID       uint   `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

// OfficeResponse represents office details in project response
type OfficeResponse struct {
	ID        uint   `json:"id"`
	OfficeName string `json:"office_name"`
	Location  string `json:"location"`
}


type ProjectResponse struct {
	ID          uint             `json:"id"`
	ProjectName string           `json:"project_name"`
	Description *string          `json:"description,omitempty"`
	ProjectCode string           `json:"project_code"`
	StartDate   *string          `json:"start_date,omitempty"`
	EndDate     *string          `json:"end_date,omitempty"`
	Budget      *float64         `json:"budget,omitempty"`
	Status      ProjectStatus    `json:"status"`
	IsActive    bool             `json:"is_active"`
	Manager     *ManagerResponse `json:"manager,omitempty"`
	Office      *OfficeResponse  `json:"office,omitempty"`
	CreatedAt   string           `json:"created_at"`
	UpdatedAt   string           `json:"updated_at"`
}


type ListProjectsResponse struct {
	Data       []ProjectResponse `json:"data"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	TotalPages int64             `json:"total_pages"`
}

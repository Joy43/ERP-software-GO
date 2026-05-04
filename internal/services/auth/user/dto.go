package user

import (
	"time"

	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/auth/role"
)

// PaginationParams represents pagination parameters
type PaginationParams struct {
	Page  int `form:"page,default=1" binding:"min=1"`
	Limit int `form:"limit,default=20" binding:"min=1,max=400"`
}

// ListUsersRequest represents the request parameters for listing users with pagination only
type ListUsersRequest struct {
	PaginationParams
}

// ListUsersResponse represents the response for listing users
type ListUsersResponse struct {
	Data       []UserProfileResponse `json:"data,omitempty"`
	Total      int64                 `json:"total"`
	Page       int                   `json:"page"`
	Limit      int                   `json:"limit"`
	TotalPages int64                 `json:"total_pages"`
}

// UserProfileResponse represents user profile in list response
type UserProfileResponse struct {
	ID            uint        `json:"id"`
	Name          string      `json:"name"`
	Email         string      `json:"email"`
	Mobile        string      `json:"mobile"`
	DOB           *time.Time  `json:"dob"`
	Gender        string      `json:"gender"`
	BloodGroup    string      `json:"blood_group"`
	EmployeeID    string      `json:"employee_id"`
	DesignationID *uint       `json:"designation_id"`
	DepartmentID  *uint       `json:"department_id"`
	OfficeID      *uint       `json:"office_id"`
	JoiningDate   *time.Time  `json:"joining_date"`
	BankName      string      `json:"bank_name"`
	AccountNumber string      `json:"account_number"`
	GrossSalary   float64     `json:"gross_salary"`
	PaymentModeID *uint       `json:"payment_mode_id"`
	Roles         []role.Role `json:"roles"`
	CreatedAt     time.Time   `json:"created_at"`
}
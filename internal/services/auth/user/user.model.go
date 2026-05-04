package user

import (
	"time"

	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/auth/role"
)

type User struct {
	ID            uint        `gorm:"primaryKey"`
	Name          string      `gorm:"size:150;not null"`
	Email         string      `gorm:"size:190;uniqueIndex;not null"`
	PasswordHash  string      `gorm:"size:255;not null"`
	Mobile        string      `gorm:"size:20"`
	DOB           *time.Time  `gorm:"type:date"`
	Gender        string      `gorm:"type:enum('Male', 'Female', 'Other')"`
	BloodGroup    string      `gorm:"type:enum('A+', 'A-', 'B+', 'B-', 'AB+', 'AB-', 'O+', 'O-')"`
	EmployeeID    string      `gorm:"size:50;unique"`
	DesignationID *uint       `gorm:"index"`
	DepartmentID  *uint       `gorm:"index"`
	OfficeID      *uint       `gorm:"index"`
	JoiningDate   *time.Time  `gorm:"type:date"`
	BankName      string      `gorm:"size:100"`
	AccountNumber string      `gorm:"size:50"`
	GrossSalary   float64     `gorm:"type:decimal(15,2)"`
	PaymentModeID *uint       `gorm:"index"`
	Roles         []role.Role `gorm:"many2many:user_roles;"`
	IsActive      bool        `gorm:"not null;default:true"`
	CreatedAt     time.Time   `gorm:"autoCreateTime"`
	UpdatedAt     time.Time   `gorm:"autoUpdateTime"`
}

type RefreshToken struct {
	ID        uint       `gorm:"primaryKey"`
	UserID    uint       `gorm:"index;not null"`
	TokenHash string     `gorm:"size:255;uniqueIndex;not null"`
	ExpiresAt time.Time  `gorm:"not null"`
	RevokedAt *time.Time `gorm:"default:null"`
	CreatedAt time.Time  `gorm:"autoCreateTime"`
}

type UserProfile struct {
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
}

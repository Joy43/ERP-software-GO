package location

import "time"

// =============================
// CREATE DTO
// =============================
type CreateLocationRequest struct {
	Code      string       `json:"code" binding:"required"`
	Name      string       `json:"name" binding:"required"`
	Type      LocationType `json:"type" binding:"required"`

	OfficeID  uint  `json:"office_id" binding:"required" `
	ParentID  *uint `json:"parent_id"`

	ManagerID *uint `json:"manager_id"`

	Location  string `json:"location"`
	IsActive  *bool  `json:"is_active"`
}

// =============================
// UPDATE DTO
// =============================
type UpdateLocationRequest struct {
	Code      *string       `json:"code,omitempty" binding:"omitempty,min=2,max=50"`
	Name      *string       `json:"name,omitempty" binding:"omitempty,min=2,max=200"`
	Type      *LocationType `json:"type,omitempty" binding:"omitempty,oneof=warehouse store showroom outlet"`

	OfficeID  *uint `json:"office_id,omitempty" binding:"omitempty,gt=0"`
	ParentID  *uint `json:"parent_id,omitempty"`

	ManagerID *uint `json:"manager_id,omitempty" binding:"omitempty,gt=0"`

	Location  *string `json:"location,omitempty" binding:"omitempty,max=500"`
	IsActive  *bool   `json:"is_active,omitempty"`
}

// =============================
// RESPONSE DTO
// =============================
type LocationResponse struct {
	ID   uint   `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
	Type string `json:"type"`

	OfficeID uint `json:"office_id"`

	ParentID *uint `json:"parent_id"`

	ManagerID *uint `json:"manager_id"`

	Location string `json:"location"`
	IsActive bool   `json:"is_active"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// =============================
//  WITH RELATIONS
// =============================
type LocationWithRelationsResponse struct {
	ID   uint   `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
	Type string `json:"type"`

	OfficeID uint                   `json:"office_id"`
	Office   *SimpleOfficeResponse  `json:"office,omitempty"`

	ParentID *uint            `json:"parent_id"`
	Parent   *SimpleLocation  `json:"parent,omitempty"`

	Children []SimpleLocation `json:"children,omitempty"`

	ManagerID *uint               `json:"manager_id"`
	Manager   *SimpleUserResponse `json:"manager,omitempty"`

	Location string `json:"location"`
	IsActive bool   `json:"is_active"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// =============================
// SMALL EMBEDDED DTOs
// =============================
type SimpleLocation struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type SimpleOfficeResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type SimpleUserResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
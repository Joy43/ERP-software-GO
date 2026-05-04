package inventorytypes
type CreateInventoryTypeRequest struct {
	TypeCode    string `json:"type_code" binding:"required,min=2,max=150"`
	TypeName    string `json:"type_name" binding:"required,min=2,max=100"`
	Description string `json:"description" binding:"omitempty,max=255"`
	IsActive    *bool  `json:"is_active"` 
}

type UpdateInventoryTypeRequest struct {
	TypeCode    *string `json:"type_code" binding:"omitempty,min=2,max=150"`
	TypeName    *string `json:"type_name" binding:"omitempty,min=2,max=100"`
	Description *string `json:"description" binding:"omitempty,max=255"`
	IsActive    *bool   `json:"is_active"`
}

type InventoryTypeResponse struct {
	ID          uint   `json:"id"`
	TypeCode    string `json:"type_code"`
	TypeName    string `json:"type_name"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
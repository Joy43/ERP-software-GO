package transfer
type CreateTransferRequest struct {
	ItemID         uint    `json:"item_id" binding:"required"`
	TransferQty    float64 `json:"transfer_qty" binding:"required,gt=0"`
	FromOfficeID   uint    `json:"from_office_id" binding:"required"`
	ToOfficeID     uint    `json:"to_office_id" binding:"required"`
	ToLocationID   uint    `json:"to_location_id" binding:"required"`
	GeneralRemarks string  `json:"general_remarks"`
	Note           string  `json:"note"`
}
type UpdateTransferRequest struct {
	TransferQty    float64 `json:"transfer_qty" binding:"required,gt=0"`
	FromOfficeID   uint    `json:"from_office_id" binding:"required"`
	ToOfficeID     uint    `json:"to_office_id" binding:"required"`
	ToLocationID   uint    `json:"to_location_id" binding:"required"`
	GeneralRemarks string  `json:"general_remarks"`
	Note           string  `json:"note"`
}

// ====== pagination response ======
type PaginatedTransferResponse struct {
	Total       int64       `json:"total"`
	Page        int         `json:"page"`
	PageLimit   int         `json:"page_limit"`
	TotalPages  int         `json:"total_pages"`
	Transfers   []Transfer  `json:"transfers"`
}

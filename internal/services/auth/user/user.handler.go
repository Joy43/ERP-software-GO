package user

import (
	"errors"

	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/shared/response"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/shared/utils"
	"gorm.io/gorm"
)

type Handler struct {
	service *Service
	rdb     *redis.Client
}

type registerRequest struct {
	Name          string  `json:"name" binding:"required,min=2,max=150" example:"John Doe"`
	Email         string  `json:"email" binding:"required,email" example:"john@example.com"`
	Password      string  `json:"password" binding:"required,min=6,max=72" example:"123456"`
	Mobile        string  `json:"mobile" binding:"required" example:"01711223344"`
	DOB           string  `json:"dob" example:"1990-01-01"`
	Gender        string  `json:"gender" example:"Male"`
	BloodGroup    string  `json:"blood_group" example:"O+"`
	EmployeeID    string  `json:"employee_id" binding:"required" example:"EMP001"`
	DesignationID *uint   `json:"designation_id" binding:"required" example:"1"`
	DepartmentID  *uint   `json:"department_id" example:"1"`
	OfficeID      *uint   `json:"office_id" binding:"required" example:"1"`
	JoiningDate   string  `json:"joining_date" binding:"required" example:"2023-01-01"`
	BankName      string  `json:"bank_name" example:"Sonali Bank"`
	AccountNumber string  `json:"account_number" example:"1234567890"`
	GrossSalary   float64 `json:"gross_salary" example:"50000"`
	PaymentModeID *uint   `json:"payment_mode_id" example:"1"`
	RoleID        uint    `json:"role_id" example:"1"`
}

type updateUserRequest struct {
	Name          string  `json:"name" example:"John Doe"`
	Email         string  `json:"email" example:"john@example.com"`
	Password      string  `json:"password" example:"123456"`
	Mobile        string  `json:"mobile" example:"01711223344"`
	DOB           string  `json:"dob" example:"1990-01-01"`
	Gender        string  `json:"gender" example:"Male"`
	BloodGroup    string  `json:"blood_group" example:"O+"`
	EmployeeID    string  `json:"employee_id" example:"EMP001"`
	DesignationID *uint   `json:"designation_id" example:"1"`
	DepartmentID  *uint   `json:"department_id" example:"1"`
	OfficeID      *uint   `json:"office_id" example:"1"`
	JoiningDate   string  `json:"joining_date" example:"2023-01-01"`
	BankName      string  `json:"bank_name" example:"Sonali Bank"`
	AccountNumber string  `json:"account_number" example:"1234567890"`
	GrossSalary   float64 `json:"gross_salary" example:"50000"`
	PaymentModeID *uint   `json:"payment_mode_id" example:"1"`
	RoleID        uint    `json:"role_id" example:"1"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"admin@admin.com"`
	Password string `json:"password" binding:"required,min=6,max=72" example:"123456"`
}

type refreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func NewHandler(service *Service, rdb *redis.Client) *Handler {
	return &Handler{service: service, rdb: rdb}
}

// Register godoc
// @Summary Register user
// @Description Registers a new user and returns access and refresh tokens
// @Tags Auth User Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body registerRequest true "register payload"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 409 {object} response.APIResponse
// @Router /auth/users/create-user [post]
func (handler *Handler) CreateUser(ctx *gin.Context) {
	var request registerRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request payload", "VALIDATION_ERROR", err.Error())
		return
	}

	dob, _ := utils.ParseDate(request.DOB)
	joiningDate, _ := utils.ParseDate(request.JoiningDate)

	result, err := handler.service.CreateUser(ctx.Request.Context(), RegisterRequest{
		Name:          request.Name,
		Email:         request.Email,
		Password:      request.Password,
		Mobile:        request.Mobile,
		DOB:           &dob,
		Gender:        request.Gender,
		BloodGroup:    request.BloodGroup,
		EmployeeID:    request.EmployeeID,
		DesignationID: request.DesignationID,
		DepartmentID:  request.DepartmentID,
		OfficeID:      request.OfficeID,
		JoiningDate:   &joiningDate,
		BankName:      request.BankName,
		AccountNumber: request.AccountNumber,
		GrossSalary:   request.GrossSalary,
		PaymentModeID: request.PaymentModeID,
		RoleID:        request.RoleID,
	})
	if err != nil {
		handler.writeServiceError(ctx, err)
		return
	}

	response.Success(ctx, "User created successful", result)
}

// Login godoc
// @Summary Login user
// @Description Authenticates user and returns access and refresh tokens
// @Tags Auth
// @Accept json
// @Produce json
// @Param payload body loginRequest true "login payload"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Router /auth/login [post]
func (handler *Handler) Login(ctx *gin.Context) {
	var request loginRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request payload", "VALIDATION_ERROR", err.Error())
		return
	}

	result, err := handler.service.Login(ctx.Request.Context(), LoginRequest{
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		handler.writeServiceError(ctx, err)
		return
	}

	response.Success(ctx, "login successful", result)
}

// RefreshToken godoc
// @Summary Refresh token
// @Description Issues a new access token and refresh token from a valid refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param payload body refreshTokenRequest true "refresh token payload"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Router /auth/refresh [post]
func (handler *Handler) RefreshToken(ctx *gin.Context) {
	var request refreshTokenRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request payload", "VALIDATION_ERROR", err.Error())
		return
	}

	result, err := handler.service.RefreshToken(ctx.Request.Context(), request.RefreshToken)
	if err != nil {
		handler.writeServiceError(ctx, err)
		return
	}

	response.Success(ctx, "token refreshed", result)
}

// Me godoc
// @Summary Current user profile
// @Description Returns current authenticated user
// @Tags Auth User Management
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Router /auth/users/me [get]
func (handler *Handler) Me(ctx *gin.Context) {
	userIDValue, exists := ctx.Get("user_id")
	if !exists {
		response.Error(ctx, http.StatusUnauthorized, "unauthorized", "UNAUTHORIZED", nil)
		return
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		response.Error(ctx, http.StatusUnauthorized, "unauthorized", "UNAUTHORIZED", nil)
		return
	}

	result, err := handler.service.Profile(ctx.Request.Context(), userID)
	if err != nil {
		handler.writeServiceError(ctx, err)
		return
	}

	response.Success(ctx, "profile fetched", result)
}

// UpdateUser godoc
// @Summary Update user
// @Description Updates user details and employee information. Password change forces logout.
// @Tags Auth User Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param payload body updateUserRequest true "update payload"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Router /auth/users/{id} [put]
func (handler *Handler) UpdateUser(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var request updateUserRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request payload", "VALIDATION_ERROR", err.Error())
		return
	}

	var dobPtr, joiningDatePtr *time.Time
	if request.DOB != "" {
		dob, err := utils.ParseDate(request.DOB)
		if err != nil {
			response.Error(ctx, http.StatusBadRequest, "invalid dob format", "VALIDATION_ERROR", err.Error())
			return
		}
		dobPtr = &dob
	}
	if request.JoiningDate != "" {
		joiningDate, err := utils.ParseDate(request.JoiningDate)
		if err != nil {
			response.Error(ctx, http.StatusBadRequest, "invalid joining_date format", "VALIDATION_ERROR", err.Error())
			return
		}
		joiningDatePtr = &joiningDate
	}

	result, err := handler.service.UpdateUser(ctx.Request.Context(), uint(id), UpdateUserRequest{
		Name:          request.Name,
		Email:         request.Email,
		Password:      request.Password,
		Mobile:        request.Mobile,
		DOB:           dobPtr,
		Gender:        request.Gender,
		BloodGroup:    request.BloodGroup,
		EmployeeID:    request.EmployeeID,
		DesignationID: request.DesignationID,
		DepartmentID:  request.DepartmentID,
		OfficeID:      request.OfficeID,
		JoiningDate:   joiningDatePtr,
		BankName:      request.BankName,
		AccountNumber: request.AccountNumber,
		GrossSalary:   request.GrossSalary,
		PaymentModeID: request.PaymentModeID,
		RoleID:        request.RoleID,
	})
	if err != nil {
		handler.writeServiceError(ctx, err)
		return
	}

	response.Success(ctx, "user updated successfully", result)
}

// DeleteUser godoc
// @Summary Delete user
// @Description Deletes a user and revokes all tokens. Superadmins cannot be deleted.
// @Tags Auth User Management
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Router /auth/users/{id} [delete]
func (handler *Handler) DeleteUser(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if err := handler.service.DeleteUser(ctx.Request.Context(), uint(id)); err != nil {
		handler.writeServiceError(ctx, err)
		return
	}

	response.Success(ctx, "user deleted successfully", nil)
}


//--------- get all user with pagination and search----------------

// List godoc
// @Summary List users
// @Description Returns a list of all users with roles and permissions with pagination.
// @Tags Auth User Management
// @Produce json
// @Security BearerAuth
// @Param page query integer false "Page number (default: 1)" default(1)
// @Param limit query integer false "Items per page (default: 20, max: 100)" default(20)
// @Success 200 {object} response.APIResponse
// @Router /auth/users [get]
func (handler *Handler) List(ctx *gin.Context) {
	var req ListUsersRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid query parameters", "VALIDATION_ERROR", err.Error())
		return
	}

	// Call service with pagination only
	result, err := handler.service.ListUsersPaginated(ctx.Request.Context(), &req)
	if err != nil {
		handler.writeServiceError(ctx, err)
		return
	}

	response.Success(ctx, "users listed successfully", result)
}

// Get godoc
// @Summary Get user
// @Description Returns a single user by ID with roles and permissions
// @Tags Auth User Management
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /auth/users/{id} [get]
func (handler *Handler) Get(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid user id", "INVALID_ID", err.Error())
		return
	}

	result, err := handler.service.GetUser(ctx.Request.Context(), uint(id))
	if err != nil {
		handler.writeServiceError(ctx, err)
		return
	}

	response.Success(ctx, "user fetched successfully", result)
}

func (handler *Handler) writeServiceError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrEmailAlreadyExists):
		response.Error(ctx, http.StatusConflict, "email already exists", "EMAIL_EXISTS", nil)
	case errors.Is(err, ErrInvalidCredentials):
		response.Error(ctx, http.StatusUnauthorized, "invalid credentials", "INVALID_CREDENTIALS", nil)
	case errors.Is(err, ErrInvalidRefreshToken):
		response.Error(ctx, http.StatusUnauthorized, "invalid refresh token", "INVALID_REFRESH_TOKEN", nil)
	case errors.Is(err, ErrUnauthorized):
		response.Error(ctx, http.StatusUnauthorized, "unauthorized", "UNAUTHORIZED", nil)
	case errors.Is(err, gorm.ErrRecordNotFound):
		response.Error(ctx, http.StatusNotFound, "user not found", "NOT_FOUND", nil)
	case errors.Is(err, ErrCannotDeleteSuperadmin):
		response.Error(ctx, http.StatusForbidden, "cannot delete a superadmin user", "FORBIDDEN", nil)
	case errors.Is(err, ErrCannotDemoteSuperadmin):
		response.Error(ctx, http.StatusForbidden, "cannot demote a superadmin user", "FORBIDDEN", nil)
	case errors.Is(err, ErrCannotAssignSuperadmin):
		response.Error(ctx, http.StatusForbidden, "cannot change user role to superadmin", "FORBIDDEN", nil)
	case errors.Is(err, ErrEmployeeIDAlreadyExists):
		response.Error(ctx, http.StatusConflict, "employee id already exists", "EMPLOYEE_ID_EXISTS", nil)
	default:
		response.Error(ctx, http.StatusInternalServerError, "internal server error", "INTERNAL_SERVER_ERROR", nil)
	}
}





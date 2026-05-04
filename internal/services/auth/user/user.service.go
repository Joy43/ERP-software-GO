package user

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/config"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/auth/role"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/shared/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrEmailAlreadyExists  = errors.New("email already exists")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrEmployeeIDAlreadyExists = errors.New("employee id already exists")
	ErrCannotDeleteSuperadmin  = errors.New("cannot delete a superadmin user")
	ErrCannotDemoteSuperadmin  = errors.New("cannot demote a superadmin user")
	ErrCannotAssignSuperadmin  = errors.New("cannot change user role to superadmin")
)

type Service struct {
	repository Repository
	roleRepo   role.Repository
	config     config.Config
}

type RegisterRequest struct {
	Name          string
	Email         string
	Password      string
	Mobile        string
	DOB           *time.Time
	Gender        string
	BloodGroup    string
	EmployeeID    string
	DesignationID *uint
	DepartmentID  *uint
	OfficeID      *uint
	JoiningDate   *time.Time
	BankName      string
	AccountNumber string
	GrossSalary   float64
	PaymentModeID *uint
	RoleID        uint
}

type UpdateUserRequest struct {
	Name          string     `json:"name"`
	Email         string     `json:"email"`
	Password      string     `json:"password"`
	Mobile        string     `json:"mobile"`
	DOB           *time.Time `json:"dob"`
	Gender        string     `json:"gender"`
	BloodGroup    string     `json:"blood_group"`
	EmployeeID    string     `json:"employee_id"`
	DesignationID *uint      `json:"designation_id"`
	DepartmentID  *uint      `json:"department_id"`
	OfficeID      *uint      `json:"office_id"`
	JoiningDate   *time.Time `json:"joining_date"`
	BankName      string     `json:"bank_name"`
	AccountNumber string     `json:"account_number"`
	GrossSalary   float64    `json:"gross_salary"`
	PaymentModeID *uint      `json:"payment_mode_id"`
	RoleID        uint       `json:"role_id"`
}

type LoginRequest struct {
	Email    string
	Password string
}

type AuthResponse struct {
	AccessToken           string      `json:"access_token"`
	RefreshToken          string      `json:"refresh_token"`
	TokenType             string      `json:"token_type"`
	AccessTokenExpiresAt  time.Time   `json:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time   `json:"refresh_token_expires_at"`
	User                  UserProfile `json:"user"`
}

func NewService(repository Repository, roleRepo role.Repository, cfg config.Config) *Service {
	return &Service{repository: repository, roleRepo: roleRepo, config: cfg}
}

func (service *Service) CreateUser(ctx context.Context, request RegisterRequest) (*UserProfile, error) {
	email := strings.ToLower(strings.TrimSpace(request.Email))

	existingUser, err := service.repository.FindByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, ErrEmailAlreadyExists
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Check if employee ID already exists
	if request.EmployeeID != "" {
		existingUserByEMP, err := service.repository.FindByEmployeeID(ctx, request.EmployeeID)
		if err == nil && existingUserByEMP != nil {
			return nil, ErrEmployeeIDAlreadyExists
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(request.Password), service.config.BCryptCost)
	if err != nil {
		return nil, fmt.Errorf("generate password hash: %w", err)
	}

	value := &User{
		Name:          request.Name,
		Email:         email,
		PasswordHash:  string(hash),
		Mobile:        request.Mobile,
		DOB:           request.DOB,
		Gender:        request.Gender,
		BloodGroup:    request.BloodGroup,
		EmployeeID:    request.EmployeeID,
		DesignationID: request.DesignationID,
		DepartmentID:  request.DepartmentID,
		OfficeID:      request.OfficeID,
		JoiningDate:   request.JoiningDate,
		BankName:      request.BankName,
		AccountNumber: request.AccountNumber,
		GrossSalary:   request.GrossSalary,
		PaymentModeID: request.PaymentModeID,
		IsActive:      true,
	}

	if request.RoleID != 0 {
		value.Roles = []role.Role{{ID: request.RoleID}}
	}

	if err := service.repository.CreateUser(ctx, value); err != nil {
		return nil, err
	}

	profile := toUserProfile(value)
	return &profile, nil
}

func (service *Service) Login(ctx context.Context, request LoginRequest) (*AuthResponse, error) {
	email := strings.ToLower(strings.TrimSpace(request.Email))

	value, err := service.repository.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if !value.IsActive {
		return nil, ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(value.PasswordHash), []byte(request.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return service.issueTokens(ctx, value)
}

func (service *Service) RefreshToken(ctx context.Context, token string) (*AuthResponse, error) {
	if _, err := utils.ParseRefreshToken(token, service.config.JWTRefreshSecret); err != nil {
		return nil, ErrInvalidRefreshToken
	}

	tokenHash := hashToken(token)
	storedToken, err := service.repository.FindRefreshTokenByHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidRefreshToken
		}
		return nil, err
	}

	if err := service.repository.RevokeRefreshToken(ctx, storedToken.ID); err != nil {
		return nil, err
	}

	value, err := service.repository.FindByID(ctx, storedToken.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUnauthorized
		}
		return nil, err
	}

	if !value.IsActive {
		return nil, ErrUnauthorized
	}

	return service.issueTokens(ctx, value)
}

func (service *Service) Profile(ctx context.Context, userID uint) (*UserProfile, error) {
	value, err := service.repository.FindByIDWithRoles(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUnauthorized
		}
		return nil, err
	}

	if !value.IsActive {
		return nil, ErrUnauthorized
	}

	profile := toUserProfile(value)
	return &profile, nil
}

func (service *Service) ListUsers(ctx context.Context) ([]UserProfile, error) {
	users, err := service.repository.FindAllWithRoles(ctx)
	if err != nil {
		return nil, err
	}

	profiles := make([]UserProfile, len(users))
	for i, u := range users {
		profiles[i] = toUserProfile(&u)
	}

	return profiles, nil
}

// ListUsersPaginated retrieves users with pagination only
func (service *Service) ListUsersPaginated(ctx context.Context, req *ListUsersRequest) (*ListUsersResponse, error) {
	// Validate pagination parameters
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	// Call repository with pagination only
	users, total, err := service.repository.FindAllWithRolesPaginated(
		ctx,
		req.Page,
		req.Limit,
	)
	if err != nil {
		return nil, err
	}

	// Convert users to profiles
	profiles := make([]UserProfileResponse, len(users))
	for i, u := range users {
		profiles[i] = toUserProfileResponse(&u)
	}

	// Calculate total pages
	totalPages := (total + int64(req.Limit) - 1) / int64(req.Limit)

	// Return nil for data if no results
	var responseData []UserProfileResponse
	if len(profiles) > 0 {
		responseData = profiles
	}

	return &ListUsersResponse{
		Data:       responseData,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
	}, nil
}

func (service *Service) GetUser(ctx context.Context, id uint) (*UserProfile, error) {
	user, err := service.repository.FindByIDWithRoles(ctx, id)
	if err != nil {
		return nil, err
	}

	profile := toUserProfile(user)
	return &profile, nil
}

func (service *Service) issueTokens(ctx context.Context, value *User) (*AuthResponse, error) {
	accessToken, accessTokenExpiresAt, err := utils.GenerateAccessToken(
		value.ID,
		value.Email,
		service.config.JWTAccessSecret,
		service.config.JWTAccessTTL,
	)
	if err != nil {
		return nil, err
	}

	refreshToken, refreshTokenExpiresAt, err := utils.GenerateRefreshToken(
		value.ID,
		service.config.JWTRefreshSecret,
		service.config.JWTRefreshTTL,
	)
	if err != nil {
		return nil, err
	}

	storedToken := &RefreshToken{
		UserID:    value.ID,
		TokenHash: hashToken(refreshToken),
		ExpiresAt: refreshTokenExpiresAt,
	}
	if err := service.repository.CreateRefreshToken(ctx, storedToken); err != nil {
		return nil, err
	}

	result := &AuthResponse{
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		TokenType:             "Bearer",
		AccessTokenExpiresAt:  accessTokenExpiresAt,
		RefreshTokenExpiresAt: refreshTokenExpiresAt,
		User:                  toUserProfile(value),
	}

	return result, nil
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func toUserProfile(value *User) UserProfile {
	return UserProfile{
		ID:            value.ID,
		Name:          value.Name,
		Email:         value.Email,
		Mobile:        value.Mobile,
		DOB:           value.DOB,
		Gender:        value.Gender,
		BloodGroup:    value.BloodGroup,
		EmployeeID:    value.EmployeeID,
		DesignationID: value.DesignationID,
		DepartmentID:  value.DepartmentID,
		OfficeID:      value.OfficeID,
		JoiningDate:   value.JoiningDate,
		BankName:      value.BankName,
		AccountNumber: value.AccountNumber,
		GrossSalary:   value.GrossSalary,
		PaymentModeID: value.PaymentModeID,
		Roles:         value.Roles,
	}
}

// toUserProfileResponse converts a User to UserProfileResponse with CreatedAt
func toUserProfileResponse(value *User) UserProfileResponse {
	return UserProfileResponse{
		ID:            value.ID,
		Name:          value.Name,
		Email:         value.Email,
		Mobile:        value.Mobile,
		DOB:           value.DOB,
		Gender:        value.Gender,
		BloodGroup:    value.BloodGroup,
		EmployeeID:    value.EmployeeID,
		DesignationID: value.DesignationID,
		DepartmentID:  value.DepartmentID,
		OfficeID:      value.OfficeID,
		JoiningDate:   value.JoiningDate,
		BankName:      value.BankName,
		AccountNumber: value.AccountNumber,
		GrossSalary:   value.GrossSalary,
		PaymentModeID: value.PaymentModeID,
		Roles:         value.Roles,
		CreatedAt:     value.CreatedAt,
	}
}

func (service *Service) UpdateUser(ctx context.Context, id uint, request UpdateUserRequest) (*UserProfile, error) {
	// Fetch existing user with roles
	user, err := service.repository.FindByIDWithRoles(ctx, id)
	if err != nil {
		return nil, err
	}

	// Security Check: Prevent role changes to/from Superadmin
	isTargetSuperadmin := service.isSuperadmin(user)

	// If changing role, check if new role is Superadmin or if we are demoting a Superadmin
	if request.RoleID != 0 {
		newRole, err := service.roleRepo.FindByID(ctx, request.RoleID)
		if err != nil {
			return nil, fmt.Errorf("invalid role: %w", err)
		}

		if newRole.Slug == "superadmin" {
			return nil, ErrCannotAssignSuperadmin
		}

		if isTargetSuperadmin && newRole.Slug != "superadmin" {
			return nil, ErrCannotDemoteSuperadmin
		}

		user.Roles = []role.Role{*newRole}
	}

	// Update basic fields
	if request.Name != "" {
		user.Name = request.Name
	}

	if request.Email != "" {
		newEmail := strings.ToLower(strings.TrimSpace(request.Email))
		if newEmail != user.Email {
			existingUser, err := service.repository.FindByEmail(ctx, newEmail)
			if err == nil && existingUser != nil {
				return nil, ErrEmailAlreadyExists
			}
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, err
			}
			user.Email = newEmail
		}
	}

	passwordChanged := false
	if request.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(request.Password), service.config.BCryptCost)
		if err != nil {
			return nil, err
		}
		user.PasswordHash = string(hash)
		passwordChanged = true
	}

	// Update employee fields
	if request.Mobile != "" {
		user.Mobile = request.Mobile
	}
	if request.DOB != nil {
		user.DOB = request.DOB
	}
	if request.Gender != "" {
		user.Gender = request.Gender
	}
	if request.BloodGroup != "" {
		user.BloodGroup = request.BloodGroup
	}
	if request.EmployeeID != "" && request.EmployeeID != user.EmployeeID {
		existing, err := service.repository.FindByEmployeeID(ctx, request.EmployeeID)
		if err == nil && existing != nil {
			return nil, ErrEmployeeIDAlreadyExists
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		user.EmployeeID = request.EmployeeID
	}
	if request.DesignationID != nil {
		user.DesignationID = request.DesignationID
	}
	if request.DepartmentID != nil {
		user.DepartmentID = request.DepartmentID
	}
	if request.OfficeID != nil {
		user.OfficeID = request.OfficeID
	}
	if request.JoiningDate != nil {
		user.JoiningDate = request.JoiningDate
	}
	if request.BankName != "" {
		user.BankName = request.BankName
	}
	if request.AccountNumber != "" {
		user.AccountNumber = request.AccountNumber
	}
	if request.GrossSalary != 0 {
		user.GrossSalary = request.GrossSalary
	}
	if request.PaymentModeID != nil {
		user.PaymentModeID = request.PaymentModeID
	}

	if err := service.repository.UpdateUser(ctx, user); err != nil {
		return nil, err
	}

	// If password changed, revoke all tokens (logout)
	if passwordChanged {
		if err := service.repository.RevokeAllUserRefreshTokens(ctx, user.ID); err != nil {
			return nil, fmt.Errorf("logout after password change failed: %w", err)
		}
	}

	profile := toUserProfile(user)
	return &profile, nil
}

func (service *Service) DeleteUser(ctx context.Context, id uint) error {
	user, err := service.repository.FindByIDWithRoles(ctx, id)
	if err != nil {
		return err
	}

	if service.isSuperadmin(user) {
		return ErrCannotDeleteSuperadmin
	}

	return service.repository.DeleteUser(ctx, id)
}

func (service *Service) isSuperadmin(u *User) bool {
	for _, r := range u.Roles {
		if r.Slug == "superadmin" {
			return true
		}
	}
	return false
}

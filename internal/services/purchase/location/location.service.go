package location

import "errors"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

//
// =============================
// VALIDATION
// =============================
func (s *Service) ValidateCreateLocationRequest(req *CreateLocationRequest) error {
	if req.Name == "" {
		return errors.New("name is required")
	}

	if req.Code == "" {
		return errors.New("code is required")
	}

	if req.Type == "" {
		return errors.New("type is required")
	}

	if req.OfficeID == 0 {
		return errors.New("office_id is required")
	}

	return nil
}

//
// =============================
// CREATE
// =============================
func (s *Service) CreateLocation(req CreateLocationRequest) (*Location, error) {
	existing, _ := s.repo.GetByCode(req.Code)
	if existing != nil {
		return nil, errors.New("location code already exists")
	}

	loc := &Location{
		Name:     req.Name,
		Code:     req.Code,
		Type:     req.Type,
		OfficeID: req.OfficeID,

		ParentID:  req.ParentID,
		ManagerID: req.ManagerID,

		Location: req.Location,
	}

	// ----------- default true if nil -----------
	if req.IsActive != nil {
		loc.IsActive = *req.IsActive
	} else {
		loc.IsActive = true
	}

	if err := s.repo.Create(loc); err != nil {
		return nil, err
	}

	return loc, nil
}

//
// =============================
//---------  GET SINGLE -------
// =============================
func (s *Service) GetLocation(id uint) (*Location, error) {
	if id == 0 {
		return nil, errors.New("invalid location id")
	}

	return s.repo.GetByID(id)
}

//
// =============================
//---------  GET ALL -------
// =============================
func (s *Service) GetAllLocations() ([]Location, error) {
	return s.repo.ListAll()
}



//
// =============================
// --------DELETE (SOFT DELETE) ------
// =============================
func (s *Service) DeleteLocation(id uint) error {
	if id == 0 {
		return errors.New("invalid location id")
	}

	return s.repo.Delete(id)
}

//
// =============================
// ------------ UPDATE --------------
// =============================
func (s *Service) UpdateLocation(id uint, req UpdateLocationRequest) (*Location, error) {

	if id == 0 {
		return nil, errors.New("invalid location id")
	}

	loc, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("location not found")
	}

	if req.Name != nil {
		loc.Name = *req.Name
	}

	if req.Code != nil {
		existing, _ := s.repo.GetByCode(*req.Code)
		if existing != nil && existing.ID != loc.ID {
			return nil, errors.New("location code already exists")
		}
		loc.Code = *req.Code
	}

	if req.Type != nil {
		loc.Type = *req.Type
	}

	if req.OfficeID != nil {
		loc.OfficeID = *req.OfficeID
	}

	if req.ParentID != nil {
		loc.ParentID = req.ParentID
	}

	if req.ManagerID != nil {
		loc.ManagerID = req.ManagerID
	}

	if req.Location != nil {
		loc.Location = *req.Location
	}

	if req.IsActive != nil {
		loc.IsActive = *req.IsActive
	}

	if err := s.repo.Update(loc); err != nil {
		return nil, err
	}

	return loc, nil
}
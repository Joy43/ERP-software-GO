package customer

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, customer *Customer) error
	Update(ctx context.Context, customer *Customer) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*Customer, error)
	FindAll(ctx context.Context) ([]Customer, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, customer *Customer) error {
	return r.db.WithContext(ctx).Create(customer).Error
}

func (r *repository) Update(ctx context.Context, customer *Customer) error {
	return r.db.WithContext(ctx).Save(customer).Error
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&Customer{}, id).Error
}

func (r *repository) FindByID(ctx context.Context, id uint) (*Customer, error) {
	var c Customer
	err := r.db.WithContext(ctx).
		Preload("Office").
		Preload("PartnerGroup").
		Preload("PartnerSubGroup").
		Preload("DefaultSalesRep").
		Preload("PriceType").
		Preload("User").
		Preload("District").
		Preload("Thana").
		Preload("TaxBracket").
		Preload("ShippingAddresses").
		Preload("ShippingAddresses.District").
		Preload("ShippingAddresses.Thana").
		Preload("BankInfos").
		First(&c, id).Error
	return &c, err
}

func (r *repository) FindAll(ctx context.Context) ([]Customer, error) {
	var customers []Customer
	err := r.db.WithContext(ctx).
		Preload("Office").
		Preload("PartnerGroup").
		Preload("PartnerSubGroup").
		Preload("DefaultSalesRep").
		Preload("PriceType").
		Preload("District").
		Preload("Thana").
		Order("id desc").
		Find(&customers).Error
	return customers, err
}

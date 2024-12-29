package service

import (
	"context"
	"errors"
	"github.com/rusystem/crm-api/internal/config"
	"github.com/rusystem/crm-api/internal/repository"
	"github.com/rusystem/crm-api/pkg/domain"
	"github.com/rusystem/crm-api/tools"
)

type Supplier interface {
	GetById(ctx context.Context, id int64, info domain.JWTInfo) (domain.Supplier, error)
	Create(ctx context.Context, spl domain.Supplier) (int64, error)
	Update(ctx context.Context, inp domain.UpdateSupplier, info domain.JWTInfo) error
	Delete(ctx context.Context, id int64, info domain.JWTInfo) error
	GetListByCompanyId(ctx context.Context, companyId int64, param domain.Param) ([]domain.Supplier, int64, error)
}

type SupplierService struct {
	cfg  *config.Config
	repo *repository.Repository
}

func NewSupplierService(cfg *config.Config, repo *repository.Repository) *SupplierService {
	return &SupplierService{
		cfg:  cfg,
		repo: repo,
	}
}

func (s *SupplierService) GetById(ctx context.Context, id int64, info domain.JWTInfo) (domain.Supplier, error) {
	supplier, err := s.repo.Suppliers.GetById(ctx, id)
	if err != nil {
		return domain.Supplier{}, err
	}

	if supplier.CompanyId != info.CompanyId && !tools.IsFullAccessSection(info.Sections) {
		return domain.Supplier{}, domain.ErrNotAllowed
	}

	return supplier, nil
}

func (s *SupplierService) Create(ctx context.Context, spl domain.Supplier) (int64, error) {
	return s.repo.Suppliers.Create(ctx, spl)
}

func (s *SupplierService) Update(ctx context.Context, inp domain.UpdateSupplier, info domain.JWTInfo) error {
	supplier, err := s.repo.Suppliers.GetById(ctx, inp.Id)
	if err != nil {
		if errors.Is(err, domain.ErrSupplierNotFound) {
			return domain.ErrSupplierNotFound
		}

		return err
	}

	if supplier.CompanyId != info.CompanyId && !tools.IsFullAccessSection(info.Sections) {
		return domain.ErrNotAllowed
	}

	if inp.Name != nil {
		supplier.Name = *inp.Name
	}

	if inp.LegalAddress != nil {
		supplier.LegalAddress = *inp.LegalAddress
	}

	if inp.ActualAddress != nil {
		supplier.ActualAddress = *inp.ActualAddress
	}

	if inp.WarehouseAddress != nil {
		supplier.WarehouseAddress = *inp.WarehouseAddress
	}

	if inp.ContactPerson != nil {
		supplier.ContactPerson = *inp.ContactPerson
	}

	if inp.Phone != nil {
		supplier.Phone = *inp.Phone
	}

	if inp.Email != nil {
		supplier.Email = *inp.Email
	}

	if inp.Website != nil {
		supplier.Website = *inp.Website
	}

	if inp.ContractNumber != nil {
		supplier.ContractNumber = *inp.ContractNumber
	}

	if inp.ProductCategories != nil {
		supplier.ProductCategories = *inp.ProductCategories
	}

	if inp.Balance != nil {
		supplier.Balance = *inp.Balance
	}

	if inp.ProductTypes != nil {
		supplier.ProductTypes = *inp.ProductTypes
	}

	if inp.Comments != nil {
		supplier.Comments = *inp.Comments
	}

	if inp.Files != nil {
		supplier.Files = *inp.Files
	}

	if inp.Country != nil {
		supplier.Country = *inp.Country
	}

	if inp.Region != nil {
		supplier.Region = *inp.Region
	}

	if inp.TaxID != nil {
		supplier.TaxID = *inp.TaxID
	}

	if inp.BankDetails != nil {
		supplier.BankDetails = *inp.BankDetails
	}

	if inp.PaymentTerms != nil {
		supplier.PaymentTerms = *inp.PaymentTerms
	}

	if inp.IsActive != nil {
		supplier.IsActive = *inp.IsActive
	}

	if inp.OtherFields != nil {
		supplier.OtherFields = *inp.OtherFields
	}

	return s.repo.Suppliers.Update(ctx, supplier)
}

func (s *SupplierService) Delete(ctx context.Context, id int64, info domain.JWTInfo) error {
	spl, err := s.repo.Suppliers.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrSupplierNotFound) {
			return domain.ErrSupplierNotFound
		}

		return err
	}

	if spl.CompanyId != info.CompanyId && !tools.IsFullAccessSection(info.Sections) {
		return domain.ErrNotAllowed
	}

	return s.repo.Suppliers.Delete(ctx, id)
}

func (s *SupplierService) GetListByCompanyId(ctx context.Context, companyId int64, param domain.Param) ([]domain.Supplier, int64, error) {
	return s.repo.Suppliers.GetListByCompanyId(ctx, companyId, param)
}

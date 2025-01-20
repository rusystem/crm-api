package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"github.com/rusystem/crm-api/internal/config"
	"github.com/rusystem/crm-api/internal/repository"
	"github.com/rusystem/crm-api/pkg/domain"
	"github.com/rusystem/crm-api/tools"
	"github.com/xuri/excelize/v2"
	"strconv"
)

type Warehouse interface {
	GetById(ctx context.Context, id int64, info domain.JWTInfo) (domain.Warehouse, error)
	Create(ctx context.Context, wh domain.Warehouse) (int64, error)
	Update(ctx context.Context, wh domain.WarehouseUpdate, info domain.JWTInfo) error
	Delete(ctx context.Context, id int64, info domain.JWTInfo) error
	GetListByCompanyId(ctx context.Context, companyId int64, param domain.Param) ([]domain.Warehouse, int64, error)
	GetResponsibleUsers(ctx context.Context, companyId int64, param domain.Param) ([]domain.UserResponse, int64, error)
	GetIncomeHistoryByWarehouseId(ctx context.Context, id int64, param domain.Param) ([]domain.Material, int64, error)
	GenerateWarehouseInfoReportXls(ctx context.Context, id int64, info domain.JWTInfo) (*excelize.File, error)
	GenerateWarehouseInfoReportPdf(ctx context.Context, id int64, info domain.JWTInfo) (*gofpdf.Fpdf, error)
}

type WarehouseServices struct {
	cfg  *config.Config
	repo *repository.Repository
}

func NewWarehouseServices(cfg *config.Config, repo *repository.Repository) *WarehouseServices {
	return &WarehouseServices{
		cfg:  cfg,
		repo: repo,
	}
}

func (s *WarehouseServices) GetById(ctx context.Context, id int64, info domain.JWTInfo) (domain.Warehouse, error) {
	wh, err := s.repo.Warehouse.GetById(ctx, id)
	if err != nil {
		return domain.Warehouse{}, err
	}

	if wh.CompanyId != info.CompanyId {
		return domain.Warehouse{}, domain.ErrNotAllowed
	}

	return wh, nil
}

func (s *WarehouseServices) Create(ctx context.Context, wh domain.Warehouse) (int64, error) {
	return s.repo.Warehouse.Create(ctx, wh)
}

func (s *WarehouseServices) Update(ctx context.Context, inp domain.WarehouseUpdate, info domain.JWTInfo) error {
	wh, err := s.repo.Warehouse.GetById(ctx, inp.ID)
	if err != nil {
		return err
	}

	if wh.CompanyId != info.CompanyId && !tools.IsFullAccessSection(info.Sections) {
		return domain.ErrNotAllowed
	}

	if inp.Name != nil {
		wh.Name = *inp.Name
	}

	if inp.Address != nil {
		wh.Address = *inp.Address
	}

	if inp.ResponsiblePerson != nil {
		wh.ResponsiblePerson = *inp.ResponsiblePerson
	}

	if inp.Phone != nil {
		wh.Phone = *inp.Phone
	}

	if inp.Email != nil {
		wh.Email = *inp.Email
	}

	if inp.MaxCapacity != nil {
		wh.MaxCapacity = *inp.MaxCapacity
	}

	if inp.CurrentOccupancy != nil {
		wh.CurrentOccupancy = *inp.CurrentOccupancy
	}

	if inp.OtherFields != nil {
		wh.OtherFields = *inp.OtherFields
	}

	if inp.Country != nil {
		wh.Country = *inp.Country
	}

	if inp.Region != nil {
		wh.Region = *inp.Region
	}

	if inp.Locality != nil {
		wh.Locality = *inp.Locality
	}

	if inp.Comments != nil {
		wh.Comments = *inp.Comments
	}

	return s.repo.Warehouse.Update(ctx, wh)
}

func (s *WarehouseServices) Delete(ctx context.Context, id int64, info domain.JWTInfo) error {
	wh, err := s.repo.Warehouse.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrWarehouseNotFound) {
			return domain.ErrWarehouseNotFound
		}

		return err
	}

	if wh.CompanyId != info.CompanyId && !tools.IsFullAccessSection(info.Sections) {
		return domain.ErrNotAllowed
	}

	return s.repo.Warehouse.Delete(ctx, id)
}

func (s *WarehouseServices) GetListByCompanyId(ctx context.Context, companyId int64, param domain.Param) ([]domain.Warehouse, int64, error) {
	return s.repo.Warehouse.GetListByCompanyId(ctx, companyId, param)
}

func (s *WarehouseServices) GetResponsibleUsers(ctx context.Context, companyId int64, param domain.Param) ([]domain.UserResponse, int64, error) {
	var resp []domain.UserResponse

	users, count, err := s.repo.Warehouse.GetResponsibleUsers(ctx, companyId, param)
	if err != nil {
		return resp, 0, err
	}

	for _, v := range users {
		resp = append(resp, domain.UserResponse{
			ID:                       v.ID,
			CompanyID:                v.CompanyID,
			Username:                 v.Username,
			Name:                     v.Name,
			Email:                    v.Email,
			Phone:                    v.Phone,
			CreatedAt:                v.CreatedAt,
			UpdatedAt:                v.UpdatedAt,
			LastLogin:                v.LastLogin.Time,
			IsActive:                 v.IsActive,
			Role:                     v.Role,
			Language:                 v.Language,
			Country:                  v.Country,
			IsApproved:               v.IsApproved,
			IsSendSystemNotification: v.IsSendSystemNotification,
			Sections:                 v.Sections,
			Position:                 v.Position,
		})
	}

	return resp, count, nil
}

func (s *WarehouseServices) GetIncomeHistoryByWarehouseId(ctx context.Context, id int64, param domain.Param) ([]domain.Material, int64, error) {
	return s.repo.Materials.GetIncomeHistoryByWarehouseId(ctx, id, param)
}

func (s *WarehouseServices) GenerateWarehouseInfoReportXls(ctx context.Context, id int64, info domain.JWTInfo) (*excelize.File, error) {
	wh, err := s.repo.Warehouse.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	if wh.CompanyId != info.CompanyId {
		return nil, domain.ErrNotAllowed
	}

	f := excelize.NewFile()
	sheetName := "Warehouse Report"
	f.SetSheetName("Sheet1", sheetName)

	responsibleUser, err := s.repo.User.GetById(ctx, wh.ResponsiblePerson)
	if err != nil {
		return nil, err
	}

	fields := []struct {
		Label string
		Value string
	}{
		{"ID", strconv.Itoa(int(wh.ID))},
		{"Наименование", wh.Name},
		{"Адрес", wh.Address},
		{"Ответственный", responsibleUser.Name},
		{"Телефон", wh.Phone},
		{"Электронная почта", wh.Email},
		{"Макс. вместимость", fmt.Sprintf("%d м²", wh.MaxCapacity)},
		{"Текущая заполняемость", fmt.Sprintf("%d м²", wh.CurrentOccupancy)},
		{"Страна", wh.Country},
		{"Регион", wh.Region},
		{"Локалитет", wh.Locality},
		{"Комментарии", wh.Comments},
	}

	maxWidth := 0
	for i, field := range fields {
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", i+1), field.Label)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", i+1), field.Value)
		if len(field.Label) > maxWidth {
			maxWidth = len(field.Label)
		}
	}

	f.SetColWidth(sheetName, "A", "A", float64(maxWidth)+2)

	return f, nil
}

func (s *WarehouseServices) GenerateWarehouseInfoReportPdf(ctx context.Context, id int64, info domain.JWTInfo) (*gofpdf.Fpdf, error) {
	wh, err := s.repo.Warehouse.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	if wh.CompanyId != info.CompanyId {
		return nil, domain.ErrNotAllowed
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddUTF8Font("Arial", "", "assets/fonts/arial/arialmt.ttf")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 16)

	pdf.Cell(190, 10, fmt.Sprintf("Отчет по складу: %s", wh.Name))
	pdf.Ln(12)

	responsibleUser, err := s.repo.User.GetById(ctx, wh.ResponsiblePerson)
	if err != nil {
		return nil, err
	}

	pdf.SetFont("Arial", "", 12)
	fields := []struct {
		Label string
		Value string
	}{
		{"ID", strconv.Itoa(int(wh.ID))},
		{"Наименование", wh.Name},
		{"Адрес", wh.Address},
		{"Ответственный", responsibleUser.Name},
		{"Телефон", wh.Phone},
		{"Электронная почта", wh.Email},
		{"Макс. вместимость", fmt.Sprintf("%d м.кв.", wh.MaxCapacity)},
		{"Текущая заполняемость", fmt.Sprintf("%d м.кв.", wh.CurrentOccupancy)},
		{"Страна", wh.Country},
		{"Регион", wh.Region},
		{"Локалитет", wh.Locality},
		{"Комментарии", wh.Comments},
	}

	maxWidth := 0.0
	for _, field := range fields {
		width := pdf.GetStringWidth(field.Label + ":")
		if width > maxWidth {
			maxWidth = width
		}
	}

	for _, field := range fields {
		pdf.CellFormat(maxWidth+2, 10, field.Label+":", "", 0, "L", false, 0, "")
		pdf.Cell(100, 10, field.Value)
		pdf.Ln(8)
	}

	return pdf, nil
}

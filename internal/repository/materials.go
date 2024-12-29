package repository

import (
	"context"
	"database/sql"
	"github.com/rusystem/crm-api/internal/config"
	"github.com/rusystem/crm-api/internal/repository/database"
	"github.com/rusystem/crm-api/pkg/domain"
)

type Materials interface {
	CreatePlanning(ctx context.Context, material domain.Material) (int64, error)
	UpdatePlanning(ctx context.Context, material domain.Material) error
	DeletePlanning(ctx context.Context, id int64) error
	GetPlanningById(ctx context.Context, id int64) (domain.Material, error)
	GetPlanningList(ctx context.Context, params domain.MaterialParams) ([]domain.Material, int64, error)
	MovePlanningToPurchased(ctx context.Context, id int64) (int64, int64, error)

	CreatePurchased(ctx context.Context, material domain.Material) (int64, int64, error)
	UpdatePurchased(ctx context.Context, material domain.Material) error
	DeletePurchased(ctx context.Context, id int64) error
	GetPurchasedById(ctx context.Context, id int64) (domain.Material, error)
	GetPurchasedList(ctx context.Context, params domain.MaterialParams) ([]domain.Material, int64, error)
	MovePurchasedToArchive(ctx context.Context, id int64) error

	GetPlanningArchiveById(ctx context.Context, id int64) (domain.Material, error)
	GetPurchasedArchiveById(ctx context.Context, id int64) (domain.Material, error)
	GetPlanningArchiveList(ctx context.Context, params domain.MaterialParams) ([]domain.Material, int64, error)
	GetPurchasedArchiveList(ctx context.Context, params domain.MaterialParams) ([]domain.Material, int64, error)
	DeletePlanningArchive(ctx context.Context, id int64) error
	DeletePurchasedArchive(ctx context.Context, id int64) error

	Search(ctx context.Context, param domain.MaterialParams) ([]domain.Material, int64, error)
	GetIncomeHistoryByWarehouseId(ctx context.Context, id int64, param domain.Param) ([]domain.Material, int64, error)
}

type MaterialsRepository struct {
	cfg  *config.Config
	psql database.Materials
}

func NewMaterialsRepository(cfg *config.Config, db *sql.DB) *MaterialsRepository {
	return &MaterialsRepository{
		cfg:  cfg,
		psql: database.NewMaterialsPostgresRepository(db),
	}
}

func (mr *MaterialsRepository) CreatePlanning(ctx context.Context, material domain.Material) (int64, error) {
	return mr.psql.CreatePlanning(ctx, material)
}

func (mr *MaterialsRepository) UpdatePlanning(ctx context.Context, material domain.Material) error {
	return mr.psql.UpdatePlanning(ctx, material)
}

func (mr *MaterialsRepository) DeletePlanning(ctx context.Context, id int64) error {
	return mr.psql.DeletePlanning(ctx, id)
}

func (mr *MaterialsRepository) GetPlanningById(ctx context.Context, id int64) (domain.Material, error) {
	return mr.psql.GetPlanningById(ctx, id)
}

func (mr *MaterialsRepository) GetPlanningList(ctx context.Context, params domain.MaterialParams) ([]domain.Material, int64, error) {
	return mr.psql.GetPlanningList(ctx, params)
}

func (mr *MaterialsRepository) MovePlanningToPurchased(ctx context.Context, id int64) (int64, int64, error) {
	return mr.psql.MovePlanningToPurchased(ctx, id)
}

func (mr *MaterialsRepository) CreatePurchased(ctx context.Context, material domain.Material) (int64, int64, error) {
	return mr.psql.CreatePurchased(ctx, material)
}

func (mr *MaterialsRepository) UpdatePurchased(ctx context.Context, material domain.Material) error {
	return mr.psql.UpdatePurchased(ctx, material)
}

func (mr *MaterialsRepository) DeletePurchased(ctx context.Context, id int64) error {
	return mr.psql.DeletePurchased(ctx, id)
}

func (mr *MaterialsRepository) GetPurchasedById(ctx context.Context, id int64) (domain.Material, error) {
	return mr.psql.GetPurchasedById(ctx, id)
}

func (mr *MaterialsRepository) GetPurchasedList(ctx context.Context, params domain.MaterialParams) ([]domain.Material, int64, error) {
	return mr.psql.GetPurchasedList(ctx, params)
}

func (mr *MaterialsRepository) MovePurchasedToArchive(ctx context.Context, id int64) error {
	return mr.psql.MovePurchasedToArchive(ctx, id)
}

func (mr *MaterialsRepository) GetPlanningArchiveById(ctx context.Context, id int64) (domain.Material, error) {
	return mr.psql.GetPlanningArchiveById(ctx, id)
}

func (mr *MaterialsRepository) GetPurchasedArchiveById(ctx context.Context, id int64) (domain.Material, error) {
	return mr.psql.GetPurchasedArchiveById(ctx, id)
}

func (mr *MaterialsRepository) GetPlanningArchiveList(ctx context.Context, params domain.MaterialParams) ([]domain.Material, int64, error) {
	return mr.psql.GetPlanningArchiveList(ctx, params)
}

func (mr *MaterialsRepository) GetPurchasedArchiveList(ctx context.Context, params domain.MaterialParams) ([]domain.Material, int64, error) {
	return mr.psql.GetPurchasedArchiveList(ctx, params)
}

func (mr *MaterialsRepository) DeletePlanningArchive(ctx context.Context, id int64) error {
	return mr.psql.DeletePlanningArchive(ctx, id)
}

func (mr *MaterialsRepository) DeletePurchasedArchive(ctx context.Context, id int64) error {
	return mr.psql.DeletePurchasedArchive(ctx, id)
}

func (mr *MaterialsRepository) Search(ctx context.Context, param domain.MaterialParams) ([]domain.Material, int64, error) {
	return mr.psql.Search(ctx, param)
}

func (mr *MaterialsRepository) GetIncomeHistoryByWarehouseId(ctx context.Context, id int64, param domain.Param) ([]domain.Material, int64, error) {
	return mr.psql.GetIncomeHistoryByWarehouseId(ctx, id, param)
}

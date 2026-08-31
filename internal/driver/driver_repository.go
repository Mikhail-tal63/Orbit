package driver

import (
	"context"

	"github.com/Mikhail-Tal63/Orbit/internal/db"
	"github.com/google/uuid"
)

type DriverRepository interface {
	CreateDriver(ctx context.Context, params db.CreateDriverParams) (*db.Driver, error)
	GetDriverByUserId(ctx context.Context, userID uuid.UUID) (*db.Driver, error)
	DriverOnline(ctx context.Context,userid uuid.UUID)error
}

type DriverRepositoryImpl struct {
	queries *db.Queries
}

func NewDriverRepository(db *db.Queries) *DriverRepositoryImpl {
	return &DriverRepositoryImpl{
		queries: db,
	}
}

func (r *DriverRepositoryImpl) CreateDriver(ctx context.Context, params db.CreateDriverParams) (*db.Driver, error) {
	driver, err := r.queries.CreateDriver(ctx, params)
	if err != nil {
		return nil, err
	}
	return &driver, nil
}

func (r *DriverRepositoryImpl) GetDriverByUserId(ctx context.Context, userID uuid.UUID) (*db.Driver, error) {
	driver, err := r.queries.GetDriverByUserId(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &driver, nil
}

func (r *DriverRepositoryImpl) DriverOnline(ctx context.Context, userid uuid.UUID) error {
	if err := r.queries.GoOnline(ctx, userid); err != nil {
		return err
	}
	return nil
}

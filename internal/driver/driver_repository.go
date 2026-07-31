package driver

import (
	"context"

	"github.com/Mikhail-Tal63/Orbit/internal/db"
)

type DriverRepository struct {
	queries *db.Queries
}

func NewDriverRepository(db *db.Queries)*DriverRepository{
	return &DriverRepository{
		queries:db,
	}
}

func (r *DriverRepository) CreateDriver(ctx context.Context,params db.CreateDriverParams)(*db.Driver,error){
	driver,err := r.queries.CreateDriver(ctx,params);
	if err!= nil {
return nil,err
	}
	return &driver,nil
}
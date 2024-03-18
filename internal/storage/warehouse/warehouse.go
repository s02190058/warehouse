package warehouse

import (
	"context"
	"errors"

	"github.com/cristalhq/builq"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	entity "github.com/s02190058/warehouse/internal/entity/warehouse"
	model "github.com/s02190058/warehouse/internal/model/warehouse"
	service "github.com/s02190058/warehouse/internal/service/warehouse"
	"github.com/s02190058/warehouse/pkg/db/postgres"
)

const (
	tableWarehouses        = "warehouses"
	tableProducts          = "products"
	tableWarehouseProducts = "warehouse_products"
)

type Storage struct {
	db *postgres.Database
}

func New(db *postgres.Database) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) Get(ctx context.Context, id int) (model.Warehouse, error) {
	q := builq.New()

	q("SELECT id, name, is_available")
	q("FROM %s", tableWarehouses)
	q("WHERE id = %$", id)

	sql, args, err := q.Build()
	if err != nil {
		return model.Warehouse{}, service.ErrInternal
	}

	var wh model.Warehouse
	err = s.db.Query(ctx).QueryRow(ctx, sql, args...).Scan(&wh.ID, &wh.Name, &wh.IsAvailable)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return model.Warehouse{}, service.ErrWarehouseNotFound
		default:
			return model.Warehouse{}, service.ErrInternal
		}
	}

	return wh, nil
}

func (s *Storage) Remains(ctx context.Context, id int) ([]entity.ProductRemains, error) {
	q := builq.New()

	q("SELECT p.code, wp.quantity - wp.reserved")
	q("FROM %s AS wp", tableWarehouseProducts)
	q("JOIN %s AS p ON wp.product_id = p.id", tableProducts)
	q("WHERE warehouse_id = %$", id)

	sql, args, err := q.Build()
	if err != nil {
		return nil, service.ErrInternal
	}

	rows, err := s.db.Query(ctx).Query(ctx, sql, args...)
	if err != nil {
		return nil, service.ErrInternal
	}
	defer rows.Close()

	var remains []entity.ProductRemains
	for rows.Next() {
		var remain entity.ProductRemains
		if err = rows.Scan(&remain.Code, &remain.Remains); err != nil {
			return nil, service.ErrInternal
		}

		remains = append(remains, remain)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return remains, nil
}

func (s *Storage) Reserve(ctx context.Context, id int, productCodes []string) ([]string, error) {
	q := builq.New()

	q("UPDATE %s AS wp", tableWarehouseProducts)
	q("SET reserved = reserved + 1")
	q("FROM %s AS p", tableProducts)
	q("WHERE wp.product_id = p.id")
	q("AND wp.warehouse_id = %$", id)
	q("AND p.code IN (%+$)", productCodes)
	q("RETURNING p.code")

	sql, args, err := q.Build()
	if err != nil {
		return nil, service.ErrInternal
	}

	rows, err := s.db.Query(ctx).Query(ctx, sql, args...)
	if err != nil {
		return nil, service.ErrInternal
	}
	defer rows.Close()

	var reservedCodes []string
	for rows.Next() {
		var code string
		if err = rows.Scan(&code); err != nil {
			return nil, service.ErrInternal
		}

		reservedCodes = append(reservedCodes, code)
	}

	if err = rows.Err(); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23514" {
			// check_violation
			return nil, service.ErrProductIsFullyReserved
		}

		return nil, service.ErrInternal
	}

	return reservedCodes, nil
}

func (s *Storage) Release(ctx context.Context, id int, productCodes []string) ([]string, error) {
	q := builq.New()

	q("UPDATE %s AS wp", tableWarehouseProducts)
	q("SET reserved = reserved - 1")
	q("FROM %s AS p", tableProducts)
	q("WHERE wp.product_id = p.id")
	q("AND wp.warehouse_id = %$", id)
	q("AND p.code IN (%+$)", productCodes)
	q("RETURNING p.code")

	sql, args, err := q.Build()
	if err != nil {
		return nil, service.ErrInternal
	}

	rows, err := s.db.Query(ctx).Query(ctx, sql, args...)
	if err != nil {
		return nil, service.ErrInternal
	}
	defer rows.Close()

	var reservedCodes []string
	for rows.Next() {
		var code string
		if err = rows.Scan(&code); err != nil {
			return nil, service.ErrInternal
		}

		reservedCodes = append(reservedCodes, code)
	}

	if err = rows.Err(); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23514" {
			// check_violation
			return nil, service.ErrProductHasNoReserves
		}

		return nil, service.ErrInternal
	}

	return reservedCodes, nil
}

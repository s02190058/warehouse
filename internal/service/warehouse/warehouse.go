package warehouse

import (
	"context"

	entity "github.com/s02190058/warehouse/internal/entity/warehouse"
	model "github.com/s02190058/warehouse/internal/model/warehouse"
	"github.com/s02190058/warehouse/internal/service"
)

type Storage interface {
	Get(ctx context.Context, id int) (wh model.Warehouse, err error)
	Remains(ctx context.Context, id int) (productsQuantity []entity.ProductRemains, err error)
	Reserve(ctx context.Context, id int, productCodes []string) (reservedCodes []string, err error)
	Release(ctx context.Context, id int, productCodes []string) (releasedCodes []string, err error)
}

type Service struct {
	transactor service.Transactor
	storage    Storage
}

func New(storage Storage, transactor service.Transactor) *Service {
	return &Service{
		transactor: transactor,
		storage:    storage,
	}
}

func (s *Service) Remains(ctx context.Context, id int) (productsQuantity []entity.ProductRemains, err error) {
	wh, err := s.storage.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if !wh.IsAvailable {
		return nil, ErrWarehouseNotAvailable
	}

	remains, err := s.storage.Remains(ctx, id)
	if err != nil {
		return nil, err
	}

	return remains, nil
}

func (s *Service) Reserve(ctx context.Context, id int, productCodes []string) ([]string, error) {
	reservedCodes := []string{}
	err := s.transactor.InTx(ctx, func(ctx context.Context) error {
		wh, err := s.storage.Get(ctx, id)
		if err != nil {
			return err
		}

		if !wh.IsAvailable {
			return ErrWarehouseNotAvailable
		}

		if len(productCodes) == 0 {
			return nil
		}

		reservedCodes, err = s.storage.Reserve(ctx, id, productCodes)

		return err
	})

	return reservedCodes, err
}

func (s *Service) Release(ctx context.Context, id int, productCodes []string) ([]string, error) {
	releasedCodes := []string{}
	err := s.transactor.InTx(ctx, func(ctx context.Context) error {
		wh, err := s.storage.Get(ctx, id)
		if err != nil {
			return err
		}

		if !wh.IsAvailable {
			return ErrWarehouseNotAvailable
		}

		if len(productCodes) == 0 {
			return nil
		}

		releasedCodes, err = s.storage.Release(ctx, id, productCodes)

		return err
	})

	return releasedCodes, err
}

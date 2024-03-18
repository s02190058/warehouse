package warehouse

import "errors"

var (
	ErrWarehouseNotFound      = errors.New("warehouse not found")
	ErrWarehouseNotAvailable  = errors.New("warehouse not available")
	ErrProductIsFullyReserved = errors.New("product is fully reserved")
	ErrProductHasNoReserves   = errors.New("product has no reserves")

	ErrInternal = errors.New("internal warehouse service error")
)

package warehouse

import "errors"

var (
	ErrBadJSONBody = errors.New("cannot decode json body")
	ErrBadURLParam = errors.New("bad url parameter")

	ErrInternal = errors.New("internal service error")
)

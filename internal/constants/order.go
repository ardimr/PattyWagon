package constants

import "errors"

var (
	ErrInvalidStartingPoint = errors.New("invalid starting point")
	ErrMerchantTooFar       = errors.New("merchant is too far")
	ErrItemNotFound         = errors.New("item is not found")
)

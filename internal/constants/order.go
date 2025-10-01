package constants

import "errors"

var (
	ErrInvalidStartingPoint = errors.New("invalid starting point")
	ErrMerchantTooFar       = errors.New("merchant is too far")
	ErrMerchantNotFound     = errors.New("merchant is not found")
	ErrItemNotFound         = errors.New("item is not found")
)

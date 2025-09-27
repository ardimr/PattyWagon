package constants

import "errors"

var (
	ErrInvalidStartingPoint = errors.New("Invalid starting point")
	ErrMerchantTooFar       = errors.New("Merchant is too far")
	ErrMerchantNotFound     = errors.New("Merchant is not  found")
)

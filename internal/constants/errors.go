package constants

import (
	"errors"
	"fmt"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserWrongPassword = errors.New("wrong password")
	ErrDuplicate         = errors.New("duplicate items")
	ErrDuplicatePhoneNum = errors.New("phone number already exists")
	ErrDuplicateEmail    = errors.New("email already exists")

	ErrInvalidFileType = errors.New("invalid file type")
	ErrMaximumFileSize = errors.New("size exceeds the maximum allowed file size")
	ErrFileNotFound    = errors.New("file not found")
	ErrInternalServer  = errors.New("internal server error")

	ErrFileIDNotValid                 = errors.New("fileId is not valid / exists")
	ErrDuplicateSKU                   = errors.New("duplicate sku")
	ErrProductNotFound                = errors.New("productId is not found")
	ErrInvalidRequest                 = errors.New("invalid request")
	ErrInvalidFileID                  = errors.New("invalid file ID")
	ErrEmailOrPhoneMustBeProvided     = errors.New("either email or phone must be provided")
	ErrCannotUpdateEmailAndPhone      = errors.New("cannot update both email and phone simultaneously")
	ErrInvalidEmailFormat             = errors.New("invalid email format")
	ErrInvalidPhoneNumberFormat       = errors.New("invalid phone number format")
	ErrNotEqualAvailableSellersInCart = errors.New("not equal to the available sellers in the cart")
	ErrPurchaseNotFound               = errors.New("purchase not found")

	// Order related errors
	ErrOrderNotFound       = errors.New("order not found")
	ErrOrderDetailNotFound = errors.New("order detail not found")
	ErrOrderItemNotFound   = errors.New("order item not found")
	ErrNoUnpurchasedOrder  = errors.New("no unpurchased order found")
	ErrMerchantNotFound    = errors.New("merchant not found")
	ErrTransactionFailed   = errors.New("transaction failed")

	// Service operation errors
	ErrFailedToBeginTransaction       = errors.New("failed to begin transaction")
	ErrFailedToGetMerchant            = errors.New("failed to get merchant")
	ErrFailedToGetItem                = errors.New("failed to get item")
	ErrFailedToCreateUnpurchasedOrder = errors.New("failed to create unpurchased order")
	ErrFailedToGetUnpurchasedOrder    = errors.New("failed to get unpurchased order")
	ErrFailedToCreateOrderDetail      = errors.New("failed to create order detail")
	ErrFailedToGetOrderDetail         = errors.New("failed to get order detail")
	ErrFailedToCreateOrderItem        = errors.New("failed to create order item")
	ErrFailedToGetOrderItem           = errors.New("failed to get order item")
	ErrFailedToUpdateOrderItem        = errors.New("failed to update order item")
	ErrFailedToCommitTransaction      = errors.New("failed to commit transaction")
	ErrFailedToCreateOrder            = errors.New("failed to create order")
	ErrFailedToGetOrder               = errors.New("failed to get order")
	ErrFailedToUpdateOrder            = errors.New("failed to update order")

	// File upload errors
	ErrErrorUploadingOriginalFile   = errors.New("error uploading original file")
	ErrErrorUploadingCompressedFile = errors.New("error uploading compressed file")
	ErrErrorInsertingFileToDatabase = errors.New("error inserting file to database")

	// Route finder errors
	ErrNotImplemented = errors.New("feature not implemented")

	// Order estimation errors
	ErrOrderEstimationNotFound        = errors.New("order estimation not found")
	ErrFailedToCalculateRoute         = errors.New("failed to calculate route")
	ErrFailedToProcessRoute           = errors.New("failed to process route")
	ErrFailedToCreateEstimation       = errors.New("failed to create estimation")
	ErrFailedToCreateEstimationItem   = errors.New("failed to create estimation item")
	ErrFailedToGetEstimation          = errors.New("failed to get estimation")
	ErrFailedToGetEstimations         = errors.New("failed to get estimations")
	ErrFailedToGetEstimationItems     = errors.New("failed to get estimation items")
	ErrFailedToUpdateEstimationStatus = errors.New("failed to update estimation status")
	ErrInvalidEstimationStatus        = errors.New("invalid estimation status")
)

// WrapError wraps an error with a base error for consistent error handling
func WrapError(baseErr error, wrappedErr error) error {
	return fmt.Errorf("%w: %w", baseErr, wrappedErr)
}

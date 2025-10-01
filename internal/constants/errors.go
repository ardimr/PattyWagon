package constants

import "errors"

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
	ErrMerchantNotFound               = errors.New("merchant not found")
)

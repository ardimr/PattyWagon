package server

type LoginRequest struct {
	Username string `json:"username" validate:"required,min=5,max=30"`
	Password string `json:"password" validate:"required,min=5,max=30"`
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=5,max=30"`
	Password string `json:"password" validate:"required,min=5,max=30"`
	Email    string `json:"email" validate:"required,email"`
}

type CreateMerchantRequest struct {
	Name     string         `json:"name" validate:"required,min=2,max=30"`
	Category string         `json:"category" validate:"required,merchantCategory"`
	ImageURL string         `json:"image_url" validate:"required"`
	Location DetailLocation `json:"location" validate:"required"`
}

type DetailLocation struct {
	Latitude  float64 `json:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" validate:"required"`
}

type GetMerchantRequest struct {
	MerchantID       string `query:"merchantId"`
	Limit            string `query:"limit"`
	Offset           string `query:"offset"`
	Name             string `query:"name"`
	MerchantCategory string `query:"merchantCategory"`
	CreatedAt        string `query:"createdAt"`
}

type CreateItemRequest struct {
	Name            string  `json:"name" validate:"required,min=2,max=30"`
	ProductCategory string  `json:"productCategory" validate:"required,productCategory"`
	Price           float64 `json:"price" validate:"required,gt=0"`
	ImageURL        string  `json:"imageUrl" validate:"required,url"`
}

type GetItemsRequest struct {
	ItemID          string `query:"itemId"`
	Limit           string `query:"limit"`
	Offset          string `query:"offset"`
	Name            string `query:"name"`
	ProductCategory string `query:"productCategory"`
	CreatedAt       string `query:"createdAt"`
}

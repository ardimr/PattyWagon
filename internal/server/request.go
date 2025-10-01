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

package server

import (
	"PattyWagon/internal/model"
	"PattyWagon/internal/types"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	constants "PattyWagon/internal/constants"

	"github.com/go-playground/validator/v10"

	_ "github.com/joho/godotenv/autoload"
)

type Service interface {
	UsernameLogin(ctx context.Context, username string, password string, role int16) (token string, err error)
	Register(ctx context.Context, userReq model.User, password string, role int16) (string, error)

	UploadFile(ctx context.Context, file io.Reader, filename string, sizeInBytes int64) (model.File, error)

	CreateMerchant(ctx context.Context, req model.Merchant) (res int64, err error)
	GetMerchants(ctx context.Context, req model.FilterMerchant) (res []model.Merchant, err error)

	CreateItems(ctx context.Context, req model.Item) (res int64, err error)
	GetItems(ctx context.Context, req model.FilterItem) (res []model.Item, err error)

	// Purchase
	// EstimateOrderPrice(ctx context.Context, req model.OrderEstimation) (model.EstimationPrice, error)
	FindNearbyMerchants(ctx context.Context, userLocation model.Location, searchParams model.FindNerbyMerchantParams) ([]model.MerchantItem, error)

	GetOrder(ctx context.Context, orderID int64) (model.Order, error)
	GetOrderDetail(ctx context.Context, orderDetailID int64) (model.OrderDetail, error)
	AddItemToCart(ctx context.Context, userID, merchantID, itemID int64, quantity int32) (model.OrderDetail, error)
	CreateUnpurchasedOrder(ctx context.Context, userID int64) (model.Order, error)
	GetPurchasedOrdersPagination(ctx context.Context, userID int64, limit, offset int) ([]model.Order, error)
	UpdateOrderToPurchased(ctx context.Context, orderID int64) error

	GetMerchant(ctx context.Context, merchantID int64) (model.Merchant, error)
	GetItem(ctx context.Context, itemID int64) (model.Item, error)

	FindOptimalRoute(ctx context.Context, startLat, startLon, userLat, userLon float64, merchantIDs []int64) (types.RouteResult, error)
	CalculateEstimatedDeliveryTime(ctx context.Context, route types.RouteResult) (float64, error)

	CreateOrderEstimation(ctx context.Context, req model.EstimationRequest) (model.EstimationResponse, error)
	GetMerchantItemDataConcurrently(ctx context.Context, items []model.EstimationRequestItem) ([]model.MerchantItemData, error)
}

type Server struct {
	port      int
	service   Service
	validator *validator.Validate
}

func NewServer(service Service) *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	v := validator.New()
	NewServer := &Server{
		port:      port,
		service:   service,
		validator: v,
	}

	// Custom validator for merchant category
	NewServer.validator.RegisterValidation("merchantCategory", func(fl validator.FieldLevel) bool {
		return constants.IsValidMerchantCategory(fl.Field().String())
	})
	// Custom validator for product category
	NewServer.validator.RegisterValidation("productCategory", func(fl validator.FieldLevel) bool {
		return constants.IsValidProductCategory(fl.Field().String())
	})

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}

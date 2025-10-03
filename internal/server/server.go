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

	"github.com/go-playground/validator/v10"

	_ "github.com/joho/godotenv/autoload"
)

type Service interface {
	EmailLogin(ctx context.Context, user string, password string) (string, string, error)
	Register(ctx context.Context, user model.User, password string) (string, error)

	IsUserExist(ctx context.Context, userID int64) (bool, error)

	UploadFile(ctx context.Context, file io.Reader, filename string, sizeInBytes int64) (model.File, error)

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

	// Custom validator for product type
	// NewServer.validator.RegisterValidation("productType", func(fl validator.FieldLevel) bool {
	// 	productType := fl.Field().String()
	// 	for _, pt := range constants.ProductTypes {
	// 		if strings.EqualFold(pt, productType) {
	// 			return true
	// 		}
	// 	}
	// 	return false
	// })

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

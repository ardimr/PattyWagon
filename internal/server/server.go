package server

import (
	"PattyWagon/internal/model"
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
	NewServer.validator.RegisterValidation("merchantCategory", func(fl validator.FieldLevel) bool {
		return constants.IsValidMerchantCategory(fl.Field().String())
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

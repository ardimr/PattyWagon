package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEstimateOrderPrice(t *testing.T) {
	s := testSetup(t)

	validReq := OrderEstimationRequest{
		UserLocation: LocationRequest{
			Lat:  6.1674,
			Long: 106.8209,
		},

		Orders: []OrderRequest{
			{
				MerchantID:      "1",
				IsStartingPoint: true,
				Items: []OrderItemRequest{
					{ItemID: "1", Quantity: 1},
					{ItemID: "2", Quantity: 2},
				},
			},
			{
				MerchantID:      "1",
				IsStartingPoint: false,
				Items: []OrderItemRequest{
					{ItemID: "3", Quantity: 1},
					{ItemID: "4", Quantity: 2},
				},
			},
		},
	}

	invalidStartingPointReq := OrderEstimationRequest{
		UserLocation: LocationRequest{
			Lat:  6.1674,
			Long: 106.8209,
		},

		Orders: []OrderRequest{
			{
				MerchantID:      "1",
				IsStartingPoint: true,
				Items: []OrderItemRequest{
					{ItemID: "1", Quantity: 1},
					{ItemID: "2", Quantity: 2},
				},
			},
			{
				MerchantID:      "2",
				IsStartingPoint: true,
				Items: []OrderItemRequest{
					{ItemID: "3", Quantity: 1},
					{ItemID: "4", Quantity: 2},
				},
			},
		},
	}

	invaliMerchantTooFarReq := OrderEstimationRequest{
		UserLocation: LocationRequest{
			Lat:  6.1674,
			Long: 106.8209,
		},

		Orders: []OrderRequest{
			{
				MerchantID:      "1",
				IsStartingPoint: true,
				Items: []OrderItemRequest{
					{ItemID: "1", Quantity: 1},
					{ItemID: "2", Quantity: 2},
				},
			},
			{
				MerchantID:      "99",
				IsStartingPoint: false,
				Items: []OrderItemRequest{
					{ItemID: "3", Quantity: 1},
					{ItemID: "4", Quantity: 2},
				},
			},
		},
	}
	t.Run("Valid", func(t *testing.T) {
		reqBody, err := json.Marshal(validReq)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, "/v1/users/estimate", bytes.NewBuffer(reqBody))
		w := httptest.NewRecorder()

		s.EstimateOrderPrice(w, req)

		resp := w.Result()
		assert.NotEqual(t, 0, resp.StatusCode)
		fmt.Printf("Response Body: %v\n", w.Body.String())
	})

	t.Run("Invalid_StartingPoint", func(t *testing.T) {
		reqBody, err := json.Marshal(invalidStartingPointReq)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, "/v1/users/estimate", bytes.NewBuffer(reqBody))
		w := httptest.NewRecorder()

		s.EstimateOrderPrice(w, req)

		resp := w.Result()
		assert.Equal(t, 400, resp.StatusCode)
		fmt.Printf("Response Body: %v\n", w.Body.String())
	})

	t.Run("Invalid_MerchantTooFar", func(t *testing.T) {
		reqBody, err := json.Marshal(invaliMerchantTooFarReq)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, "/v1/users/estimate", bytes.NewBuffer(reqBody))
		w := httptest.NewRecorder()

		s.EstimateOrderPrice(w, req)

		resp := w.Result()
		assert.Equal(t, 400, resp.StatusCode)
		fmt.Printf("Response Body: %v\n", w.Body.String())
	})
}

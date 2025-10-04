package repository

import (
	"PattyWagon/internal/database"
	"PattyWagon/internal/model"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func setupRepo(t *testing.T) *Queries {
	t.Helper()

	db := database.New(
		"localhost",
		"5432",
		"patty-wagon-dev",
		"postgres",
		"postgres",
		"public",
		&database.ConnectionPoolConfig{
			MaxOpenConns:    int(database.MaxOpenConns),
			MaxIdleConns:    int(database.MaxIdleConns),
			ConnMaxIdleTime: time.Duration(database.ConnMaxIdleTime * int64(time.Second)),
			ConnMaxLifeTime: time.Duration(database.ConnMaxLifeTime * int64(time.Second)),
		},
	)

	repo := New(db)
	return repo
}

func TestListMerchantWitItems(t *testing.T) {
	repo := setupRepo(t)
	t.Run("Valid", func(t *testing.T) {
		filter := model.ListMerchantWithItemParams{
			Cell: &model.Cell{
				CellID:     610049360213835775,
				Resolution: 8,
			},
			MerchantParams: model.MerchantParams{
				// Name: "",
			},
		}
		merchantItems, err := repo.ListMerchantWithItems(context.TODO(), filter)
		assert.Nil(t, err)
		assert.NotEmpty(t, merchantItems)

		for i, merchant := range merchantItems {
			t.Logf("merchant %d: %v", i, merchant)
		}
	})

	t.Run("Valid_WithoutCell", func(t *testing.T) {
		filter := model.ListMerchantWithItemParams{
			MerchantParams: model.MerchantParams{
				// Name: "",
			},
		}
		merchantItems, err := repo.ListMerchantWithItems(context.TODO(), filter)
		if err != nil {
			fmt.Println(err)
		}

		assert.Nil(t, err)
		assert.NotEmpty(t, merchantItems)

		for i, merchant := range merchantItems {
			t.Logf("merchant %d: %v", i, merchant.Merchant.Name)
		}
	})

	t.Run("Valid_WithFilterName", func(t *testing.T) {
		filterName := "bat"
		filter := model.ListMerchantWithItemParams{
			Cell: &model.Cell{
				CellID: 614348827586985983,
			},
			MerchantParams: model.MerchantParams{
				Name: &filterName,
			},
		}

		merchantItems, err := repo.ListMerchantWithItems(context.TODO(), filter)
		assert.Nil(t, err)
		assert.NotEmpty(t, merchantItems)

		for i, merchant := range merchantItems {
			t.Logf("merchant %d: %v", i, merchant)
		}
	})

	t.Run("Valid_WithWrongName", func(t *testing.T) {
		filterName := "kopi"
		filter := model.ListMerchantWithItemParams{
			MerchantParams: model.MerchantParams{
				Name: &filterName,
			},
		}

		merchantItems, err := repo.ListMerchantWithItems(context.TODO(), filter)
		assert.Nil(t, err)
		assert.Empty(t, merchantItems)
	})

	t.Run("Valid_WithCategory", func(t *testing.T) {
		filterCategory := "makanan"
		filter := model.ListMerchantWithItemParams{
			MerchantParams: model.MerchantParams{
				MerchantCategory: &filterCategory,
			},
		}

		merchantItems, err := repo.ListMerchantWithItems(context.TODO(), filter)
		assert.Nil(t, err)
		assert.NotEmpty(t, merchantItems)
		for i, merchant := range merchantItems {
			t.Logf("merchant %d: %v", i, merchant.Merchant.Name)
		}
	})
}

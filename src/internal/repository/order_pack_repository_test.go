package repository_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/hayrullahcansu/order-packs-calculator/src/internal/model"
	"github.com/hayrullahcansu/order-packs-calculator/src/internal/repository"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB creates a fresh in-memory SQLite database with migrated schema.
func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	assert.NoError(t, db.AutoMigrate(&model.OrderPack{}))
	return db
}

func TestAddAndFetchOrderPacks(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewOrderPackRepository(db)
	ctx := context.Background()

	// add two packs
	assert.NoError(t, repo.AddOrderPack(ctx, &model.OrderPack{Items: 250}))
	assert.NoError(t, repo.AddOrderPack(ctx, &model.OrderPack{Items: 500}))

	// fetch all
	packs, err := repo.FetchAvailableOrderPacks(ctx)
	assert.NoError(t, err)
	assert.Len(t, packs, 2)
}

func TestAddDuplicateItems_ReturnsError(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewOrderPackRepository(db)
	ctx := context.Background()

	assert.NoError(t, repo.AddOrderPack(ctx, &model.OrderPack{Items: 250}))

	// duplicate items should fail due to unique index
	err := repo.AddOrderPack(ctx, &model.OrderPack{Items: 250})
	assert.Error(t, err)
}

func TestGetOrderPackByID(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewOrderPackRepository(db)
	ctx := context.Background()

	pack := &model.OrderPack{Items: 1000}
	assert.NoError(t, repo.AddOrderPack(ctx, pack))

	// fetch by the generated UUID
	found, err := repo.GetOrderPackByID(ctx, pack.ID)
	assert.NoError(t, err)
	assert.Equal(t, 1000, found.Items)
	assert.Equal(t, pack.ID, found.ID)
}

func TestGetOrderPackByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewOrderPackRepository(db)
	ctx := context.Background()

	_, err := repo.GetOrderPackByID(ctx, uuid.New())
	assert.Error(t, err)
}

func TestUpdateOrderPack(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewOrderPackRepository(db)
	ctx := context.Background()

	pack := &model.OrderPack{Items: 250}
	assert.NoError(t, repo.AddOrderPack(ctx, pack))

	// update items
	pack.Items = 500
	assert.NoError(t, repo.UpdateOrderPack(ctx, pack))

	// verify the update
	updated, err := repo.GetOrderPackByID(ctx, pack.ID)
	assert.NoError(t, err)
	assert.Equal(t, 500, updated.Items)
}

func TestRemoveOrderPack(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewOrderPackRepository(db)
	ctx := context.Background()

	pack := &model.OrderPack{Items: 250}
	assert.NoError(t, repo.AddOrderPack(ctx, pack))

	// remove
	assert.NoError(t, repo.RemoveOrderPack(ctx, pack.ID))

	// should no longer exist
	_, err := repo.GetOrderPackByID(ctx, pack.ID)
	assert.Error(t, err)
}

func TestFetchAvailableOrderPacks_Empty(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewOrderPackRepository(db)
	ctx := context.Background()

	packs, err := repo.FetchAvailableOrderPacks(ctx)
	assert.NoError(t, err)
	assert.Empty(t, packs)
}

func TestFullCRUDCycle(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewOrderPackRepository(db)
	ctx := context.Background()

	// create
	pack := &model.OrderPack{Items: 2000}
	assert.NoError(t, repo.AddOrderPack(ctx, pack))

	// read
	found, err := repo.GetOrderPackByID(ctx, pack.ID)
	assert.NoError(t, err)
	assert.Equal(t, 2000, found.Items)

	// update
	found.Items = 3000
	assert.NoError(t, repo.UpdateOrderPack(ctx, found))

	// verify update
	updated, err := repo.GetOrderPackByID(ctx, pack.ID)
	assert.NoError(t, err)
	assert.Equal(t, 3000, updated.Items)

	// delete
	assert.NoError(t, repo.RemoveOrderPack(ctx, pack.ID))

	// verify delete
	_, err = repo.GetOrderPackByID(ctx, pack.ID)
	assert.Error(t, err)

	// list should be empty
	packs, err := repo.FetchAvailableOrderPacks(ctx)
	assert.NoError(t, err)
	assert.Empty(t, packs)
}

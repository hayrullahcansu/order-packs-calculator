package db

import (
	"fmt"
	"time"

	"github.com/hayrullahcansu/order-packs-calculator/src/internal/model"
	"github.com/hayrullahcansu/order-packs-calculator/src/shared/logging"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// InitiateSqliteDbContext creates an sqlite3 database context
// File: order-pack.db
func InitiateSqliteDbContext() *gorm.DB {
	connectionString := "order-pack.db"
	return initateDbContext(connectionString)
}

// InitiateSqliteDbContext creates an sqlite3 database context in memory
func InitiateInmemorySqliteDbContext() *gorm.DB {
	connectionString := fmt.Sprintf("file:memdb_%d?mode=memory&cache=shared&_busy_timeout=2000", time.Now().UnixNano())
	return initateDbContext(connectionString)
}

func initateDbContext(connection string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(connection), &gorm.Config{})
	if err != nil {
		logging.Fatalf("db couldn't initiated %v", err)
	}
	logging.Info("database has initiated")
	migrateErr := db.AutoMigrate(&model.OrderPack{})
	if migrateErr != nil {
		logging.Fatalf("migration error %v", migrateErr)
	}
	logging.Info("migration has completed")

	// seed default order packs on first run
	var count int64
	db.Model(&model.OrderPack{}).Count(&count)
	if count == 0 {
		seedPacks := []model.OrderPack{
			{Items: 23},
			{Items: 31},
			{Items: 53},
		}
		db.Create(&seedPacks)
		logging.Info("seed data has been inserted")
	}

	return db
}

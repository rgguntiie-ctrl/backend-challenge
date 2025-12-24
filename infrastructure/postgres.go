package infrastructure

import (
	"fmt"
	"log"
	"time"

	"github.com/kanta/backend-challenge/internal/adapters/repositories/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresClient(host, user, password, dbname, port string) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port,
	)

	var db *gorm.DB
	var err error

	for i := 0; i < 3; i++ { // retry 3 times
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("failed to connect to postgres (attempt %d): %v\n", i+1, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatalf("failed to connect to postgres after retries: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get sql.DB from gorm: %v", err)
	}

	// setting connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := db.AutoMigrate(models.GetAllModels()...); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	log.Println("Connected to PostgreSQL successfully")
	return db
}

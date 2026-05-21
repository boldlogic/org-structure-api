package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase(dsn string) (*gorm.DB, error) {
	db, err := initializeDatabase(dsn)
	if err != nil {
		return nil, err
	}
	if err := CheckDatabaseConnection(db); err != nil {
		return nil, err
	}
	return db, nil
}

func initializeDatabase(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к БД: %w", err)
	}
	return db, nil
}

func CheckDatabaseConnection(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("не удалось получить экземпляр БД: %w", err)
	}
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("потеря БД: %w", err)
	}
	return nil
}

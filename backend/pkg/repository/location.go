package repository

import (
	"context"
	"fmt"
	"sync"

	"github.com/wakka-2/Namless/backend/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Location models the DB operations available for locations.
type Location struct {
	db    *gorm.DB
	mutex sync.RWMutex
}

// NewLocation builds a new Location repository.
//
// When silent is true, it will use a custom logger that does not output anything to the console.
func NewLocation(dsn string, silent bool) (*Location, error) {
	result := &Location{}

	cfg := &gorm.Config{}
	if silent {
		cfg.Logger = NewNoopLogger()
	}

	var err error

	result.db, err = gorm.Open(postgres.Open(dsn), cfg)
	if err != nil {
		return nil, fmt.Errorf("could not open Import DB: %w", err)
	}

	err = result.db.AutoMigrate(&models.Location{})
	if err != nil {
		return nil, fmt.Errorf("could not auto migrate models.Data: %w", err)
	}

	return result, nil
}

// NewLocationTruncate builds a new Location repo, and deletes its previous contents.
//
// Meant to be used in tests.
func NewLocationTruncate(dsn string, silent bool) (*Location, error) {
	result, err := NewLocation(dsn, silent)

	if err == nil {
		success := result.db.Exec("TRUNCATE TABLE locations;")
		if success.Error != nil {
			return nil, fmt.Errorf("could not truncate: %w", err)
		}
	}

	return result, err
}

// GetAll returns all Location items.
func (l *Location) GetAll(ctx context.Context) ([]models.Location, error) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	var result []models.Location

	err := l.db.WithContext(ctx).Find(&result).Error
	if err != nil {
		return nil, fmt.Errorf("could not get all import items: %w", err)
	}

	return result, nil
}

// Create a new Location item.
func (l *Location) Create(ctx context.Context, item models.Location) (models.Location, error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	success := l.db.WithContext(ctx).Create(&item)
	if success.Error != nil {
		return models.Location{}, fmt.Errorf("could not create Location item: %w", success.Error)
	}

	return item, nil
}

// Update a given data item.
func (l *Location) Update(ctx context.Context, item models.Location) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if item.ID < 0 {
		return ErrDoesNotExist
	}

	var result models.Data

	success := l.db.WithContext(ctx).First(&result, "id = ?", item.ID)
	if success.Error != nil {
		return ErrDoesNotExist
	}

	success = l.db.WithContext(ctx).Save(item)
	if success.Error != nil {
		return fmt.Errorf("could not update data item: %w", success.Error)
	}

	return nil
}

// ByID returns the data item with a given ID.
func (l *Location) ByID(ctx context.Context, itemID int) (models.Location, error) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	var result models.Location

	success := l.db.WithContext(ctx).First(&result, "id = ?", itemID)

	if success.Error != nil {
		return models.Location{}, fmt.Errorf("could not find Location item with ID %q: %w", itemID, success.Error)
	}

	return result, nil
}

// Delete a given Location item.
//
//nolint:dupl
func (l *Location) Delete(ctx context.Context, locationID int) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if locationID < 0 {
		return ErrDoesNotExist
	}

	var result models.Data

	success := l.db.WithContext(ctx).First(&result, "id = ?", locationID)
	if success.Error != nil {
		return ErrDoesNotExist
	}

	success = l.db.WithContext(ctx).Delete(&result)
	if success.Error != nil {
		return fmt.Errorf("could not delete Location item %q: %w", locationID, success.Error)
	}

	return nil
}

// Close closes the DB connection.
func (l *Location) Close(ctx context.Context) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	database, err := l.db.WithContext(ctx).DB()
	if err != nil {
		return fmt.Errorf("could not get DB: %w", err)
	}

	err = database.Close()
	if err != nil {
		return fmt.Errorf("could not close DB: %w", err)
	}

	return nil
}

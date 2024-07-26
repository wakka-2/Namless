/*
Package repository offers objects that can do CRUD operations on the import table.
*/
package repository

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	// ErrDoesNotExist for when we try to update/delete a non existing search.
	ErrDoesNotExist = errors.New("item does not exit")
)

// Store models the DB operations available for search items.
type Store struct {
	db    *gorm.DB
	mutex sync.RWMutex
}

// New builds a new import repository.
//
// When silent is true, it will use a custom logger that does not output anything to the console.
func New(dsn string, silent bool) (*Store, error) {
	result := &Store{}

	cfg := &gorm.Config{}
	if silent {
		cfg.Logger = NewNoopLogger()
	}

	var err error

	result.db, err = gorm.Open(postgres.Open(dsn), cfg)
	if err != nil {
		return nil, fmt.Errorf("could not open Import DB: %w", err)
	}

	err = result.db.AutoMigrate(&models.Import{})
	if err != nil {
		return nil, fmt.Errorf("could not auto migrate models.Import: %w", err)
	}

	return result, nil
}

// NewTruncate builds a new search repo, and deletes its previous contents.
//
// Meant to be used in tests.
func NewTruncate(dsn string, silent bool) (*Store, error) {
	result, err := New(dsn, silent)

	if err == nil {
		success := result.db.Exec("TRUNCATE TABLE imports;")
		if success.Error != nil {
			return nil, fmt.Errorf("could not truncate: %w", err)
		}
	}

	return result, err
}

// GetAll returns all import items.
func (c *Store) GetAll(ctx context.Context) ([]models.Import, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	var result []models.Import

	err := c.db.WithContext(ctx).Find(&result).Error
	if err != nil {
		return nil, fmt.Errorf("could not get all import items: %w", err)
	}

	return result, nil
}

// Monitored returns all import items with status StatusMonitored, that were not previously deleted.
func (c *Store) Monitored(ctx context.Context) ([]models.Import, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	var result []models.Import

	err := c.db.WithContext(ctx).Where("status = ?", models.StatusMonitored).Find(&result).Error
	if err != nil {
		return nil, fmt.Errorf("could not get serch items: %w", err)
	}

	return result, nil
}

// Create a new search item.
//
// Sets the CreatedAt and UpdatedAt fields.
func (c *Store) Create(ctx context.Context, item models.Import) (models.Import, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	item.CreatedAt = time.Now()
	item.UpdatedAt = item.CreatedAt

	success := c.db.WithContext(ctx).Create(&item)
	if success.Error != nil {
		return models.Import{}, fmt.Errorf("could not create import item: %w", success.Error)
	}

	return item, nil
}

// Update a given search item.
func (c *Store) Update(ctx context.Context, item models.Import) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if item.ID == "" {
		return ErrDoesNotExist
	}

	var result models.Import

	success := c.db.WithContext(ctx).First(&result, "id = ?", item.ID)
	if success.Error != nil {
		return ErrDoesNotExist
	}

	item.UpdatedAt = time.Now()

	success = c.db.WithContext(ctx).Save(item)
	if success.Error != nil {
		return fmt.Errorf("could not update search item: %w", success.Error)
	}

	return nil
}

// ByID returns the search item with a given ID.
func (c *Store) ByID(ctx context.Context, itemID string) (models.Import, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	var result models.Import

	success := c.db.WithContext(ctx).First(&result, "id = ?", itemID)

	if success.Error != nil {
		return models.Import{}, fmt.Errorf("could not find import item with ID %q: %w", itemID, success.Error)
	}

	return result, nil
}

// Delete a given search item.
func (c *Store) Delete(ctx context.Context, importID string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if importID == "" {
		return ErrDoesNotExist
	}

	var result models.Import

	success := c.db.WithContext(ctx).First(&result, "id = ?", importID)
	if success.Error != nil {
		return ErrDoesNotExist
	}

	success = c.db.WithContext(ctx).Delete(&result)
	if success.Error != nil {
		return fmt.Errorf("could not delete import %q: %w", importID, success.Error)
	}

	return nil
}

// DeleteOld deletes all entries created before a given timestamp.
func (c *Store) DeleteOld(ctx context.Context, startingFrom time.Time) (int64, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	success := c.db.WithContext(ctx).Where("created_at < ?", startingFrom).Delete(&models.Import{})
	if success.Error != nil {
		return 0, fmt.Errorf("could not delete search items older than %q: %w",
			startingFrom.Format(time.DateTime), success.Error)
	}

	return success.RowsAffected, nil
}

// Close closes the DB connection.
func (c *Store) Close(ctx context.Context) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	database, err := c.db.WithContext(ctx).DB()
	if err != nil {
		return fmt.Errorf("could not get DB: %w", err)
	}

	err = database.Close()
	if err != nil {
		return fmt.Errorf("could not close DB: %w", err)
	}

	return nil
}

// ByProjectID returns the search item with a given project ID.
//
// Returns an error if a project with the given projectID does not exist.
func (c *Store) ByProjectID(ctx context.Context, projectID string) (models.Import, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	var result models.Import

	success := c.db.WithContext(ctx).Where("project_id = ?", projectID).First(&result)

	if success.Error != nil {
		return models.Import{}, fmt.Errorf("could not find import item with projectID %q: %w", projectID, success.Error)
	}

	return result, nil
}

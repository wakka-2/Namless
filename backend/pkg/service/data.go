/*
Package service offers data-related functionality.
*/
package service

import (
	"context"
	"fmt"

	"github.com/wakka-2/Namless/backend/pkg/models"
	"github.com/wakka-2/Namless/backend/pkg/repository"
	"github.com/wakka-2/Namless/backend/pkg/types"
)

// Data offers data-related functionality.
type Data struct {
	db        *repository.Store
	serverCtx context.Context
}

// New builds a new data service.
func New(ctx context.Context, db *repository.Store) *Data {
	return &Data{
		db:        db,
		serverCtx: ctx,
	}
}

// Add a new key-value pair.
func (d *Data) Add(ctx context.Context, key string, value string) error {
	if d.serverCtx.Err() != nil || ctx.Err() != nil {
		return types.ErrCancelledContext
	}

	_, err := d.db.Create(ctx, models.Data{
		ID:    key,
		Value: value,
	})

	if err != nil {
		return fmt.Errorf("could not create data entry: %w", err)
	}

	return nil
}

// Get the key-value pair with a given key.
func (d *Data) Get(ctx context.Context, key string) (string, error) {
	if d.serverCtx.Err() != nil || ctx.Err() != nil {
		return "", types.ErrCancelledContext
	}

	result, err := d.db.ByID(ctx, key)
	if err != nil {
		return "", fmt.Errorf("could not retrieve data entry: %w", err)
	}

	return result.Value, nil
}

// Update a given key-value pair.
func (d *Data) Update(ctx context.Context, key string, value string) error {
	if d.serverCtx.Err() != nil || ctx.Err() != nil {
		return types.ErrCancelledContext
	}

	err := d.db.Update(ctx, models.Data{
		ID:    key,
		Value: value,
	})

	if err != nil {
		return fmt.Errorf("could not update data entry: %w", err)
	}

	return nil
}

// Delete a given key-value pair.
func (d *Data) Delete(ctx context.Context, key string) error {
	if d.serverCtx.Err() != nil || ctx.Err() != nil {
		return types.ErrCancelledContext
	}

	err := d.db.Delete(ctx, key)
	if err != nil {
		return fmt.Errorf("could not delete data entry: %w", err)
	}

	return nil
}

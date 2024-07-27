package service

import (
	"context"
	"fmt"

	"github.com/wakka-2/Namless/backend/pkg/models"
	"github.com/wakka-2/Namless/backend/pkg/repository"
	"github.com/wakka-2/Namless/backend/pkg/types"
)

// Location offers Location-related functionality.
type Location struct {
	db        *repository.Location
	serverCtx context.Context
}

// NewLocation builds a new Location service.
func NewLocation(ctx context.Context, db *repository.Location) *Location {
	return &Location{
		db:        db,
		serverCtx: ctx,
	}
}

// Add a new key-value pair.
func (l *Location) Add(ctx context.Context, location models.Location) error {
	if l.serverCtx.Err() != nil || ctx.Err() != nil {
		return types.ErrCancelledContext
	}

	_, err := l.db.Create(ctx, location)

	if err != nil {
		return fmt.Errorf("could not create Location entry: %w", err)
	}

	return nil
}

// Get the location with a given ID.
func (l *Location) Get(ctx context.Context, id int) (models.Location, error) {
	if l.serverCtx.Err() != nil || ctx.Err() != nil {
		return models.Location{}, types.ErrCancelledContext
	}

	result, err := l.db.ByID(ctx, id)
	if err != nil {
		return models.Location{}, fmt.Errorf("could not retrieve Location entry: %w", err)
	}

	return result, nil
}

// Get the key-value pairs.
func (l *Location) GetAll(ctx context.Context) ([]models.Location, error) {
	if l.serverCtx.Err() != nil || ctx.Err() != nil {
		return nil, types.ErrCancelledContext
	}

	result, err := l.db.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve locations: %w", err)
	}

	return result, nil
}

// Update a given key-value pair.
func (l *Location) Update(ctx context.Context, location models.Location) error {
	if l.serverCtx.Err() != nil || ctx.Err() != nil {
		return types.ErrCancelledContext
	}

	err := l.db.Update(ctx, location)
	if err != nil {
		return fmt.Errorf("could not update Location entry: %w", err)
	}

	return nil
}

// Delete a given key-value pair.
func (l *Location) Delete(ctx context.Context, id int) error {
	if l.serverCtx.Err() != nil || ctx.Err() != nil {
		return types.ErrCancelledContext
	}

	err := l.db.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("could not delete Location entry: %w", err)
	}

	return nil
}

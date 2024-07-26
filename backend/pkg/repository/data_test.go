package repository

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wakka-2/Namless/backend/pkg/models"
	"github.com/wakka-2/Namless/backend/pkg/types"
)

func Test_Create(t *testing.T) {
	repo, err := buildRepo(true)
	assert.NoError(t, err)

	defer func() {
		err := repo.Close(context.TODO())
		assert.NoError(t, err)
	}()

	Import := models.Data{
		ID: "BBB171C4-00E8-4B0F-97EB-2F3EC3394A87",
	}

	got, err := repo.Create(context.TODO(), Import)
	assert.NoError(t, err)
	assert.Equal(t, "BBB171C4-00E8-4B0F-97EB-2F3EC3394A87", got.ID)
}

func Test_Update(t *testing.T) {
	repo, err := buildRepo(false)
	assert.NoError(t, err)

	defer func() {
		err := repo.Close(context.TODO())
		assert.NoError(t, err)
	}()

	Import := models.Data{
		ID: "AAA-00E8-4B0F-97EB-2F3EC3394A87",
	}

	created, err := repo.Create(context.TODO(), Import)
	assert.NoError(t, err)
	assert.Equal(t, "AAA-00E8-4B0F-97EB-2F3EC3394A87", created.ID)

	err = repo.Update(context.TODO(), created)
	assert.NoError(t, err)

	updated, err := repo.ByID(context.TODO(), "AAA-00E8-4B0F-97EB-2F3EC3394A87")
	assert.NoError(t, err)
	assert.Equal(t, "AAA-00E8-4B0F-97EB-2F3EC3394A87", updated.ID)
}

func Test_Delete(t *testing.T) {
	repo, err := buildRepo(true)
	assert.NoError(t, err)

	defer func() {
		err := repo.Close(context.TODO())
		assert.NoError(t, err)
	}()

	Import := models.Data{
		ID: "BBB-00E8-4B0F-97EB-2F3EC3394A87",
	}

	created, err := repo.Create(context.TODO(), Import)
	assert.NoError(t, err)
	assert.Equal(t, "BBB-00E8-4B0F-97EB-2F3EC3394A87", created.ID)

	found, err := repo.ByID(context.TODO(), created.ID)
	assert.NoError(t, err)
	assert.Equal(t, created.ID, found.ID)

	err = repo.Delete(context.TODO(), found.ID)
	assert.NoError(t, err)

	found, err = repo.ByID(context.TODO(), created.ID)
	assert.Error(t, err)
}

func Test_DeleteOld(t *testing.T) {
	repo, err := buildRepo(true)
	assert.NoError(t, err)

	defer func() {
		err := repo.Close(context.TODO())
		assert.NoError(t, err)
	}()

	creationTime := time.Date(2024, 4, 1, 1, 1, 1, 1, time.UTC)

	Imports := []models.Data{
		{
			ID:        "1BB-00E8-4B0F-97EB-2F3EC3394A87",
			CreatedAt: creationTime,
		},
		{
			ID:        "2BB-00E8-4B0F-97EB-2F3EC3394A87",
			CreatedAt: creationTime,
		},
		{
			ID:        "3BB-00E8-4B0F-97EB-2F3EC3394A87",
			CreatedAt: creationTime,
		},
		{
			ID:        "4BB-00E8-4B0F-97EB-2F3EC3394A87",
			CreatedAt: creationTime.Add(2 * time.Hour),
		},
	}

	for _, item := range Imports {
		created, err := repo.Create(context.TODO(), item)
		assert.NoError(t, err)
		assert.Equal(t, item.ID, created.ID)

		found, err := repo.ByID(context.TODO(), created.ID)
		assert.NoError(t, err)
		assert.Equal(t, created.ID, found.ID)

		found.CreatedAt = item.CreatedAt // create automatically overrides the CreatedAt field
		err = repo.Update(context.TODO(), found)
		assert.NoError(t, err)
	}

	got, err := repo.GetAll(context.TODO())
	assert.NoError(t, err)
	assert.Len(t, got, len(Imports))

	for i := range Imports {
		assert.Equal(t, got[i].ID, Imports[i].ID)
	}

	deleted, err := repo.DeleteOld(context.TODO(), creationTime.Add(time.Hour))
	assert.Equal(t, int64(3), deleted)
	assert.NoError(t, err)

	got, err = repo.GetAll(context.TODO())
	assert.NoError(t, err)
	assert.Len(t, got, 1)

	assert.Equal(t, got[0].ID, "4BB-00E8-4B0F-97EB-2F3EC3394A87")
}

func buildRepo(silent bool) (*Store, error) {
	err := os.RemoveAll("testdata")
	if err != nil {
		return nil, fmt.Errorf("could not remove all: %w", err)
	}

	err = os.Mkdir("testdata", types.PermissionReadWrite)
	if err != nil {
		return nil, fmt.Errorf("could not make dir: %w", err)
	}

	return NewTruncate(
		"user=postgres password=postgres dbname=test_nameles host=127.0.0.1 port=5432 sslmode=disable",
		silent,
	)
}

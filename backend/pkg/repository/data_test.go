package repository

import (
	"context"
	"fmt"
	"os"
	"testing"

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

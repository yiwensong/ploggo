package postgres

import (
	context "context"
	testing "testing"

	assert "github.com/stretchr/testify/assert"
)

func TestGetGames(t *testing.T) {
	ctx := context.Background()
	storage, cleanup, err := NewAvalonPostgresStorage(ctx, testConnectionString)
	assert.NoError(t, err)
	defer cleanup()

	games, err := storage.GetGames(ctx)
	assert.NotNil(t, games)
	assert.NoError(t, err)
}

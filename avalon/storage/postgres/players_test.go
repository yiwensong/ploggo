package postgres

import (
	context "context"
	testing "testing"

	assert "github.com/stretchr/testify/assert"
	avalon "github.com/yiwensong/ploggo/avalon"
)

func TestCreatePlayer(t *testing.T) {
	ctx := context.Background()
	storage, cleanup, err := NewAvalonPostgresStorage(ctx, testConnectionString)
	assert.NoError(t, err)
	defer cleanup()

	err = storage.CreatePlayer(ctx, &avalon.PlayerImpl{
		Id:       avalon.PlayerId("id"),
		Name:     "yiwen",
		Rating:   300.1,
		NumGames: 120,
	})
	assert.NoError(t, err)
}

func TestGetPlayer(t *testing.T) {
	ctx := context.Background()
	storage, cleanup, err := NewAvalonPostgresStorage(ctx, testConnectionString)
	assert.NoError(t, err)
	defer cleanup()

	id := avalon.NewPlayerId()

	createdPlayer := &avalon.PlayerImpl{
		Id:       id,
		Name:     "yiwen",
		Rating:   300.1,
		NumGames: 120,
	}
	err = storage.CreatePlayer(ctx, createdPlayer)
	assert.NoError(t, err)

	fetchedPlayer, err := storage.GetPlayer(ctx, id)
	assert.Equal(t, fetchedPlayer, createdPlayer)
}

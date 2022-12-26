package postgres

import (
	context "context"
	fmt "fmt"
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

func TestGetPlayersById(t *testing.T) {
	ctx := context.Background()
	storage, cleanup, err := NewAvalonPostgresStorage(ctx, testConnectionString)
	assert.NoError(t, err)
	defer cleanup()

	// Make 4 players
	players := make([]*avalon.PlayerImpl, 4)
	playerIds := make([]avalon.PlayerId, 4)
	for idx := range players {
		players[idx] = avalon.NewPlayer(fmt.Sprintf("player-%d", idx))

		err = storage.CreatePlayer(ctx, players[idx])
		assert.NoError(t, err)

		playerIds[idx] = players[idx].Id
	}

	// Get first 3 player ids
	returnedPlayers, err := storage.GetPlayersById(ctx, playerIds[:3])
	assert.NoError(t, err)
	assert.Equal(t, 3, len(returnedPlayers), "expected returned players to be size 3")
	for _, player := range players[:3] {
		returnedPlayer, ok := returnedPlayers[player.Id]
		assert.True(t, ok, "returned players did not contain player with id %q", player.Id)
		assert.Equal(t, player, returnedPlayer)
	}
}

func TestUpdatePlayers(t *testing.T) {
	ctx := context.Background()
	storage, cleanup, err := NewAvalonPostgresStorage(ctx, testConnectionString)
	assert.NoError(t, err)
	defer cleanup()

	// Make 4 players
	players := make([]*avalon.PlayerImpl, 4)
	playerIds := make([]avalon.PlayerId, 4)
	for idx := range players {
		players[idx] = avalon.NewPlayer(fmt.Sprintf("player-%d", idx))

		err = storage.CreatePlayer(ctx, players[idx])
		assert.NoError(t, err)

		playerIds[idx] = players[idx].Id
	}

	for _, player := range players {
		player.Name = fmt.Sprintf("player-%s", player.Id)
		player.Rating = 1241.
		player.NumGames++
	}

	err = storage.UpdatePlayers(ctx, players)
	assert.NoError(t, err)

	playersById, err := storage.GetPlayersById(ctx, playerIds)
	for _, player := range playersById {
		assert.Equal(t, 1241., player.Rating)
		assert.Equal(t, 1, player.NumGames)
	}
}

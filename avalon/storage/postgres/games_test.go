package postgres

import (
	context "context"
	"fmt"
	testing "testing"
	"time"

	assert "github.com/stretchr/testify/assert"
	"github.com/yiwensong/ploggo/avalon"
)

var fakeTimestamp = time.Date(2000, 1, 1, 0, 0, 0, 0, time.FixedZone("UTC-8", -8*60*60))

func TestGetGames(t *testing.T) {
	ctx := context.Background()
	storage, cleanup, err := NewAvalonPostgresStorage(ctx, testConnectionString)
	assert.NoError(t, err)
	defer cleanup()

	games, err := storage.GetGames(ctx)
	assert.NotNil(t, games)
	assert.NoError(t, err)
}

func TestSaveGame(t *testing.T) {
	ctx := context.Background()
	storage, cleanup, err := NewAvalonPostgresStorage(ctx, testConnectionString)
	assert.NoError(t, err)
	defer cleanup()

	// Make 10 players, give first 6 good role and last 4 bad role
	playersById := map[avalon.PlayerId]*avalon.PlayerImpl{}
	rolesByPlayerId := map[avalon.PlayerId]avalon.Role{}
	for idx := range make([]interface{}, 10) {
		player := avalon.NewPlayer(fmt.Sprintf("player-%d", idx))

		role := avalon.LoyalServant
		if idx > 5 {
			role = avalon.MinionOfMordred
		}

		rolesByPlayerId[player.Id] = role
		playersById[player.Id] = player
	}

	game := avalon.NewGame(playersById, rolesByPlayerId)
	game.SetWinner(avalon.Red)
	err = storage.SaveGame(ctx, game)
	assert.NoError(t, err)

	// Ensure game is saved by loading
	returnedGame, err := storage.GetGame(ctx, game.Id)
	assert.NoError(t, err)

	// Save times from both and compare later
	gameTime := game.CreatedAt
	returnedTime := returnedGame.CreatedAt

	// Fake a date so we can compare structs
	game.CreatedAt = fakeTimestamp
	returnedGame.CreatedAt = fakeTimestamp
	assert.Equal(t, game, returnedGame)

	// Check that the two's UTC time matches
	assert.Equal(t, gameTime.UTC(), returnedTime.UTC())
}

package postgres

import (
	context "context"
	"fmt"
	testing "testing"

	assert "github.com/stretchr/testify/assert"
	"github.com/yiwensong/ploggo/avalon"
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
	err = storage.SaveGame(ctx, game)
	assert.NoError(t, err)

	// Ensure game is saved by loading
}

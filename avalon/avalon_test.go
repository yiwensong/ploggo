package avalon

import (
	fmt "fmt"
	math "math"
	runtime "runtime"
	testing "testing"

	assert "github.com/stretchr/testify/assert"
)

var _ *runtime.Func = nil

type gameTestScenario struct {
	game    *GameImpl
	players *[10]*PlayerImpl
	roles   [10]Role
}

func generateGameTestScenario(t *testing.T) *gameTestScenario {
	var roles [10]Role = [10]Role{
		LoyalServant,
		LoyalServant,
		LoyalServant,
		LoyalServant,
		Percival,
		Merlin,
		Assassin,
		Mordred,
		Morgana,
		MinionOfMordred,
	}

	players := new([10]*PlayerImpl)
	rolesByPlayerId := map[PlayerId]Role{}
	playersById := map[PlayerId]*PlayerImpl{}
	for i := range players {
		players[i] = NewPlayer(fmt.Sprintf("player %d", i))

		rolesByPlayerId[players[i].Id] = roles[i]
		playersById[players[i].Id] = players[i]
	}

	game := &GameImpl{
		Id:             "gameId",
		Winner:         Blue,
		RoleByPlayerId: rolesByPlayerId,
		PlayersById:    playersById,
	}

	return &gameTestScenario{
		game:    game,
		roles:   roles,
		players: players,
	}
}

func Test_GameImpl_GetWinPercentage(t *testing.T) {
	scenario := generateGameTestScenario(t)
	game := scenario.game

	winnersWinPercent, err := game.GetWinPercentage(game.Winner)
	assert.NoError(t, err)

	losersWinPercent, err := game.GetWinPercentage(game.Winner.OtherTeam())
	assert.NoError(t, err)

	assert.True(t, math.Abs(winnersWinPercent+losersWinPercent-1) < 0.0000001)
}

func Test_GameImpl_GetNewRatingForPlayer(t *testing.T) {
	scenario := generateGameTestScenario(t)
	game := scenario.game
	players := scenario.players

	for i, player := range players {
		newRating, err := game.GetNewRatingForPlayer(player.Id)
		assert.NoError(t, err)

		// first 6 players are good
		if i < 6 {
			assert.Truef(t, newRating > player.Rating, "player number %d (id=%q) should have positive rating update new_rating=%f", i, player.Id, newRating)
		} else {
			assert.Truef(t, player.Rating > newRating, "player number %d (id=%q) should have negative rating update new_rating=%f", i, player.Id, newRating)
		}
	}
}

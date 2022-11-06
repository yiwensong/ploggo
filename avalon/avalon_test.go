package avalon

import (
	fmt "fmt"
	math "math"
	runtime "runtime"
	testing "testing"

	assert "github.com/stretchr/testify/assert"
)

var _ *runtime.Func = nil

type gameTestScenarioOpts struct {
	ratings *[10]float64
	winner  *Team
}

type gameTestScenario struct {
	game    *GameImpl
	players *[10]*PlayerImpl
	roles   [10]Role
}

func generateGameTestScenario(t *testing.T, opts *gameTestScenarioOpts) *gameTestScenario {
	if opts == nil {
		opts = &gameTestScenarioOpts{}
	}

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

		if opts.ratings != nil {
			players[i].Rating = opts.ratings[i]
		}

	}

	game := NewGame(
		playersById,
		rolesByPlayerId,
	)
	winner := Blue
	if opts.winner != nil {
		winner = *opts.winner
	}
	game.SetWinner(winner)

	return &gameTestScenario{
		game:    game,
		roles:   roles,
		players: players,
	}
}

func Test_GameImpl_GetWinPercentage(t *testing.T) {
	tests := []struct {
		testName string
		ratings  [10]float64
	}{
		{
			testName: "all starting",
			ratings:  [10]float64{1500, 1500, 1500, 1500, 1500, 1500, 1500, 1500, 1500, 1500},
		},
		{
			testName: "one really high, winning",
			ratings:  [10]float64{3000, 1500, 1500, 1500, 1500, 1500, 1500, 1500, 1500, 1500},
		},
		{
			testName: "one really high, losing",
			ratings:  [10]float64{1500, 1500, 1500, 1500, 1500, 1500, 1500, 1500, 1500, 3000},
		},
		{
			testName: "all low",
			ratings:  [10]float64{300, 300, 300, 300, 300, 300, 300, 300, 300, 300},
		},
		{
			testName: "various",
			ratings:  [10]float64{300, 200, 1500, 3000, 1500, 100, 2700, 1800, 1500, 1500},
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			scenario := generateGameTestScenario(t, nil)
			game := scenario.game

			winnersWinPercent, err := game.GetWinPercentage(game.Winner)
			assert.NoError(t, err)
			winnersWinPercent2, err := game.GetWinPercentage(game.Winner)
			assert.NoError(t, err)
			assert.Equal(t, winnersWinPercent, winnersWinPercent2)

			losersWinPercent, err := game.GetWinPercentage(game.Winner.OtherTeam())
			assert.NoError(t, err)
			losersWinPercent2, err := game.GetWinPercentage(game.Winner.OtherTeam())
			assert.NoError(t, err)
			assert.Equal(t, losersWinPercent, losersWinPercent2)

			assert.True(t, math.Abs(winnersWinPercent+losersWinPercent-1) < 0.0000001)
		})
	}
}

func Test_GameImpl_GetNewRatingForPlayer(t *testing.T) {
	scenario := generateGameTestScenario(t, nil)
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

func Test_GameImpl_UpdatePlayersAfterGame(t *testing.T) {
	tests := []struct {
		testName string
		ratings  [10]float64
	}{
		{
			testName: "all starting",
			ratings:  [10]float64{1500, 1500, 1500, 1500, 1500, 1500, 1500, 1500, 1500, 1500},
		},
		{
			testName: "one really high, winning",
			ratings:  [10]float64{3000, 1500, 1500, 1500, 1500, 1500, 1500, 1500, 1500, 1500},
		},
		{
			testName: "one really high, losing",
			ratings:  [10]float64{1500, 1500, 1500, 1500, 1500, 1500, 1500, 1500, 1500, 3000},
		},
		{
			testName: "all low",
			ratings:  [10]float64{300, 300, 300, 300, 300, 300, 300, 300, 300, 300},
		},
		{
			testName: "various",
			ratings:  [10]float64{300, 200, 1500, 3000, 1500, 100, 2700, 1800, 1500, 1500},
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			scenario := generateGameTestScenario(t, nil)
			game := scenario.game
			game.SetWinner(Blue)

			updatedPlayers, err := game.UpdatePlayersAfterGame()
			assert.NoError(t, err)
			assert.Equal(t, len(updatedPlayers), 10)
		})
	}
}

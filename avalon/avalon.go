package avalon

import (
	fmt "fmt"
	math "math"

	uuid "github.com/google/uuid"
	errors "github.com/pkg/errors"
)

// id types
type (
	PlayerId string
	GameId   string
)

type PlayerImpl struct {
	Id     PlayerId
	Name   string
	Rating float64

	NumGames int
}

func NewPlayer(name string) *PlayerImpl {
	var playerId PlayerId = PlayerId(fmt.Sprintf("usr_%s", uuid.NewString()))
	return &PlayerImpl{
		Id:       playerId,
		Name:     name,
		Rating:   1500,
		NumGames: 0,
	}
}

func (p *PlayerImpl) GetUpdateConstant() float64 {
	return float64(300) / math.Log(float64(p.NumGames+2))
}

type Team string

const (
	NoTeam Team = "no_team"
	Blue   Team = "blue_team"
	Red         = "red_team"
)

func (t Team) OtherTeam() Team {
	if t == NoTeam {
		return NoTeam
	}

	if t == Blue {
		return Red
	}

	return Blue
}

type Role struct {
	Name string
	Team Team
}

var (
	LoyalServant = Role{
		Name: "loyal_servant",
		Team: Blue,
	}
	Merlin = Role{
		Name: "merlin",
		Team: Blue,
	}
	Percival = Role{
		Name: "percival",
		Team: Blue,
	}

	Morgana = Role{
		Name: "morgana",
		Team: Red,
	}
	Assassin = Role{
		Name: "assassin",
		Team: Red,
	}
	MinionOfMordred = Role{
		Name: "minion_of_mordred",
		Team: Red,
	}
	Mordred = Role{
		Name: "mordred",
		Team: Red,
	}
	Oberon = Role{
		Name: "oberon",
		Team: Red,
	}
)

type GameImpl struct {
	Id             GameId
	Winner         Team
	RoleByPlayerId map[PlayerId]Role

	// Keep a static snapshot of players so the game calculation is repeatable
	PlayersById map[PlayerId]*PlayerImpl
}

// Returns a list of players on a team
func (g *GameImpl) GetTeam(
	team Team,
) (players []*PlayerImpl, err error) {
	for playerId, role := range g.RoleByPlayerId {
		player, ok := g.PlayersById[playerId]
		if !ok {
			return nil, errors.Errorf("Game contained a player id that did not exist game_id=%q player_id=%q", g.Id, playerId)
		}

		if role.Team == team {
			players = append(players, player)
		}
	}

	return players, nil
}

// Returns the arthmetic average team rating for a team in the game
func (g *GameImpl) GetTeamRating(
	team Team,
) (rating float64, err error) {
	players, err := g.GetTeam(team)
	if err != nil {
		return 0, errors.Wrapf(err, "GetTeam(%q, %q)", team, g.Id)
	}

	var total float64 = 0.0
	var numPlayers float64 = 0.0
	for _, player := range players {
		total += player.Rating
		numPlayers += 1.0
	}

	return total / numPlayers, nil
}

// Returns the expectation of the team winning, which is:
//
//   1/(1 + 10 ** ((r_other - r_team)/400))
func (g *GameImpl) GetWinPercentage(
	team Team,
) (winRate float64, err error) {
	playerTeamRating, err := g.GetTeamRating(team)
	if err != nil {
		return 0, errors.Wrapf(err, "GetTeamRating(%q)", team)
	}

	otherTeam := team.OtherTeam()
	otherTeamRating, err := g.GetTeamRating(otherTeam)
	if err != nil {
		return 0, errors.Wrapf(err, "GetTeamRating(%q)", otherTeam)
	}

	return 1 / (1 + (math.Pow(10, (otherTeamRating-playerTeamRating)/400))), nil
}

// Given a game, return what the update should be for that player
func (g *GameImpl) GetNewRatingForPlayer(
	playerId PlayerId,
) (rating float64, err error) {
	if g.Winner == NoTeam {
		return 0, errors.Errorf("Game had no winner, update is impossible game_id=%q", g.Id)
	}

	player, ok := g.PlayersById[playerId]
	if !ok {
		return 0, errors.Errorf("player id not found player_id=%q game_id=%q", playerId, g.Id)
	}

	role, ok := g.RoleByPlayerId[playerId]
	if !ok {
		return 0, errors.Errorf("Player was not found in game player_id=%q game_id=%q", playerId, g.Id)
	}

	// Elo update is:
	// rating_1 = rating_0 + K * (Score - Expected)
	// where:
	//     * K is the value from PlayerImpl.GetNewRatingForPlayer
	//     * score is 1 if the player won, otherwise zero
	//     * expected is the calculated win rate for the player's team
	updateConstant := player.GetUpdateConstant()
	var gameScore float64
	if g.Winner == role.Team {
		gameScore = 1.0
	} else {
		gameScore = 0.0
	}
	expectedWin, err := g.GetWinPercentage(role.Team)
	if err != nil {
		return 0, errors.Wrapf(err, "GetWinPercentage(%q, %q)", g.Id, role.Team)
	}

	newRating := player.Rating + updateConstant*(gameScore-expectedWin)

	return newRating, nil
}

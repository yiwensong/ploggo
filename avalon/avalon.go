package avalon

import (
	fmt "fmt"
	math "math"
	"time"

	"github.com/google/uuid"
	errors "github.com/pkg/errors"
	ksuid "github.com/segmentio/ksuid"
)

// id types
type (
	PlayerId string
	GameId   string
)

type PlayerImpl struct {
	Id       PlayerId
	Name     string
	Rating   float64
	NumGames int
}

func NewPlayerId() PlayerId {
	return PlayerId(uuid.New().String())
}

func NewPlayer(name string) *PlayerImpl {
	playerId := NewPlayerId()
	return &PlayerImpl{
		Id:       playerId,
		Name:     name,
		Rating:   1500,
		NumGames: 0,
	}
}

func (p *PlayerImpl) GetUpdateConstant() float64 {
	return float64(150) / math.Log(float64(p.NumGames+2))
}

type Team string

const (
	NoTeam Team = "no_team"
	Blue   Team = "blue_team"
	Red    Team = "red_team"
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

type Role string

const (
	LoyalServant    Role = "loyal_servant"
	Merlin               = "merlin"
	Percival             = "percival"
	Morgana              = "morgana"
	Assassin             = "assassin"
	MinionOfMordred      = "minion_of_mordred"
	Mordred              = "mordred"
	Oberon               = "oberon"
)

var TeamByRole map[Role]Team = map[Role]Team{
	LoyalServant:    Blue,
	Merlin:          Blue,
	Percival:        Blue,
	Morgana:         Red,
	Assassin:        Red,
	MinionOfMordred: Red,
	Mordred:         Red,
	Oberon:          Red,
}

type GameImpl struct {
	Id             GameId
	Winner         Team
	RoleByPlayerId map[PlayerId]Role
	Timestamp      time.Time

	// Keep a static snapshot of players so the game calculation is repeatable
	PlayersById map[PlayerId]*PlayerImpl

	teamWinPercentages map[Team]float64
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

		if TeamByRole[role] == team {
			players = append(players, player)
		}
	}

	return players, nil
}

// Creates a new game
func NewGame(
	playersById map[PlayerId]*PlayerImpl,
	roleByPlayerId map[PlayerId]Role,
) *GameImpl {
	var gameId GameId = GameId(fmt.Sprintf("gam_%s", ksuid.New().String()))
	return &GameImpl{
		Id:             gameId,
		RoleByPlayerId: roleByPlayerId,
		PlayersById:    playersById,
		Timestamp:      time.Now(),
	}
}

func (g *GameImpl) SetWinner(team Team) {
	g.Winner = team
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
//	1/(1 + 10 ** ((r_other - r_team)/400))
//
// cached via teamWinPercentages in game
func (g *GameImpl) GetWinPercentage(
	team Team,
) (winRate float64, err error) {
	if g.teamWinPercentages == nil {
		g.teamWinPercentages = map[Team]float64{}
	}

	_, ok := g.teamWinPercentages[team]
	if !ok {

		playerTeamRating, err := g.GetTeamRating(team)
		if err != nil {
			return 0, errors.Wrapf(err, "GetTeamRating(%q)", team)
		}

		otherTeam := team.OtherTeam()
		otherTeamRating, err := g.GetTeamRating(otherTeam)
		if err != nil {
			return 0, errors.Wrapf(err, "GetTeamRating(%q)", otherTeam)
		}

		g.teamWinPercentages[team] = 1 / (1 + (math.Pow(10, (otherTeamRating-playerTeamRating)/400)))
	}
	return g.teamWinPercentages[team], nil
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
	if g.Winner == TeamByRole[role] {
		gameScore = 1.0
	} else {
		gameScore = 0.0
	}
	expectedWin, err := g.GetWinPercentage(TeamByRole[role])
	if err != nil {
		return 0, errors.Wrapf(err, "GetWinPercentage(%q, %q)", g.Id, TeamByRole[role])
	}

	newRating := player.Rating + updateConstant*(gameScore-expectedWin)

	return newRating, nil
}

// Returns a list of updated players for when the game is finished
//
// The list of players is deep copied to preserve the game's state
func (g *GameImpl) UpdatePlayersAfterGame() (updatedPlayers []*PlayerImpl, err error) {
	for playerId, player := range g.PlayersById {
		updatedRating, err := g.GetNewRatingForPlayer(playerId)
		if err != nil {
			return nil, errors.Wrapf(err, "GetNewRatingForPlayer(%q)", playerId)
		}

		updatedPlayers = append(updatedPlayers, &PlayerImpl{
			Id:       playerId,
			Name:     player.Name,
			Rating:   updatedRating,
			NumGames: player.NumGames + 1,
		})
	}
	return
}

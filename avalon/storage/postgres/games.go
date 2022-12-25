package postgres

import (
	context "context"

	dbr "github.com/gocraft/dbr"
	dialect "github.com/gocraft/dbr/dialect"
	pgx "github.com/jackc/pgx/v5"
	errors "github.com/pkg/errors"
	avalon "github.com/yiwensong/ploggo/avalon"
)

// A row for `games` table in db
type GameEntry struct {
	Id      string
	BlueWon bool
}

// A row for `players_by_game` table in db
type PlayersByGameEntry struct {
	GameId         string
	PlayerId       string
	PlayerName     string
	PlayerRating   float64
	PlayerNumGames int64
	PlayerRole     string
}

// A row from the results when joining `players_by_name`
// and `games` tables
type JoinedPlayersByGame struct {
	GameEntry
	PlayersByGameEntry
}

// Maximum number of games to query
const MAX_GAMES = 100

func (s *AvalonPostgresStorage) SaveGame(ctx context.Context, game *avalon.GameImpl) error {
	return errors.New("Not implemented")
}

func (s *AvalonPostgresStorage) GetGames(ctx context.Context) (games []*avalon.GameImpl, err error) {
	stmt := dbr.
		Select(
			"games.id",
			"games.blue_won",
			"games.blue_win_expectation",
			"p.player_id",
			"p.player_name",
			"p.player_rating",
			"p.player_num_games",
			"p.player_role",
		).
		From(
			dbr.Select("players_by_game").As("p"),
		).
		Join(
			dbr.Select("games").Limit(MAX_GAMES),
			"p.game_id = games.id",
		).
		OrderBy("games.id")

	buffer := dbr.NewBuffer()
	err = stmt.Build(dialect.PostgreSQL, buffer)
	if err != nil {
		return nil, errors.Wrapf(err, "dbr.Build")
	}

	resultRows := []*JoinedPlayersByGame{}
	perform := func(ctx context.Context, tx pgx.Tx) error {
		rows, err := tx.Query(ctx, buffer.String())
		if err != nil {
			return errors.Wrapf(err, "tx.Query(%q)", buffer.String())
		}
		defer rows.Close()

		for rows.Next() {
			var resultRow JoinedPlayersByGame

			err = rows.Scan(
				&resultRow.Id,
				&resultRow.BlueWon,
				&resultRow.PlayerId,
				&resultRow.PlayerName,
				&resultRow.PlayerRating,
				&resultRow.PlayerNumGames,
				&resultRow.PlayerRole,
			)
			if err != nil {
				return errors.Wrapf(err, "rows.Scan")
			}

			// Repeated value, added for consistency
			resultRow.GameId = resultRow.Id

			resultRows = append(resultRows, &resultRow)
		}

		return nil
	}

	err = s.WithTx(ctx, perform)
	if err != nil {
		return nil, errors.Wrapf(err, "WithTx(%q)", buffer.String)
	}

	games = []*avalon.GameImpl{}
	lastGameId := ""
	var nextGame *avalon.GameImpl

	for _, row := range resultRows {
		playerId := avalon.PlayerId(row.PlayerId)
		player := &avalon.PlayerImpl{
			Id:       playerId,
			Name:     row.PlayerName,
			Rating:   row.PlayerRating,
			NumGames: int(row.PlayerNumGames),
		}

		if nextGame != nil && row.Id == lastGameId {
			// if the game is already added, only update player info
			nextGame.RoleByPlayerId[playerId] = avalon.Role(row.PlayerRole)
			nextGame.PlayersById[playerId] = player
		} else {
			// this is a new game, add a new one
			winner := avalon.Red
			if row.BlueWon {
				winner = avalon.Blue
			}

			nextGame = &avalon.GameImpl{
				Id:     avalon.GameId(row.Id),
				Winner: winner,
				RoleByPlayerId: map[avalon.PlayerId]avalon.Role{
					playerId: avalon.Role(row.PlayerRole),
				},
				PlayersById: map[avalon.PlayerId]*avalon.PlayerImpl{
					playerId: player,
				},
			}
			games = append(games, nextGame)
			lastGameId = row.Id
		}
	}

	return games, nil
}

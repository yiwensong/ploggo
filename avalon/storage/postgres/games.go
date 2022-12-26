package postgres

import (
	context "context"

	glog "github.com/golang/glog"
	pgx "github.com/jackc/pgx/v5"
	pgtype "github.com/jackc/pgx/v5/pgtype"
	errors "github.com/pkg/errors"
	avalon "github.com/yiwensong/ploggo/avalon"
)

// A row for `games` table in db
type GameEntry struct {
	Id        pgtype.Text
	BlueWon   pgtype.Bool
	Timestamp pgtype.Timestamp
}

// A row for `players_by_game` table in db
type PlayersByGameEntry struct {
	GameId         pgtype.Text
	PlayerId       pgtype.Text
	PlayerName     pgtype.Text
	PlayerRating   pgtype.Float8
	PlayerNumGames pgtype.Int8
	PlayerRole     pgtype.Text
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
	query := `
		SELECT
			games.id,
			games.blue_won,
			games.created_at,
			p.player_id,
			p.player_name,
			p.player_rating,
			p.player_num_games,
			p.player_role
		FROM players_by_game AS p
		JOIN (
			SELECT
				id,
				blue_won,
				created_at
			FROM games
			LIMIT $1
		) AS games
		ON p.game_id = games.id
		ORDER BY games.id`

	resultRows := []*JoinedPlayersByGame{}
	perform := func(ctx context.Context, tx pgx.Tx) error {
		glog.Infof("GetGames")

		rows, err := tx.Query(ctx, query, MAX_GAMES)
		if err != nil {
			return errors.Wrapf(err, "tx.Query(%q)", query)
		}
		defer rows.Close()

		for rows.Next() {
			var resultRow JoinedPlayersByGame

			err = rows.Scan(
				&resultRow.Id,
				&resultRow.BlueWon,
				&resultRow.Timestamp,
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
		return nil, errors.Wrapf(err, "WithTx(%q)", query)
	}

	games = []*avalon.GameImpl{}
	lastGameId := ""
	var nextGame *avalon.GameImpl

	for _, row := range resultRows {
		playerId := avalon.PlayerId(row.PlayerId.String)
		player := &avalon.PlayerImpl{
			Id:       playerId,
			Name:     row.PlayerName.String,
			Rating:   row.PlayerRating.Float64,
			NumGames: int(row.PlayerNumGames.Int64),
		}

		if nextGame != nil && row.Id.String == lastGameId {
			// if the game is already added, only update player info
			nextGame.RoleByPlayerId[playerId] = avalon.Role(row.PlayerRole.String)
			nextGame.PlayersById[playerId] = player
		} else {
			// this is a new game, add a new one
			winner := avalon.Red
			if row.BlueWon.Bool {
				winner = avalon.Blue
			}

			nextGame = &avalon.GameImpl{
				Id:        avalon.GameId(row.Id.String),
				Winner:    winner,
				Timestamp: row.Timestamp.Time,
				RoleByPlayerId: map[avalon.PlayerId]avalon.Role{
					playerId: avalon.Role(row.PlayerRole.String),
				},
				PlayersById: map[avalon.PlayerId]*avalon.PlayerImpl{
					playerId: player,
				},
			}
			games = append(games, nextGame)
			lastGameId = row.Id.String
		}
	}

	return games, nil
}

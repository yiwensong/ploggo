package postgres

import (
	context "context"
	"fmt"

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
	CreatedAt pgtype.Timestamp
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

type PlayersByGameEntryForWriting struct {
	GameId         avalon.GameId
	PlayerId       avalon.PlayerId
	PlayerName     string
	PlayerRating   float64
	PlayerNumGames int
	PlayerRole     avalon.Role
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
	insertGameQuery := `
		INSERT INTO games (
			id,
			blue_won,
			created_at
		)
		VALUES (
			$1,
			$2,
			$3
		)`

	insertPlayersByGameQuery := `
		INSERT INTO players_by_game (
			game_id,
			player_id,
			player_name,
			player_rating,
			player_num_games,
			player_role
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6
		)`

	playersByGameRows := []*PlayersByGameEntryForWriting{}
	for playerId, role := range game.RoleByPlayerId {
		player, ok := game.PlayersById[playerId]
		if !ok {
			return errors.Errorf("Player %q had a role but no player data", playerId)
		}

		playersByGameRows = append(playersByGameRows, &PlayersByGameEntryForWriting{
			GameId:         game.Id,
			PlayerId:       playerId,
			PlayerName:     player.Name,
			PlayerRating:   player.Rating,
			PlayerNumGames: player.NumGames,
			PlayerRole:     role,
		})
	}

	blueWon := game.Winner == avalon.Blue

	return s.WithTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		// Insert the game
		_, err := tx.Exec(ctx, insertGameQuery, game.Id, blueWon, game.CreatedAt)
		if err != nil {
			return errors.Wrapf(err, "tx.Exec(%q, %+v)", insertGameQuery, game)
		}

		for _, row := range playersByGameRows {
			_, err := tx.Exec(
				ctx,
				insertPlayersByGameQuery,
				row.GameId,
				row.PlayerId,
				row.PlayerName,
				row.PlayerRating,
				row.PlayerNumGames,
				row.PlayerRole,
			)
			if err != nil {
				return errors.Wrapf(err, "tx.Exec(%q, %+v)", insertPlayersByGameQuery, row)
			}
		}

		return nil
	})
}

type GameQueryOpts struct {
	Id avalon.GameId
}

func (s *AvalonPostgresStorage) getGamesQueryLoader(
	ctx context.Context,
	opts *GameQueryOpts,
) (games []*avalon.GameImpl, err error) {
	queryTemplate := `
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
			%s
			LIMIT $1
		) AS games
		ON p.game_id = games.id
		ORDER BY games.id`

	args := []interface{}{
		MAX_GAMES,
	}

	whereStmt := ""
	argN := 2
	if opts != nil && opts.Id != "" {
		whereStmt = fmt.Sprintf("WHERE games.id = $%d", argN)
		args = append(args, opts.Id)
		argN++
	}

	query := fmt.Sprintf(
		queryTemplate,
		whereStmt,
	)
	resultRows := []*JoinedPlayersByGame{}
	perform := func(ctx context.Context, tx pgx.Tx) error {
		glog.Infof("GetGames")

		rows, err := tx.Query(ctx, query, args...)
		if err != nil {
			return errors.Wrapf(err, "tx.Query(%q)", query)
		}
		defer rows.Close()

		for rows.Next() {
			var resultRow JoinedPlayersByGame

			err = rows.Scan(
				&resultRow.Id,
				&resultRow.BlueWon,
				&resultRow.CreatedAt,
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
				CreatedAt: row.CreatedAt.Time,
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

func (s *AvalonPostgresStorage) GetGame(
	ctx context.Context,
	id avalon.GameId,
) (game *avalon.GameImpl, err error) {
	games, err := s.getGamesQueryLoader(ctx, &GameQueryOpts{Id: id})
	if err != nil {
		return nil, errors.Wrapf(err, "getGamesQueryLoader(%q)", id)
	}

	if len(games) < 1 {
		return nil, errors.Errorf("Game with id=%q does not exist", id)
	}

	return games[0], nil
}

func (s *AvalonPostgresStorage) GetGames(
	ctx context.Context,
) (games []*avalon.GameImpl, err error) {
	return s.getGamesQueryLoader(ctx, nil)
}

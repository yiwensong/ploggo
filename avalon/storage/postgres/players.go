package postgres

import (
	context "context"

	glog "github.com/golang/glog"
	pgx "github.com/jackc/pgx/v5"
	pgtype "github.com/jackc/pgx/v5/pgtype"
	errors "github.com/pkg/errors"
	avalon "github.com/yiwensong/ploggo/avalon"
)

type PlayerEntry struct {
	Id       pgtype.Text
	Name     pgtype.Text
	Rating   pgtype.Float8
	NumGames pgtype.Int8
}

func (s *AvalonPostgresStorage) GetPlayersById(
	ctx context.Context,
	playerIds []avalon.PlayerId,
) (map[avalon.PlayerId]*avalon.PlayerImpl, error) {
	query := `
		SELECT
			id,
			name,
			rating,
			num_games
		FROM players
		WHERE
			id = ANY($1)`

	players := map[avalon.PlayerId]*avalon.PlayerImpl{}
	err := s.WithTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		rows, err := tx.Query(ctx, query, playerIds)
		if err != nil {
			return errors.Wrapf(err, "Query(%q, %+v)", query, playerIds)
		}
		defer rows.Close()

		for rows.Next() {
			var entry PlayerEntry
			err = rows.Scan(
				&entry.Id,
				&entry.Name,
				&entry.Rating,
				&entry.NumGames,
			)
			if err != nil {
				return errors.Wrapf(err, "rows.Scan")
			}

			players[avalon.PlayerId(entry.Id.String)] = &avalon.PlayerImpl{
				Id:       avalon.PlayerId(entry.Id.String),
				Name:     entry.Name.String,
				Rating:   entry.Rating.Float64,
				NumGames: int(entry.NumGames.Int64),
			}
		}
		return nil
	})

	if err != nil {
		return nil, errors.Wrapf(err, "WithTx")
	}

	return players, nil
}

func (s *AvalonPostgresStorage) GetPlayer(
	ctx context.Context,
	playerId avalon.PlayerId,
) (player *avalon.PlayerImpl, err error) {
	query := `
		SELECT
			id,
			name,
			rating,
			num_games
		FROM players
		WHERE
			id = $1`

	err = s.WithTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		glog.Infof("GetPlayer")
		rows, err := tx.Query(ctx, query, string(playerId))
		if err != nil {
			return errors.Wrapf(err, "Query(%s, %s)", query, playerId)
		}
		defer rows.Close()

		if !rows.Next() {
			return errors.Errorf("Player with id %q does not exist", playerId)
		}

		var entry PlayerEntry
		err = rows.Scan(&entry.Id, &entry.Name, &entry.Rating, &entry.NumGames)
		if err != nil {
			return errors.Wrapf(err, "rows.Scan")
		}

		player = &avalon.PlayerImpl{
			Id:       avalon.PlayerId(entry.Id.String),
			Name:     entry.Name.String,
			Rating:   entry.Rating.Float64,
			NumGames: int(entry.NumGames.Int64),
		}

		return nil
	})

	return player, err
}

func (s *AvalonPostgresStorage) CreatePlayer(
	ctx context.Context,
	player *avalon.PlayerImpl,
) error {
	return s.WithTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		glog.Infof("CreatePlayer")

		query := `
			INSERT INTO players (
				id,
				name,
				rating,
				num_games
			)
			VALUES (
				$1,
				$2,
				$3,
				$4
			)`
		_, err := tx.Exec(
			ctx,
			query,
			string(player.Id),
			player.Name,
			player.Rating,
			player.NumGames,
		)
		if err != nil {
			return errors.Wrapf(err, "Exec(%q, %+v)", query, player)
		}

		return nil
	})
}

func (s *AvalonPostgresStorage) UpdatePlayers(
	ctx context.Context,
	players []*avalon.PlayerImpl,
) error {
	query := `
		UPDATE players
		SET
			name = $2,
			rating = $3,
			num_games = $4
		WHERE
			id = $1`

	// Do this for each player
	return s.WithTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		for _, player := range players {
			ct, err := tx.Exec(
				ctx,
				query,
				player.Id,
				player.Name,
				player.Rating,
				player.NumGames,
			)
			if err != nil {
				return errors.Wrapf(err, "tx.Exec(%q, %+v)", query, player)
			}
			if rowsAffected := ct.RowsAffected(); rowsAffected == 0 {
				return errors.Errorf("Player with id %q was not found", player.Id)
			}
		}

		return nil
	})
}

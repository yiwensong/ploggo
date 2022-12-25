package postgres

import (
	context "context"

	glog "github.com/golang/glog"
	pgx "github.com/jackc/pgx/v5"
	pgxpool "github.com/jackc/pgx/v5/pgxpool"
	errors "github.com/pkg/errors"
	storage "github.com/yiwensong/ploggo/avalon/storage"
)

type AvalonPostgresStorage struct {
	Pool *pgxpool.Pool
}

func NewAvalonPostgresStorage(
	ctx context.Context,
	connString string,
) (
	store *AvalonPostgresStorage,
	cleanup func(),
	err error,
) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, func() {}, errors.Wrapf(err, "pgxpool.New")
	}

	store = &AvalonPostgresStorage{
		Pool: pool,
	}

	return store, pool.Close, nil
}

func (s *AvalonPostgresStorage) WithTx(
	ctx context.Context,
	perform func(context.Context, pgx.Tx) error,
) error {
	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return errors.Wrapf(err, "pgxpool.Begin")
	}
	defer func() {
		// If there is no error, try to commit the transaction
		if err == nil {
			err = tx.Commit(ctx)
		}

		// If there is an error (including the commit attempt), roll back
		if err != nil {
			rollBackErr := tx.Rollback(ctx)
			glog.Errorf("WithTx encountered an error but db rollback failed: %s", rollBackErr.Error())
		}
	}()

	err = perform(ctx, tx)
	if err != nil {
		err = errors.Wrapf(err, "AvalonPostgresStorage.WithTx.perform")
		return err
	}

	return err
}

var _ storage.AvalonStorage = (*AvalonPostgresStorage)(nil)

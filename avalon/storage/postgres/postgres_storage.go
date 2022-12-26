package postgres

import (
	context "context"
	"fmt"
	"os"
	"path/filepath"

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
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, func() {}, errors.Wrapf(err, "pgx.ParseConfig")
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
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
) (err error) {
	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return errors.Wrapf(err, "pgxpool.Begin")
	}
	defer func() {
		// If there is an error (including the commit attempt), roll back
		if tx != nil && err != nil {
			glog.Errorf("WithTx encountered an error and is attempting to rollback: %s", err.Error())
			rollBackErr := tx.Rollback(ctx)
			if rollBackErr != nil {
				glog.Errorf("WithTx encountered an error but db rollback failed: %s", rollBackErr.Error())
			}
		}
	}()

	err = perform(ctx, tx)
	if err != nil {
		err = errors.Wrapf(err, "AvalonPostgresStorage.WithTx.perform")
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return errors.Wrapf(err, "tx.Commit")
	}

	return err
}

func (s *AvalonPostgresStorage) RunMigrations(
	ctx context.Context,
	migrationFilesPath string,
) error {
	// Load all migration files
	globPattern := fmt.Sprintf("%s/*.sql", migrationFilesPath)
	matches, err := filepath.Glob(globPattern)
	if err != nil {
		return errors.Wrapf(err, "filepath.Glob(%s)", globPattern)
	}

	migrations := []string{}
	for _, file := range matches {
		bytes, err := os.ReadFile(file)
		if err != nil {
			return errors.Wrapf(err, "os.ReadFile(%q)", file)
		}

		migrations = append(migrations, string(bytes))
	}

	return s.WithTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		glog.Infof("RunMigrations")

		for _, migration := range migrations {
			_, err := tx.Exec(ctx, migration)
			if err != nil {
				return errors.Wrapf(err, "tx.Exec(%q)", migration)
			}
		}

		return nil
	})
}

var _ storage.AvalonStorage = (*AvalonPostgresStorage)(nil)

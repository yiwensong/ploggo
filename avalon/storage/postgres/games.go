package postgres

import (
	context "context"

	errors "github.com/pkg/errors"
	"github.com/yiwensong/ploggo/avalon"
)

func (s *AvalonPostgresStorage) SaveGame(ctx context.Context, game *avalon.GameImpl) error {
	return errors.New("Not implemented")
}

func (s *AvalonPostgresStorage) GetGames(ctx context.Context) (games []*avalon.GameImpl, err error) {
	return nil, errors.New("Not implemented")
}

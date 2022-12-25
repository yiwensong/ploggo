package postgres

import (
	context "context"

	errors "github.com/pkg/errors"
	avalon "github.com/yiwensong/ploggo/avalon"
)

func (s *AvalonPostgresStorage) GetPlayersById(ctx context.Context, playerIds []avalon.PlayerId) (map[avalon.PlayerId]*avalon.PlayerImpl, error) {
	return nil, errors.New("Not implemented")
}

func (s *AvalonPostgresStorage) GetPlayer(ctx context.Context, playerId avalon.PlayerId) (*avalon.PlayerImpl, error) {
	return nil, errors.New("Not implemented")
}

func (s *AvalonPostgresStorage) CreatePlayer(ctx context.Context, player *avalon.PlayerImpl) error {
	return errors.New("Not implemented")
}

func (s *AvalonPostgresStorage) UpdatePlayers(ctx context.Context, players []*avalon.PlayerImpl) error {
	return errors.New("Not implemented")
}

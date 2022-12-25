package storage

import (
	"context"

	avalon "github.com/yiwensong/ploggo/avalon"
)

type AvalonStorage interface {
	PlayerStorage
	GameStorage
}

type PlayerStorage interface {
	GetPlayersById(context.Context, []avalon.PlayerId) (map[avalon.PlayerId]*avalon.PlayerImpl, error)
	GetPlayer(context.Context, avalon.PlayerId) (*avalon.PlayerImpl, error)
	CreatePlayer(context.Context, *avalon.PlayerImpl) error
	UpdatePlayers(context.Context, []*avalon.PlayerImpl) error
}

type GameStorage interface {
	SaveGame(context.Context, *avalon.GameImpl) error
	GetGames(context.Context) (games []*avalon.GameImpl, err error)
}

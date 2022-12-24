package storage

import (
	avalon "github.com/yiwensong/ploggo/avalon"
)

type AvalonStorage interface {
	PlayerStorage
	GameStorage
}

type PlayerStorage interface {
	GetPlayersById([]avalon.PlayerId) (map[avalon.PlayerId]*avalon.PlayerImpl, error)
	GetPlayer(avalon.PlayerId) (*avalon.PlayerImpl, error)
	CreatePlayer(*avalon.PlayerImpl) error
	UpdatePlayers([]*avalon.PlayerImpl) error
}

type GameStorage interface {
	SaveGame(*avalon.GameImpl) error
	GetGames() (games []*avalon.GameImpl, err error)
}

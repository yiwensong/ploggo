package storage

import (
	avalon "github.com/yiwensong/ploggo/avalon"
)

type AvalonStorage interface {
	GetPlayersById([]avalon.PlayerId) (map[avalon.PlayerId]*avalon.PlayerImpl, error)
	GetPlayer(avalon.PlayerId) (*avalon.PlayerImpl, error)
	CreatePlayer(*avalon.PlayerImpl) error
	UpdatePlayers([]*avalon.PlayerImpl) error

	SaveGame(*avalon.GameImpl) error
	GetGames() (games []*avalon.GameImpl, err error)
}

const AVALON_PLAYER_JSON_FILE = "players.json"
const AVALON_GAME_JSON_FILE = "games.json"

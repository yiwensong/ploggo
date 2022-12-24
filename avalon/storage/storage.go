package storage

import (
	json "encoding/json"
	os "os"
	path "path"

	errors "github.com/pkg/errors"
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

type AvalonJsonStorage struct {
	BasePath    string
	PlayersById map[avalon.PlayerId]*avalon.PlayerImpl
	Games       []*avalon.GameImpl
}

func NewAvalonJsonStorage(basePath string) (j *AvalonJsonStorage, err error) {
	j = &AvalonJsonStorage{
		BasePath:    basePath,
		Games:       []*avalon.GameImpl{},
		PlayersById: map[avalon.PlayerId]*avalon.PlayerImpl{},
	}

	err = j.saveGameJson()
	if err != nil {
		return nil, errors.Wrapf(err, "saveGameJson")
	}

	err = j.savePlayerJson()
	if err != nil {
		return nil, errors.Wrapf(err, "savePlayerJson")
	}

	return j, nil
}

func LoadAvalonJsonStorageFromPath(basePath string) (j *AvalonJsonStorage, err error) {
	j = &AvalonJsonStorage{
		BasePath: basePath,
	}

	err = j.loadGameJson()
	if err != nil {
		return nil, errors.Wrapf(err, "loadGameJson")
	}

	err = j.loadPlayerJson()
	if err != nil {
		return nil, errors.Wrapf(err, "loadPlayerJson")
	}

	return j, nil
}

func (j *AvalonJsonStorage) loadPlayerJson() error {
	playerJsonPath := path.Join(j.BasePath, AVALON_PLAYER_JSON_FILE)

	playerData, err := os.ReadFile(playerJsonPath)
	if err != nil {
		return errors.Wrapf(err, "os.ReadFile(%q)", playerJsonPath)
	}

	json.Unmarshal(playerData, &j.PlayersById)
	return nil
}

func (j *AvalonJsonStorage) savePlayerJson() error {
	playerJsonPath := path.Join(j.BasePath, AVALON_PLAYER_JSON_FILE)

	f, err := os.OpenFile(playerJsonPath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return errors.Wrapf(err, "os.OpenFile(%q)", playerJsonPath)
	}
	defer f.Close()

	bytes, err := json.Marshal(j.PlayersById)

	if err != nil {
		return errors.Wrapf(err, "json.Marshal")
	}

	err = f.Truncate(0)
	if err != nil {
		return errors.Wrapf(err, "f.Truncate(0)")
	}

	_, err = f.Seek(0, 0)
	if err != nil {
		return errors.Wrapf(err, "f.Seek(0, 0)")
	}

	_, err = f.Write(bytes)
	if err != nil {
		return errors.Wrapf(err, "f.Write")
	}

	return nil
}

func (j *AvalonJsonStorage) loadGameJson() error {
	gameJsonPath := path.Join(j.BasePath, AVALON_GAME_JSON_FILE)

	gameData, err := os.ReadFile(gameJsonPath)
	if err != nil {
		return errors.Wrapf(err, "os.ReadFile(%q)", gameJsonPath)
	}

	json.Unmarshal(gameData, &j.Games)
	return nil
}

func (j *AvalonJsonStorage) saveGameJson() error {
	gameJsonPath := path.Join(j.BasePath, AVALON_GAME_JSON_FILE)

	f, err := os.OpenFile(gameJsonPath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return errors.Wrapf(err, "os.OpenFile(%q)", gameJsonPath)
	}
	defer f.Close()

	bytes, err := json.Marshal(j.Games)
	if err != nil {
		return errors.Wrapf(err, "json.Marshal")
	}

	err = f.Truncate(0)
	if err != nil {
		return errors.Wrapf(err, "f.Truncate(0)")
	}

	_, err = f.Seek(0, 0)
	if err != nil {
		return errors.Wrapf(err, "f.Seek(0, 0)")
	}

	_, err = f.Write(bytes)
	if err != nil {
		return errors.Wrapf(err, "f.Write")
	}

	return nil
}

func (j *AvalonJsonStorage) GetPlayersById(playerIds []avalon.PlayerId) (playersById map[avalon.PlayerId]*avalon.PlayerImpl, err error) {
	err = j.loadPlayerJson()
	if err != nil {
		return nil, errors.Wrapf(err, "loadPlayerJson")
	}

	return j.PlayersById, nil
}

func (j *AvalonJsonStorage) GetPlayer(playerId avalon.PlayerId) (player *avalon.PlayerImpl, err error) {
	players, err := j.GetPlayersById([]avalon.PlayerId{playerId})

	if err != nil {
		return nil, errors.Wrapf(err, "GetPlayersById(%q)", playerId)
	}

	player, ok := players[playerId]
	if !ok {
		return nil, errors.Errorf("Player %q does not exist", playerId)
	}

	return player, nil
}

func (j *AvalonJsonStorage) CreatePlayer(player *avalon.PlayerImpl) error {
	err := j.loadPlayerJson()
	if err != nil {
		return errors.Wrapf(err, "loadPlayerJson")
	}
	defer func() {
		err = j.savePlayerJson()
		if err != nil {
			err = errors.Wrapf(err, "savePlayerJson")
		}
	}()

	_, playerAlreadyExists := j.PlayersById[player.Id]
	if playerAlreadyExists {
		return errors.Errorf("Player already exists player_id=%q", player.Id)
	}

	j.PlayersById[player.Id] = player

	return err
}

func (j *AvalonJsonStorage) UpdatePlayers(players []*avalon.PlayerImpl) error {
	err := j.loadPlayerJson()
	if err != nil {
		return errors.Wrapf(err, "loadPlayerJson")
	}
	defer func() {
		err = j.savePlayerJson()
		if err != nil {
			err = errors.Wrapf(err, "savePlayerJson")
		}
	}()

	for _, player := range players {
		_, ok := j.PlayersById[player.Id]
		if !ok {
			return errors.Errorf("Player does not exist player_id=%q", player.Id)
		}

		j.PlayersById[player.Id] = player
	}

	return err
}

func (j *AvalonJsonStorage) SaveGame(game *avalon.GameImpl) error {
	err := j.loadGameJson()
	if err != nil {
		return errors.Wrapf(err, "loadGameJson")
	}
	defer func() {
		err = j.saveGameJson()
		if err != nil {
			err = errors.Wrapf(err, "saveGameJson")
		}
	}()

	j.Games = append(j.Games, game)

	return nil
}

func (j *AvalonJsonStorage) GetGames() (games []*avalon.GameImpl, err error) {
	err = j.loadGameJson()
	if err != nil {
		return nil, errors.Wrapf(err, "loadGameJson")
	}

	return j.Games, nil
}

var _ AvalonStorage = (*AvalonJsonStorage)(nil)

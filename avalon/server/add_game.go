package server

import (
	json "encoding/json"
	http "net/http"

	glog "github.com/golang/glog"
	errors "github.com/pkg/errors"
	avalon "github.com/yiwensong/ploggo/avalon"
)

type AddGameArgs struct {
	RoleByPlayerId map[avalon.PlayerId]avalon.Role
	Winner         avalon.Team
}

func (s *AvalonServer) AddGame(w http.ResponseWriter, req *http.Request) {
	var decodedArgs AddGameArgs
	ctx := req.Context()

	err := json.NewDecoder(req.Body).Decode(&decodedArgs)
	if err != nil {
		glog.Errorf("json.NewDecoder: %s", errors.Wrapf(err, "json.Decode"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var playerIds []avalon.PlayerId
	for playerId := range decodedArgs.RoleByPlayerId {
		playerIds = append(playerIds, playerId)
	}

	playersById, err := s.Storage.GetPlayersById(ctx, playerIds)
	if err != nil {
		glog.Errorf("GetPlayersById: %s", errors.Wrapf(err, "GetPlayersById"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	game := avalon.NewGame(playersById, decodedArgs.RoleByPlayerId)
	game.SetWinner(decodedArgs.Winner)

	updatedPlayers, err := game.UpdatePlayersAfterGame()
	if err != nil {
		glog.Errorf("UpdatePlayersAfterGame: %s", errors.Wrapf(err, "UpdatePlayersAfterGame"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = s.Storage.SaveGame(ctx, game); err != nil {
		glog.Errorf("SaveGame: %s", errors.Wrapf(err, "SaveGame"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = s.Storage.UpdatePlayers(ctx, updatedPlayers); err != nil {
		glog.Errorf("UpdatePlayers: %s", errors.Wrapf(err, "UpdatePlayers"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

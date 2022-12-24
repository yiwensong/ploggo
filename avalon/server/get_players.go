package server

import (
	json "encoding/json"
	http "net/http"

	glog "github.com/golang/glog"
	errors "github.com/pkg/errors"
	avalon "github.com/yiwensong/ploggo/avalon"
)

type GetPlayersArgs struct{}

func (s *AvalonServer) GetPlayers(w http.ResponseWriter, req *http.Request) {
	players, err := s.Storage.GetPlayersById([]avalon.PlayerId{})
	if err != nil {
		glog.Errorf(errors.Wrapf(err, "GetPlayersById").Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(players)
	if err != nil {
		glog.Errorf(errors.Wrapf(err, "json.Encode").Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

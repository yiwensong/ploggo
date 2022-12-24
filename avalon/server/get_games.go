package server

import (
	json "encoding/json"
	http "net/http"

	glog "github.com/golang/glog"
	errors "github.com/pkg/errors"
)

type GetGamesArgs struct{}

func (s *AvalonServer) GetGames(w http.ResponseWriter, req *http.Request) {
	games, err := s.Storage.GetGames()
	if err != nil {
		glog.Errorf(errors.Wrapf(err, "GetGames").Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(games)
	if err != nil {
		glog.Errorf(errors.Wrapf(err, "json.Encode").Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

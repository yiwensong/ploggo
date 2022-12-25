package server

import (
	json "encoding/json"
	http "net/http"

	glog "github.com/golang/glog"
	errors "github.com/pkg/errors"
	avalon "github.com/yiwensong/ploggo/avalon"
)

type CreatePlayerArgs struct {
	Name string
}

func (s *AvalonServer) CreatePlayer(w http.ResponseWriter, req *http.Request) {
	var decodedArgs CreatePlayerArgs
	ctx := req.Context()

	err := json.NewDecoder(req.Body).Decode(&decodedArgs)
	if err != nil {
		glog.Errorf(errors.Wrapf(err, "json.Decode").Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = s.Storage.CreatePlayer(ctx, avalon.NewPlayer(decodedArgs.Name))
	if err != nil {
		glog.Errorf(errors.Wrapf(err, "Avalon.NewPlayer(%s)", decodedArgs.Name).Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

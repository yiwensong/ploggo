package server

import (
	json "encoding/json"
	fmt "fmt"
	http "net/http"
	os "os"
	path "path"

	glog "github.com/golang/glog"
	errors "github.com/pkg/errors"
	storage "github.com/yiwensong/ploggo/avalon/storage"
)

func Healthcheck(w http.ResponseWriter, req *http.Request) {
	healthy, err := json.Marshal(map[string]string{
		"healthcheck": "healthy",
	})
	if err != nil {
		glog.Fatalf("json.Marshal")
	}

	w.Write(healthy)
}

type AvalonServer struct {
	Storage storage.AvalonStorage
}

func (s *AvalonServer) HandleGames(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		s.GetGames(w, req)
	case "POST":
		s.AddGame(w, req)
	default:
		http.Error(w, "Unexpected method", http.StatusMethodNotAllowed)
	}
}

func (s *AvalonServer) HandlePlayers(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		s.GetPlayers(w, req)
	case "POST":
		s.CreatePlayer(w, req)
	default:
		http.Error(w, "Unexpected method", http.StatusMethodNotAllowed)
	}
}

type ServerOpts struct {
	BasePath string
	Port     int64
}

func StartServer(opts *ServerOpts) error {
	mux := http.NewServeMux()

	basePath := opts.BasePath
	if basePath == "" {
		basePath = path.Join(os.Getenv("HOME"), ".avalon")
	}

	avalonStorage, err := storage.NewAvalonJsonStorage(basePath)
	if err != nil {
		return errors.Wrapf(err, "NewAvalonJsonStorage(%s)", basePath)
	}

	server := &AvalonServer{
		Storage: avalonStorage,
	}

	mux.HandleFunc("/healthcheck", Healthcheck)
	mux.HandleFunc("/game", server.HandleGames)
	mux.HandleFunc("/player", server.HandlePlayers)

	port := opts.Port
	if port == 0 {
		port = 4000
	}

	portString := fmt.Sprintf(":%d", port)

	err = http.ListenAndServe(portString, mux)
	if err != nil {
		return errors.Wrapf(err, "ListenAndServe(%s)", portString)
	}

	return nil
}

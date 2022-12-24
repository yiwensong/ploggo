package main

import (
	os "os"
	path "path"

	server "github.com/yiwensong/ploggo/avalon/server"
)

func main() {
	server.StartServer(&server.ServerOpts{
		BasePath: path.Join(os.Getenv("HOME"), ".avalon_test"),
	})
}

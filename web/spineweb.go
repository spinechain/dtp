package web

import (
	"context"
	"encoding/json"
	"net/http"
	"spinedtp/util"
	"time"
)

var server *http.Server

type NetworkInfo struct {
	NetworkName string
	ClientCount int
	TaskCount   int
}

// Start the website
func Start() {
	go ServeAPI()
}

// stop the website
func Stop() error {
	util.PrintBlue("Stopping web server")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	data := NetworkInfo{NetworkName: "SpineChainTestNet1", ClientCount: 0, TaskCount: 0}
	json.NewEncoder(w).Encode(data)
}

func ServeAPI() {

	util.PrintBlue("Starting web server")
	http.HandleFunc("/api/v1/info/", infoHandler)

	server = &http.Server{Addr: ":8080"}

	server.ListenAndServe()
}

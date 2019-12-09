package router

import (
	"github.com/BurntSushi/toml"
	"github.com/Sterks/ftpReader/config"
	"log"
	"net/http"
	"time"

	"github.com/Sterks/ftpReader/controller"
	"github.com/gorilla/mux"
)

// Start ...
func Start() (*http.Server, error) {
	var dir string
	r := mux.NewRouter()
	r.HandleFunc("/", controller.HomeHandler)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(dir))))
	r.HandleFunc("/download", controller.DownloadHandler)
	r.HandleFunc("/settings", controller.SettingsHandler)
	http.Handle("/", r)

	config := config.NewConfig()

	toml.DecodeFile("config/config.toml", config)

	srv := &http.Server{
		Handler:      r,
		Addr:         config.BindAddr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
	return srv, nil
}

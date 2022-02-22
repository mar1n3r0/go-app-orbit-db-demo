//go:build !wasm

// file: http_server.go

package main

import (
	"log"
	"net/http"

	"berty.tech/go-orbit-db/iface"
	"github.com/NYTimes/gziphandler"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	// This package is needed so that all the preloaded plugins are loaded automatically
)

var KvStore iface.KeyValueStore

func initServer() {
	withGz := gziphandler.GzipHandler(&app.Handler{
		Name:        "messenger",
		Description: "A orbit-db example with ipfs",
		Styles:      []string{"https://stackpath.bootstrapcdn.com/bootstrap/4.5.2/css/bootstrap.min.css"},
	})
	http.Handle("/", withGz)

	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
}

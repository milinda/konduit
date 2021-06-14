package main

import (
	"embed"
	"io/fs"
	"net/http"
)


//go:embed client/build
var content embed.FS

func rootHandler() http.Handler {
	fsys := fs.FS(content)
	contentStatic, _ := fs.Sub(fsys, "client/build")
	return http.FileServer(http.FS(contentStatic))
}

func serve() {
	mux := http.NewServeMux()
	mux.Handle("/", rootHandler())
	err := http.ListenAndServe(":3000", mux)
	if err != nil {
		panic("Failed to start")
	}
}
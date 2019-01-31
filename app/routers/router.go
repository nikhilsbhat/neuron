package router

import (
	"github.com/gorilla/mux"
	"io"
	"net/http"
	handle "neuron/app/handlers"
	mid "neuron/app/middleware"
	"github.com/gobuffalo/packr/v2"
)

type MuxIn struct {
	UiDir          string
	UiTemplatePath string
	Apilog         io.Writer
}

func (log *MuxIn) NewRouter() *mux.Router {

	if log.UiTemplatePath != "" {
		handle.UiTemplatePath = log.UiTemplatePath
	}
	rout := mux.NewRouter().StrictSlash(true)
	//rout.Use(mid.authenticate)
	rout.Use(mid.TimeoutHandler)

	//initializing logger with log path
	test := new(mid.Login)
	test.Logpath = log.Apilog
	rout.Use(test.Logger)

	for _, route := range handle.Routes {
		rout.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
		//rout.Use(mid.JsonHandler)
	}

	if log.UiDir != "" {
		for _, route := range handle.UiRoutes {
			rout.
				Path(route.Pattern).
				Name(route.Name).
				Handler(route.HandlerFunc)
			//rout.Use(mid.GzipHandler)
		}
		//rout.NotFoundHandler = http.HandlerFunc(notfound)
		box := packr.New("neuronBox", "/root/neuron/web")
		rout.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(box)))
		http.Handle("/", rout)
	} else {
		rout.HandleFunc("/", handle.Nouifound)
		http.Handle("/", rout)
	}

	return rout
}

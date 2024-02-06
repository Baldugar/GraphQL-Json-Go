package muxRouter

import (
	"graphql_json_go/graph"
	"graphql_json_go/graph/gentypes"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/gorilla/mux"
)

type HTMLDir struct {
	d http.Dir
}

func (d HTMLDir) Open(name string) (http.File, error) {
	// Try name as supplied
	f, err := d.d.Open(name)
	if os.IsNotExist(err) {
		// Not found, try with .html
		if f, err := d.d.Open(name + ".html"); err == nil {
			return f, nil
		}
	}
	return f, err
}

func CreateMux() *mux.Router {
	router := mux.NewRouter()

	graphQLServer := handler.NewDefaultServer(gentypes.NewExecutableSchema(gentypes.Config{Resolvers: &graph.Resolver{}}))

	// router.HandleFunc("/config.js", configHandler)            // for sending config.js to client

	// graphQL for client
	router.Handle("/next-public", graphQLServer)

	//router.HandleFunc("/___version", versionHandler)

	// this router serves the public folder after building the app
	admin := http.StripPrefix("/admin/", http.FileServer(http.Dir("./admin")))
	public := http.FileServer(HTMLDir{http.Dir("./public")})
	router.PathPrefix("/admin/").Handler(admin)
	router.PathPrefix("/").Handler(public)
	return router
}

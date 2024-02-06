package main

import (
	"flag"
	"graphql_json_go/settings"
	"graphql_json_go/util/logging"
	"graphql_json_go/util/muxRouter"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/rs/cors"
	"github.com/rs/zerolog/log"
)

type Flags = struct {
	settingsFile string
	listen       string
}

var flags Flags

func main() {

	//the server recover
	defer func() {
		r := recover()
		if r != nil {
			log.Error().Err(r.(error)).Bytes("stack", debug.Stack()).Msgf("Recovering panic")
		}
	}()

	// parse the command line arguments
	parseFlags()

	// load the settings
	settings.Load(flags.settingsFile)

	// set up the logging based on the settings
	logging.Configure(settings.Current.Logging)

	log.Info().Msg("########## Starting server ##########")

	router := muxRouter.CreateMux()

	var loggingRouter http.Handler
	loggingRouter = loggingHandler(router)

	if settings.Current.AllowCrossOrigin {
		c := cors.New(cors.Options{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{
				http.MethodHead,
				http.MethodGet,
				http.MethodPost,
				http.MethodPut,
				http.MethodPatch,
				http.MethodDelete,
			},
			AllowedHeaders:   []string{"*"},
			AllowCredentials: true,
		})
		loggingRouter = c.Handler(loggingRouter)
	}
	var srv *http.Server

	var httpListen string
	httpListen = settings.Current.HTTPListen
	if flags.listen != "" {
		httpListen = flags.listen
	}

	log.Info().Msgf("starting server on %s", httpListen)
	srv = &http.Server{
		Addr:              ":" + httpListen,
		Handler:           loggingRouter,
		ReadTimeout:       time.Duration(settings.Current.ReadTimeout) * time.Second,
		ReadHeaderTimeout: time.Duration(settings.Current.ReadHeaderTimeout) * time.Second,
		WriteTimeout:      time.Duration(settings.Current.WriteTimeout) * time.Second,
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start HTTP server")
	}
}

func parseFlags() {
	flag.StringVar(&flags.settingsFile, "settings", "", "Determines if the settings are loaded from a file. If it's not set the settings are loaded from environment variables")
	flag.StringVar(&flags.listen, "listen", "", "The port number to listen on")
	flag.Parse()
}

func loggingHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		next.ServeHTTP(w, r)
		t2 := time.Now()
		log.Debug().Str("method", r.Method).Str("path", r.URL.String()).Msgf("request executed in %v", t2.Sub(t1))
	}
	return http.HandlerFunc(fn)
}

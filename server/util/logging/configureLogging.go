package logging

import (
	"graphql_json_go/settings"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func Configure(cfg settings.LogConfig) {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimestampFieldName = "t"
	zerolog.LevelFieldName = "l"
	zerolog.MessageFieldName = "msg"
	zerolog.ErrorFieldName = "err"
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.DurationFieldUnit = time.Duration(zerolog.DurationFieldUnit.Milliseconds())
	zerolog.DurationFieldInteger = true
	zerolog.CallerMarshalFunc = func(file string, line int) string {
		s := strings.Split(file, "/")
		return s[len(s)-1] + ":" + strconv.Itoa(line)
	}

	if cfg.LogLevel == "Debug" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else if cfg.LogLevel == "Info" {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else if cfg.LogLevel == "Warn" {
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	} else if cfg.LogLevel == "Error" {
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	} else if cfg.LogLevel == "Fatal" {
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	}

	if cfg.LogFile != "" {
		//ensure the path to the logfile exist
		err := ensureDir(cfg.LogFilePath)
		if err != nil {
			log.Error().Err(err).Msg("loadSettings: could not ensure logfile path")
		}

		//write logs into a file
		dkd := filepath.Join(cfg.LogFilePath, cfg.LogFile)
		f, err := os.OpenFile(dkd, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Error().Err(err).Msg("error opening log file")
		} else {
			log.Logger = zerolog.New(f).With().Caller().Timestamp().Logger()
		}
		log.Info().Str("level", cfg.LogLevel).Str("path", dkd).Msg("Logging initialized")
	} else {
		log.Logger = log.With().Caller().Timestamp().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr})
		log.Info().Str("level", cfg.LogLevel).Msg("Logging initialized")
	}
}

// This function takes the given path and ensures that the directories of this path exist. If not the directories will be created.
func ensureDir(path string) error {
	_, err := os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	} else if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

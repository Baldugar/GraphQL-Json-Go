package settings

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
)

type LogConfig struct {
	LogLevel    string
	LogFile     string
	LogFilePath string
}

// Settings is the main struct that contains the configuration of the application
type Settings struct {
	AllowCrossOrigin  bool
	Logging           LogConfig
	Domain            string
	HTTPListen        string
	ReadHeaderTimeout int
	ReadTimeout       int
	WriteTimeout      int
}

var Current Settings

// Loads the settings
func Load(settingsFile string) {

	loadSettingsFromFile(settingsFile)

}

func loadSettingsFromFile(settingsFile string) {
	// read the settings file
	jsonData, err := ioutil.ReadFile(settingsFile)
	if os.IsNotExist(err) {
		log.Fatal().Err(err).Msg("error could not find settings file")
	} else if err != nil {
		log.Fatal().Err(err).Msg("error reading settings file")
	}

	// parse the json data into a struct
	var newSettings Settings
	err = json.Unmarshal(jsonData, &newSettings)
	if err != nil {
		log.Fatal().Err(err).Msg("error parsing settings into struct")
	}

	if isEmpty(newSettings.Domain) {
		log.Fatal().Err(err).Msg("error: Domain is missing in the settings file")
	}
	if isEmpty(newSettings.HTTPListen) {
		log.Fatal().Err(err).Msg("error: HTTPListen is missing in the settings file")
	}

	if isEmpty(newSettings.Logging.LogFilePath) {
		newSettings.Logging.LogFilePath = "./logs/"
	}
	if newSettings.ReadTimeout == 0 {
		newSettings.ReadTimeout = 120
	}
	if newSettings.ReadHeaderTimeout == 0 {
		newSettings.ReadHeaderTimeout = 40
	}
	if newSettings.WriteTimeout == 0 {
		newSettings.WriteTimeout = 120
	}

	Current = newSettings
}

func isEmpty(text string) bool {
	return len(strings.Trim(text, " ")) == 0
}

package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/koding/multiconfig"
	log "github.com/sirupsen/logrus"
)

// Config stores configuration of service
type Config struct {
	LogsFilesListJSON string            `required:"true"` // (example: '{"/log1.txt":"first_format ", "/dir/log2.log":"second_format "}'
	LogLevel          string            `default:"Info"`  // tool's logs level in stdout
	LogsFilesList     map[string]string `json:"-"`        // LogsFilesList store unmarshalled json  LogsFilesListJSON
	MongoURL          string            `required:"true"` // Mongo URL
	MongoUsername     string            `default:""`      // MongoUsername
	MongoPassword     string            `default:""`      // MongoPassword
	MongoDB           string            `default:"myDB"`  // Mongo DB name
	MongoCollection   string            `default:"logs"`  // Mongo DB collection
	DropDB            bool              `default:"false"` // if true - will dorp whole collection
}

// Help output for flags when program run with -h flag
func setFlagsHelp() map[string]string {
	usageMsg := make(map[string]string)

	usageMsg["LogsFilesListJSON"] = `JSON with list of all files that need to be looked at and converted
					example of JSON:
						{
							"/log1.txt":"first_format", 
							"/dir/log2.log":"second_format",
							"/dir2/log3.txt":"first_format"
						}
			`
	usageMsg["LogLevel"] = `LogLevel level: All, Debug, Info, Error, Fatal, Panic, Warn`

	usageMsg["MongoURL"] = "Mongo URL"
	usageMsg["MongoCollection"] = "Mongo DB collection"
	usageMsg["MongoDB"] = "Mongo DB name"
	usageMsg["DropDB"] = "if true - will drop whole collection before starting to store all logs"
	usageMsg["MongoPassword"] = "MongoDB Password"
	usageMsg["MongoUsername"] = "MongoDB Username"

	return usageMsg
}

// LoadConfig loads configuration struct from env vars, flags or from toml file
func LoadConfig(configPath string) (*Config, error) {

	svcConfig := new(Config)

	log.Infof("Loading configuration\n")

	m := newConfig(configPath, "LogsConverter", true)

	err := m.Load(svcConfig)
	if err != nil {

		log.Errorf("Failed to load configuration: %v", err)
		return nil, fmt.Errorf("failed to load configuration: %v", err)
	}

	setLogger(svcConfig)

	// Parse log files names with formats
	svcConfig.LogsFilesList = make(map[string]string)
	errUnMarshal := json.Unmarshal([]byte(svcConfig.LogsFilesListJSON), &svcConfig.LogsFilesList)
	if errUnMarshal != nil {
		log.Errorf("Failed to unmarshal json with files [%s] to struct: %v", svcConfig.LogsFilesListJSON, errUnMarshal)
		return nil, fmt.Errorf("failed to unmarshal json with files [%s] to struct: %v", svcConfig.LogsFilesListJSON, errUnMarshal)
	}

	err = m.Validate(svcConfig)
	if err != nil {

		log.Errorf("Config struct is invalid: %v\n", err)
		return nil, fmt.Errorf("config struct is invalid: %v", err)
	}
	if len(svcConfig.LogsFilesList) == 0 {

		log.Errorf("No log files provided: [%+v]", svcConfig.LogsFilesList)
		return nil, fmt.Errorf("No log files provided: [%+v]", svcConfig.LogsFilesList)
	}

	log.Infof("Configuration loaded\n")

	prettyConfig, err := json.MarshalIndent(svcConfig, "", "")
	if err != nil {

		log.Errorf("Failed to marshal indent config: %v", err)
		return nil, fmt.Errorf("failed to marshal indent config: %v", err)

	}
	log.Infof("Current config:\n %s", string(prettyConfig))

	return svcConfig, nil

}

// Implementation of default loader for multiconfig
func newConfig(path string, prefix string, camelCase bool) *multiconfig.DefaultLoader {
	var loaders []multiconfig.Loader

	// Read default values defined via tag fields "default"
	loaders = append(loaders, &multiconfig.TagLoader{})

	if path != "" {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			log.Warnf("Provided local config file [%s] does not exist. Flags and Environment variables will be used ", path)
		} else {
			log.Infof("Local config file [%s] will be used", path)
			// Choose what while is passed
			if strings.HasSuffix(path, "toml") {
				log.Debugf("Toml detected")
				loaders = append(loaders, &multiconfig.TOMLLoader{Path: path})
			}

		}
	}

	e := &multiconfig.EnvironmentLoader{
		Prefix:    prefix,
		CamelCase: camelCase,
	}

	usageMsg := setFlagsHelp()
	f := &multiconfig.FlagLoader{
		Prefix:        "",
		Flatten:       false,
		CamelCase:     camelCase,
		EnvPrefix:     prefix,
		ErrorHandling: 0,
		Args:          nil,
		FlagUsageFunc: func(s string) string { return usageMsg[s] },
	}

	loaders = append(loaders, e, f)
	loader := multiconfig.MultiLoader(loaders...)

	d := &multiconfig.DefaultLoader{}
	d.Loader = loader
	d.Validator = multiconfig.MultiValidator(&multiconfig.RequiredValidator{})
	return d

}

func setLogger(cfg *Config) {
	formatter := &log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		DisableSorting:  false,
		ForceColors:     true,
	}
	log.SetFormatter(formatter)
	lvl, err := log.ParseLevel(cfg.LogLevel)
	if err != nil {
		log.Warnf("Could not parse log level [%s], will be used Info", cfg.LogLevel)
		lvl = log.InfoLevel
	}
	log.SetLevel(lvl)

}

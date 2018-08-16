package config

import (
	"encoding/json"
	"os"
	"path"
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

	return usageMsg
}

// LoadConfig loads configuration struct from env vars, flags or from toml file
func LoadConfig() *Config {

	svcConfig := new(Config)

	log.Infof("Loading configuration\n")

	configPath := path.Join("config.toml")

	m := newConfig(configPath, "KafkaDump", false)

	err := m.Load(svcConfig)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	logLevel, err := log.ParseLevel(svcConfig.LogLevel)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(logLevel)

	// Parse log files names with formats
	svcConfig.LogsFilesList = make(map[string]string)
	errUnMarshal := json.Unmarshal([]byte(svcConfig.LogsFilesListJSON), &svcConfig.LogsFilesList)
	if errUnMarshal != nil {
		log.Fatal(errUnMarshal)
	}

	err = m.Validate(svcConfig)
	if err != nil {

		log.Fatalf("Config struct is invalid: %v\n", err)
	}

	log.Infof("Configuration loaded\n")
	prettyConfig, err := json.MarshalIndent(svcConfig, "", "")
	if err != nil {
		log.Fatalf("Failed to marshal indent config: %v", err)

	}
	log.Infof("Current config:\n %s", string(prettyConfig))

	return svcConfig

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

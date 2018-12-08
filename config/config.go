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
	logsFilesList     map[string]string `json:"-"`        // LogsFilesList store unmarshalled json  LogsFilesListJSON
	DBURL             string            `required:"true"` // Database URL
	DBUsername        string            `default:""`      // Database Username
	DBPassword        string            `default:""`      // DBPassword
	DBName            string            `default:"myDB"`  // DB name
	MongoCollection   string            `default:"logs"`  // Mongo DB collection
	DropDB            bool              `default:"false"` // if true - will dorp whole collection
	FollowFiles       bool              `default:"true"`  // if true - will tail file and wait for updates
	FilesMustExist    bool              `default:"true"`  // if true - will throw error when file is not exist; when false - wait for file create

}

// Help output for flags when program run with -h flag
func setFlagsHelp() map[string]string {
	usageMsg := make(map[string]string)

	usageMsg["logs-files-list-json"] = `JSON with list of all files that need to be looked at and converted
								example of JSON:
									{
										"/log1.txt":"first_format", 
										"/dir/log2.log":"second_format",
										"/dir2/log3.txt":"first_format"
									}`
	usageMsg["LogLevel"] = `LogLevel level: All, Debug, Info, Error, Fatal, Panic, Warn`
	usageMsg["DBURL"] = "Mongo URL"
	usageMsg["MongoCollection"] = "Mongo DB collection"
	usageMsg["DBName"] = "Mongo DB name"
	usageMsg["DropDB"] = "if true - will drop whole collection before starting to store all logs"
	usageMsg["DBPassword"] = "DBName Password"
	usageMsg["DBUsername"] = "DBName Username"
	usageMsg["FollowFiles"] = ` if true - will tail file and wait for updates; when false - end file reading after EOF`
	usageMsg["FilesMustExist"] = `if true - will throw error when file is not exist; when false - wait for file create`

	return usageMsg
}

// GetFilesList returns list of log files with filename and format
func (cfg *Config) GetFilesList() map[string]string {
	return cfg.logsFilesList
}

// LoadConfig loads configuration struct from env vars, flags or from toml file
func LoadConfig(configPath string) (*Config, error) {

	svcConfig := new(Config)

	log.Infof("Loading configuration\n")

	m := newConfig(configPath, "LogsConverter", true)

	err := m.Load(svcConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %v", err)
	}

	setLogger(svcConfig)

	// Parse log files names with formats
	svcConfig.logsFilesList, err = parseLogsFilesList(svcConfig.LogsFilesListJSON)
	if err != nil {
		return nil, err
	}

	err = m.Validate(svcConfig)
	if err != nil {
		return nil, fmt.Errorf("config struct is invalid: %v", err)
	}
	if len(svcConfig.logsFilesList) == 0 {
		return nil, fmt.Errorf("no log files provided: [%+v]", svcConfig.logsFilesList)
	}

	log.Infof("Configuration loaded\n")

	prettyConfig, err := json.MarshalIndent(svcConfig, "", "")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal indent config: %v", err)
	}
	log.Infof("Current config:\n %s", string(prettyConfig))

	return svcConfig, nil

}

func parseLogsFilesList(filesListJSON string) (map[string]string, error) {
	filesList := make(map[string]string)

	errUnMarshal := json.Unmarshal([]byte(filesListJSON), &filesList)
	if errUnMarshal != nil {
		return nil, fmt.Errorf("failed to unmarshal json with files [%s] to struct: %v", filesListJSON, errUnMarshal)
	}
	return filesList, nil

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

package config

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testConfigCreate(t *testing.T, configPath string, logsFilesListJSON string, logLevel string,
	mongoURL string, mongoDB string, mongoCollection string, dropDB bool) error {

	fmt.Println("Helper func in action")
	if err := os.MkdirAll(filepath.Dir(configPath), 0700); err != nil {
		t.Fatalf("failed creating all dirs for config file [%s]: %v", filepath.Dir(configPath), err)

	}

	configFile, err := os.OpenFile(configPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)

	if err != nil {
		t.Fatalf("error opening config file: %v", err)
	}

	if err != nil {
		return fmt.Errorf("failed to create test config: %v", err)
	}

	_, err = configFile.WriteString(fmt.Sprintf(`LogLevel="%s"
	LogsFilesListJSON='%s'
	DBURL="%s"
	DBName="%s"
	MongoCollection="%s"
	DropDB=%t`, logLevel, logsFilesListJSON, mongoURL, mongoDB, mongoCollection, dropDB))

	if err != nil {
		return fmt.Errorf("failed to write test config: %v", err)
	}
	if err := configFile.Close(); err != nil {
		return fmt.Errorf("failed to close test config file descriptor: %v", err)
	}

	return nil
}

func testConfigDelete(t *testing.T, configPath string) {
	err := os.Remove(configPath)
	if err != nil {
		t.Fatalf("Failed to delete old config file")
	}
}

func TestLoadConfig(t *testing.T) {
	type input struct {
		logsFilesListJSON string
		logLevel          string
		DBURL             string
		DBName            string
		DBUsername        string
		DBPassword        string
		mongoCollection   string
		dropDB            bool
		followFiles       bool
		fileMustExist     bool
	}
	type expectedResult struct {
		wantConfig *Config
		wantErr    bool
	}
	var tests = []struct {
		id             int
		description    string
		input          input
		expectedResult expectedResult
	}{
		{
			id:          1,
			description: `Check configuration loading from cofig file`,
			input: input{
				logsFilesListJSON: `{"testdata/testfile1.log":"second_format",
"testdata/dir1/testfile2.log":"first_format"}`,
				logLevel:        "Info",
				DBURL:           "localhost:27017",
				DBUsername:      "",
				DBPassword:      "",
				DBName:          "myDB",
				mongoCollection: "logs",
				dropDB:          true,
				followFiles:     true,
				fileMustExist:   true,
			},
			expectedResult: expectedResult{
				wantConfig: &Config{
					LogsFilesListJSON: `{"testdata/testfile1.log":"second_format",
"testdata/dir1/testfile2.log":"first_format"}`,
					LogLevel:        "Info",
					DBURL:           "localhost:27017",
					DBUsername:      "",
					DBPassword:      "",
					DBName:          "myDB",
					MongoCollection: "logs",
					DropDB:          true,
					logsFilesList: map[string]string{"testdata/testfile1.log": "second_format",
						"testdata/dir1/testfile2.log": "first_format"},
					FilesMustExist: true,
					FollowFiles:    true,
				},
				wantErr: false,
			},
		},
		{
			id:          2,
			description: `Broken config: incorrect json with files`,
			input: input{
				logsFilesListJSON: `{"testdata/testfile1.log":"second_format",
"testdata/dir1/testfile2.log":"first_format`,
				logLevel:        "Debug",
				DBURL:           "localhost:29800",
				DBUsername:      "",
				DBPassword:      "",
				DBName:          "myDB1",
				mongoCollection: "logs1",
				dropDB:          true,
			},
			expectedResult: expectedResult{
				wantConfig: nil,
				wantErr:    true,
			},
		},
		{
			id:          3,
			description: `Broken config: empty json with files`,
			input: input{
				logsFilesListJSON: `{}`,
				logLevel:          "Warn",
				DBURL:             "localhost:2980",
				DBUsername:        "",
				DBPassword:        "",
				DBName:            "myDB2",
				mongoCollection:   "logs2",
				dropDB:            true,
			},
			expectedResult: expectedResult{
				wantConfig: nil,
				wantErr:    true,
			},
		},
	}

	currentDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		t.Fatal(err)
	}
	configPath := filepath.Join(currentDir, "testdata", "config.toml")

	for _, tc := range tests {
		tc := tc

		t.Run(fmt.Sprintf("Test%d:%s", tc.id, tc.description), func(t *testing.T) {

			err = testConfigCreate(t, configPath, tc.input.logsFilesListJSON, tc.input.logLevel, tc.input.DBURL,
				tc.input.DBName, tc.input.mongoCollection, tc.input.dropDB)

			if err != nil {
				t.Fatalf("Error while creating test config: %v", err)
			}

			gotModel, err := LoadConfig(configPath)

			switch tc.expectedResult.wantErr {
			case true:
				assert.Error(t, err, "Expected to receive error from LoadConfig()")

			case false:
				assert.NoError(t, err, "Unexpected error from LoadConfig()")

			}
			assert.Equal(t, tc.expectedResult.wantConfig, gotModel)

			testConfigDelete(t, configPath) // Delete old config file after it was loaded

		})
	}
}

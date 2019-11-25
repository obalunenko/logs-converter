package config

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

type expectedResult struct {
	wantConfig *Config
	wantErr    bool
}

type test struct {
	id             int
	description    string
	inputFile      string
	expectedResult expectedResult
}

func TestLoadConfig(t *testing.T) {
	for _, tc := range tests() {
		tc := tc
		t.Run(fmt.Sprintf("Test%d:%s", tc.id, tc.description), func(t *testing.T) {
			gotModel, err := LoadConfig(tc.inputFile)
			if tc.expectedResult.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedResult.wantConfig, gotModel)
		})
	}
}

func tests() []test {
	return []test{
		{
			id:          1,
			description: `Check configuration loading from cofig file`,
			inputFile:   filepath.Join("testdata", "valid-config.toml"),
			expectedResult: expectedResult{
				wantConfig: &Config{
					LogsFilesListJSON: "{\"testdata/testfile1.log\":\"second_format\"," +
						"\"testdata/dir1/testfile2.log\":\"first_format\"}",
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
			inputFile:   filepath.Join("testdata", "broken-config.toml"),
			expectedResult: expectedResult{
				wantConfig: nil,
				wantErr:    true,
			},
		},
		{
			id:          3,
			description: `Broken config: empty json with files`,
			inputFile:   filepath.Join("testdata", "broken-config2.toml"),
			expectedResult: expectedResult{
				wantConfig: nil,
				wantErr:    true,
			},
		},
	}
}

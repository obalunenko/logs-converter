package converter

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/oleg-balunenko/logs-converter/internal/models"
)

func Test_processLine(t *testing.T) {
	type input struct {
		logName    string
		line       string
		format     string
		lineNumber uint64
	}

	type expectedResult struct {
		wantModel *models.LogModel
		wantErr   bool
	}

	var tests = []struct {
		id             int
		description    string
		input          input
		expectedResult expectedResult
	}{
		{
			id:          1,
			description: `Invalid format of line. Separator "|"  not found`,
			input: input{
				logName:    "test",
				line:       `Feb 1, 2018 at 3:04:05pm (UTC)  This is log message`,
				format:     "first_format",
				lineNumber: 1,
			},
			expectedResult: expectedResult{
				wantModel: nil,
				wantErr:   true,
			},
		},
		{
			id:          2,
			description: `Positive case. First format - one "|" separator`,
			input: input{
				logName:    "test",
				line:       `Feb 1, 2018 at 3:04:05pm (UTC) | This is log message`,
				format:     "first_format",
				lineNumber: 1,
			},
			expectedResult: expectedResult{
				wantModel: &models.LogModel{
					LogTime:   time.Date(2018, 02, 01, 15, 04, 05, 0, time.UTC),
					LogMsg:    `This is log message`,
					LogFormat: `first_format`,
					FileName:  "test",
				},
				wantErr: false,
			},
		},
		{
			id:          3,
			description: `Positive case. First format - more that one  "|" separator`,
			input: input{
				logName:    "test",
				line:       `Feb 1, 2018 at 3:04:05pm (UTC) | This is log message | that has|more than one separator`,
				format:     "first_format",
				lineNumber: 1,
			},
			expectedResult: expectedResult{
				wantModel: &models.LogModel{
					LogTime:   time.Date(2018, 02, 01, 15, 04, 05, 0, time.UTC),
					LogMsg:    `This is log message | that has|more than one separator`,
					LogFormat: `first_format`,
					FileName:  "test",
				},
				wantErr: false,
			},
		},
		{
			id:          4,
			description: `Positive case. Second format`,
			input: input{
				logName:    "test",
				line:       `2018-02-01T15:04:05Z | This is log message`,
				format:     "second_format",
				lineNumber: 1,
			},
			expectedResult: expectedResult{
				wantModel: &models.LogModel{
					LogTime:   time.Date(2018, 02, 01, 15, 04, 05, 0, time.UTC),
					LogMsg:    `This is log message`,
					LogFormat: `second_format`,
					FileName:  "test",
				},
				wantErr: false,
			},
		},
		{
			id:          5,
			description: `Negative case. Format missmatch - time format is second, but in config specified that file has first`,
			input: input{
				logName:    "test",
				line:       `2018-02-01T15:04:05Z | This is log message`,
				format:     "first_format",
				lineNumber: 1,
			},
			expectedResult: expectedResult{
				wantModel: nil,
				wantErr:   true,
			},
		},
		{
			id:          6,
			description: `Negative case. Not supported format specified in config`,
			input: input{
				logName:    "test",
				line:       `02/01/06 03:04:05 PM Jan | This is log message`,
				format:     "third_format",
				lineNumber: 1,
			},
			expectedResult: expectedResult{
				wantModel: nil,
				wantErr:   true,
			},
		},
		{
			id:          7,
			description: `Negative case. Format missmatch - time format is first, but in config specified that file has second`,
			input: input{
				logName:    "test",
				line:       `Feb 1, 2018 at 3:04:05pm (UTC) | This is log message`,
				format:     "second_format",
				lineNumber: 1,
			},
			expectedResult: expectedResult{
				wantModel: nil,
				wantErr:   true,
			},
		},
		{
			id:          8,
			description: `Negative case. Empty line received`,
			input: input{
				logName:    "test",
				line:       ``,
				format:     "second_format",
				lineNumber: 1,
			},
			expectedResult: expectedResult{
				wantModel: nil,
				wantErr:   true,
			},
		},
		{
			id:          9,
			description: `Negative case. Failed to parse time`,
			input: input{
				logName:    "test",
				line:       `2018-02-01T25:04:05Z | This is log message`,
				format:     "second_format",
				lineNumber: 1,
			},
			expectedResult: expectedResult{
				wantModel: nil,
				wantErr:   true,
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(fmt.Sprintf("Test%d:%s", tc.id, tc.description), func(t *testing.T) {
			gotModel, err := processLine(tc.input.logName, tc.input.line, tc.input.format, tc.input.lineNumber)

			switch tc.expectedResult.wantErr {
			case true:
				assert.Error(t, err, "Expected to receive error from processLine()")
			case false:
				assert.NoError(t, err, "Unexpected error from processLine()")
			}

			assert.Equal(t, tc.expectedResult.wantModel, gotModel)
		})
	}
}

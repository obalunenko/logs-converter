package converter

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func Test_processLine(t *testing.T) {
	type input struct {
		logName    string
		line       string
		format     string
		lineNumber uint64
	}
	type expectedResult struct {
		wantModel *LogModel
		wantErr   bool
	}
	type testCase struct {
		id             int
		description    string
		input          input
		expectedResult expectedResult
	}

	type testSuite []testCase

	var ts = testSuite{
		testCase{
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
		testCase{
			id:          2,
			description: `Positive case. First format - one "|" separator`,
			input: input{
				logName:    "test",
				line:       `Feb 1, 2018 at 3:04:05pm (UTC) | This is log message`,
				format:     "first_format",
				lineNumber: 1,
			},
			expectedResult: expectedResult{
				wantModel: &LogModel{
					LogTime:   time.Date(2018, 02, 01, 15, 04, 05, 0, time.UTC),
					LogMsg:    `This is log message`,
					LogFormat: `first_format`,
					FileName:  "test",
				},
				wantErr: false,
			},
		},
		testCase{
			id:          3,
			description: `Positive case. First format - more that one  "|" separator`,
			input: input{
				logName:    "test",
				line:       `Feb 1, 2018 at 3:04:05pm (UTC) | This is log message | that has|more than one separator`,
				format:     "first_format",
				lineNumber: 1,
			},
			expectedResult: expectedResult{
				wantModel: &LogModel{
					LogTime:   time.Date(2018, 02, 01, 15, 04, 05, 0, time.UTC),
					LogMsg:    `This is log message | that has|more than one separator`,
					LogFormat: `first_format`,
					FileName:  "test",
				},
				wantErr: false,
			},
		},
		testCase{
			id:          4,
			description: `Positive case. Second format`,
			input: input{
				logName:    "test",
				line:       `2018-02-01T15:04:05Z | This is log message`,
				format:     "second_format",
				lineNumber: 1,
			},
			expectedResult: expectedResult{
				wantModel: &LogModel{
					LogTime:   time.Date(2018, 02, 01, 15, 04, 05, 0, time.UTC),
					LogMsg:    `This is log message`,
					LogFormat: `second_format`,
					FileName:  "test",
				},
				wantErr: false,
			},
		},
		testCase{
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
		testCase{
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
		testCase{
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
		testCase{
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
		testCase{
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

	for _, tc := range ts {
		//if tc.id == 8 {
		t.Run(fmt.Sprintf("Test%d:%s", tc.id, tc.description), func(t *testing.T) {
			gotModel, err := processLine(tc.input.logName, tc.input.line, tc.input.format, tc.input.lineNumber)
			if (err != nil) != tc.expectedResult.wantErr {
				t.Errorf("processLine() error = %+v, \nwantErr %+v", err, tc.expectedResult.wantErr)
				return
			}
			if !reflect.DeepEqual(gotModel, tc.expectedResult.wantModel) {
				t.Errorf("processLine() = %+v, \nwant %+v", gotModel, tc.expectedResult.wantModel)
			}
		})
	}
	//}
}

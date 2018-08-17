package converter

import (
	"fmt"
	"strings"
	"time"

	"github.com/hpcloud/tail"
	log "github.com/sirupsen/logrus"
)

// Start starts converting of logfile
func Start(logName string, format string, resultChan chan *LogModel) {
	log.Infof("Starting tailing and converting file [%s] with logs format [%s]", logName, format)

	t, err := tail.TailFile(logName, tail.Config{
		Follow: true,
	})
	if err != nil {
		log.Fatalf("failed to tail file [%s]: %v", logName, err)
	}

	var cnt uint64
	for line := range t.Lines {
		cnt++
		log.Debugf("File:[%s] Line tailed: [%v]", logName, line)
		model, err := processLine(logName, line.Text, format, cnt)
		if err != nil {
			log.Errorf("Failed to process line [%s]: %v", line.Text, err)
		}

		log.Debugf("Go routine for file [%s] sending model to chanel", logName)
		resultChan <- model
		log.Debugf("Go routine for file [%s] sent model to chanel", logName)

	}

}

func processLine(logName string, line string, format string, lineNumber uint64) (model *LogModel, err error) {

	lineElements := strings.Split(line, " | ")
	if len(lineElements) < 1 {
		log.Errorf("[%s]: Line [%d] has wrong log structure: %s", logName, lineNumber, line)
		return nil, err
	}
	logTime, err := parseTime(lineElements[0], format)
	if err != nil {
		return nil, err
	}

	var msg string
	// For case when log message contain additional " | " to not miss other part of message
	if len(lineElements) > 2 {
		msg = strings.Join(lineElements[1:], " | ")
	} else {
		msg = lineElements[1]
	}

	model = &LogModel{
		LogTime:   logTime,
		LogMsg:    msg,
		FileName:  logName,
		LogFormat: format,
	}

	return model, err
}

// parseTime parses logTime string as format that was passed and return time.Time representation of it
func parseTime(logTimeStr string, format string) (time.Time, error) {

	switch format {
	case firstFormat:
		logTime, errParse := time.Parse(firstFormatLayout, logTimeStr)
		if errParse != nil {
			log.Errorf("failed to parse logTime [%s] as format [%s]: %v", logTimeStr, format, errParse)
			return time.Time{}, fmt.Errorf("failed to parse logTime [%s] as format [%s]: %v", logTimeStr, format, errParse)
		}

		return logTime, nil
	case secondFormat:
		logTime, errParse := time.Parse(secondFormatLayout, logTimeStr)
		if errParse != nil {
			log.Errorf("failed to parse logTime [%s] as format [%s]: %v", logTimeStr, format, errParse)
			return time.Time{}, fmt.Errorf("failed to parse logTime [%s] as format [%s]: %v", logTimeStr, format, errParse)
		}

		return logTime, nil
	default:
		log.Errorf("unexpected time format received (%s)", logTimeStr)
		return time.Time{}, fmt.Errorf("unexpected time format received (%s)", logTimeStr)
	}

}

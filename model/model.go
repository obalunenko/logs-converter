package model

import "time"

// LogModel to store log line
type LogModel struct {
	ID        string    `bson:"_id"`
	LogTime   time.Time `bson:"log_time"`
	LogMsg    string    `bson:"log_msg"`
	FileName  string    `bson:"file_name"`
	LogFormat string    `bson:"log_format"`
}

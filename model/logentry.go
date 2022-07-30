package model

import "time"

type LogEntry struct {
	RequestTime  time.Time
	ClientIP     string
	ClientName   string
	Reason       string
	ResponseType string
	QuestionType string
	QuestionName string
}

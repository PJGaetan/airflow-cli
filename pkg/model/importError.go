package model

import (
	"time"
)

type ImportErrors struct {
	ImportError  []ImportError `json:"import_errors"`
	TotalEntries int           `json:"total_entries"`
}

type ImportError struct {
	ImportErrrorId int       `json:"import_error_id"`
	Timestamp      time.Time `json:"timestamp"`
	Filename       string    `json:"finlename"`
	Stacktrace     bool      `json:"stack_trace"`
}

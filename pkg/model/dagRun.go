package model

import (
	"encoding/json"
	"time"
)

type DagRuns struct {
	DagRun       []DagRun `json:"dag_runs"`
	TotalEntries int      `json:"total_entries"`
}

type DagRun struct {
	Conf            json.RawMessage `json:"conf"`
	DagId           string          `json:"dag_id"`
	DagRunId        string          `json:"dag_run_id"`
	ExternalTrigger bool            `json:"external_trigger"`
	State           string          `json:"state"`
	ExecutionDate   time.Time       `json:"execution_date"`
	LogicalDate     time.Time       `json:"logical_date"`
	StartDate       time.Time       `json:"start_date"`
	EndDate         time.Time       `json:"end_date"`
}

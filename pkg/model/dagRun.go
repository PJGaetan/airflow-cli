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
	Conf             json.RawMessage `json:"conf"`
	Dag_id           string          `json:"dag_id"`
	Dag_run_id       string          `json:"dag_run_id"`
	External_trigger bool            `json:"external_trigger"`
	State            string          `json:"state"`
	End_date         time.Time       `json:"end_date"`
	Execution_date   time.Time       `json:"execution_date"`
	Logical_date     time.Time       `json:"logical_date"`
	Start_date       time.Time       `json:"start_date"`
}

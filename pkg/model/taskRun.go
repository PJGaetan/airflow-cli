package model

import (
	"encoding/json"
	"time"
)

type TaskInstances struct {
	TaskInstance []TaskInstance `json:"task_instances"`
	TotalEntries int            `json:"total_entries"`
}

type TaskInstance struct {
	DagId          string          `json:"dag_id"`
	Duration       float64         `json:"duration"`
	EndDate        time.Time       `json:"end_date"`
	ExecutionDate  time.Time       `json:"execution_date"`
	ExecutorConfig json.RawMessage `json:"executor_config"`
	Hostname       string          `json:"hostname"`
	MaxTries       float64         `json:"max_tries"`
	Operator       string          `json:"operator"`
	Pid            float64         `json:"pid"`
	Pool           string          `json:"pool"`
	PoolSlots      float64         `json:"pool_slots"`
	PriorityWeight float64         `json:"priority_weight"`
	Queue          string          `json:"queue"`
	QueuedWhen     string          `json:"queued_when"`
	Sla_miss       string          `json:"sla_miss"`
	StartDate      time.Time       `json:"start_date"`
	State          string          `json:"state"`
	TaskId         string          `json:"task_id"`
	TryNumber      float64         `json:"try_number"`
	Unixname       string          `json:"unixname"`
}

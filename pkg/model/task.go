package model

import (
	"encoding/json"
	"time"
)

type Tasks struct {
	Task         []Task `json:"tasks"`
	TotalEntries int    `json:"total_entries"`
}

type ClassRef struct {
	ClassName  string `json:"class_name"`
	ModulePath string `json:"module_path"`
}

type Task struct {
	ClassRef                ClassRef        `json:"class_ref"`
	DependsOnPast           bool            `json:"depends_on_past"`
	DownstreamTaskIds       []string        `json:"downstream_task_ids"`
	EndDate                 time.Time       `json:"end_date"`
	ExecutionTimeout        string          `json:"execution_timeout"`
	ExtraLinks              json.RawMessage `json:"extra_links"`
	Owner                   string          `json:"owner"`
	Params                  json.RawMessage `json:"params"`
	Pool                    string          `json:"pool"`
	PoolSlots               float32         `json:"pool_slots"`
	PriorityWeight          float32         `json:"priority_weight"`
	Queue                   string          `json:"queue"`
	Retries                 float32         `json:"retries"`
	RetryDelay              json.RawMessage `json:"retry_delay"`
	RetryExponentialBackoff bool            `json:"retry_exponential_backoff"`
	StartDate               time.Time       `json:"start_date"`
	TaskId                  string          `json:"task_id"`
	TemplateFields          json.RawMessage `json:"template_fields"`
	TriggerRule             string          `json:"trigger_rule"`
	UiColor                 string          `json:"ui_color"`
	UiFgcolor               string          `json:"ui_fgcolor"`
	WaitForDownstream       bool            `json:"wait_for_downstream"`
	WeightRule              string          `json:"weight_rule"`
}

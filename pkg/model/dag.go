package model

import "time"

type Dags struct {
	Dags         []Dag `json:"dags"`
	TotalEntries int   `json:"total_entries"`
}

type scheduleInterval struct {
	Type  string `json:"__type"`
	Value string `json:"value"`
}

type Tag struct {
	Tag string `json:"name"`
}

type Dag struct {
	DagId            string           `json:"dag_id"`
	Description      string           `json:"description"`
	FileToken        string           `json:"file_token"`
	Fileloc          string           `json:"fileloc"`
	IsActive         bool             `json:"is_active"`
	IsPaused         bool             `json:"is_paused"`
	IsSubdag         bool             `json:"s_subdag"`
	Owners           []string         `json:"owners"`
	RootDagId        string           `json:"root_dag_id"`
	ScheduleInterval scheduleInterval `json:"schedule_interval"`
	Tags             []Tag            `json:"tags"`
	DefaultView      string           `json:"default_view"`
	NextDagrun       time.Time        `json:"next_dagrun"`
}

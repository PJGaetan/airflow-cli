package model

type Dags struct {
	Dags         []Dag `json:"dags"`
	TotalEntries int   `json:"total_entries"`
}

type ResponseDag struct {
	Response Dags   `json:"response"`
	Status   string `json:"stat"`
}

type scheduleInterval struct {
	Type  string `json:"__type"`
	Value string `json:"value"`
}

type Tag struct {
	Tag map[string]string `json:"value"`
}

type Dag struct {
	Dag_id            string           `json:"dag_id"`
	Description       string           `json:"description"`
	File_token        string           `json:"file_token"`
	Fileloc           string           `json:"fileloc"`
	Is_active         bool             `json:"is_active"`
	Is_paused         bool             `json:"is_paused"`
	S_subdag          bool             `json:"s_subdag"`
	Owners            []string         `json:"owners"`
	Root_dag_id       string           `json:"root_dag_id"`
	Schedule_interval scheduleInterval `json:"schedule_interval"`
	Tags              []Tag            `json:"tags"`
}

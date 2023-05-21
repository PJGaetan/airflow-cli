package model

import "encoding/json"

type Logs struct {
	Content            json.RawMessage `json:"content"`
	ContinutationToken string          `json:"continuation_token"`
}

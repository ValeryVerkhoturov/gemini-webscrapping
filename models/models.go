package models

import "encoding/json"

type Review struct {
	Name    string   `json:"name"`
	Message string   `json:"message"`
	Mark    *float32 `json:"mark,omitempty"`
	Date    string   `json:"date"` // YYYY-MM-DD
}

type Reviews []*Review

func (r *Reviews) UnmarshalLLMText(data string) error {
	reviews := Reviews{}
	if err := json.Unmarshal([]byte(data), &reviews); err != nil {
		return err
	}

	*r = reviews
	return nil
}

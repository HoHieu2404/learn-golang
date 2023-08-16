package models

type Data struct {
	Rates map[string]map[string]float32 `json:"rates"`
}

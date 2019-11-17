package model

type ModelItem struct {
	ID    string `json:"id"`
	Oem   string `json:"oem"`
	Model string `json:"model"`
	Total int    `json:"total"`
}

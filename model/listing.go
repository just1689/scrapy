package model

type ListingItem struct {
	ID    string `json:"id"`
	Oem   string `json:"oem"`
	Model string `json:"model"`
	Url   string `json:"url"`
}

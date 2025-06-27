package model

type Recipe struct {
	Base
	Name       string `json:"name"`
	Cookbook   string `json:"cookbook"`
	Pagenumber int    `json:"pagenumber"`
	Image      string `json:"image"`
	UserID     string `json:"userID"`
	User       User
}

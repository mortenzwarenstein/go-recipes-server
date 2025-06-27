package dto

type Response struct {
	Data    interface{} `json:"data"`
	Error   string      `json:"error"`
	Message string      `json:"message"`
}

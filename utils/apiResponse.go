package utils

type APIResponse[T any] struct {
	Error  string `json:"error"`
	Status int    `json:"status"`
	Data   T      `json:"data"`
}

package types

type LCDError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

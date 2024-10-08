package types

import (
	"time"
)

type BlockResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Block   struct {
		Header struct {
			Time   time.Time `json:"time"`
			Height int64     `json:"height,string"`
		} `json:"header"`
	} `json:"block"`
}

type LCDError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

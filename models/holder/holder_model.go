package holder

import (
	"encoding/json"
)

type Holder struct {
	Address string `json:"address"`
	Balance json.Number `json:"balance"`
	Share float32 `json:"share"`
}
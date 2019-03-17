package api

import "github.com/rijdendetreinen/gotrain/stores"

// Statistics includes counters and the inventory
type Statistics struct {
	Counters  stores.Counters `json:"counters"`
	Inventory int             `json:"inventory"`
}

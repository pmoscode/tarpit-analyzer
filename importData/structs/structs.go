package structs

import time2 "time"

type ImportItem struct {
	Begin    time2.Time
	End      time2.Time
	Duration float32
	Ip       string
	Success  bool
}

package structs

import time2 "time"

type ImportItem struct {
	Begin    time2.Time `gorm:"index"`
	End      time2.Time `gorm:"index"`
	Duration float32
	Ip       string
	Success  bool
}

package structs

import time2 "time"

type ImportItem struct {
	Begin    time2.Time `gorm:"primaryKey"`
	End      time2.Time
	Duration float32
	Ip       string `gorm:"primaryKey"`
	Success  bool
}

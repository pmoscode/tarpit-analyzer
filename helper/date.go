package helper

import (
	log "github.com/sirupsen/logrus"
	time2 "time"
)

func GetDate(dateString string) *time2.Time {
	if dateString != "unset" {
		var date time2.Time
		var err error

		date, err = time2.Parse("2006-01-02", dateString)
		if err != nil {
			log.Panicln("Parameter 'StartDate' is not a valid date!")
		}

		return &date
	}

	return nil
}

func IsBefore(now time2.Time, then *time2.Time) bool {
	if then == nil {
		return true
	}

	return now.AddDate(0, 0, 1).Before(*then)
}

func IsAfter(now time2.Time, then *time2.Time) bool {
	if then == nil {
		return true
	}

	return now.AddDate(0, 0, -1).After(*then)
}

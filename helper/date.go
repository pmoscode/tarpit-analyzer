package helper

import (
	log "github.com/sirupsen/logrus"
	time2 "time"
)

func GetDate(dateString string) *time2.Time {
	if dateString != "unset" {
		var date *time2.Time
		var err error

		*date, err = time2.Parse("2006-01-02", dateString)
		if err != nil {
			log.Panicln("Parameter 'StartDate' is not a valid date!")
		}

		return date
	}

	return nil
}

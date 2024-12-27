package cron

import (
	"time"

	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

var parser = cron.NewParser(
	cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow,
)

func Next(cronExpr string, start time.Time) time.Time {
	s, err := parser.Parse(cronExpr)
	if err != nil {
		log.Fatal(err)
	}
	return s.Next(start)
}

func IsValid(cronExpr string) bool {
	_, err := parser.Parse(cronExpr)
	return err == nil
}

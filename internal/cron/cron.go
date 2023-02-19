package cron

import (
	"time"

	"github.com/adhocore/gronx"
	log "github.com/sirupsen/logrus"
)

func NextTickAfter(cronExpr string, start time.Time, includeStart bool) time.Time {
	t, err := gronx.NextTickAfter(cronExpr, start, includeStart)
	if err != nil {
		log.Fatal(err)
	}
	return t
}

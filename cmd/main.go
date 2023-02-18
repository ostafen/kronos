package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/ostafen/kronos/internal/api"
	"github.com/ostafen/kronos/internal/config"
	"github.com/ostafen/kronos/internal/db"
	"github.com/ostafen/kronos/internal/model"
	"github.com/ostafen/kronos/internal/notification"
	"github.com/ostafen/kronos/internal/service"
	log "github.com/sirupsen/logrus"
)

func printLogo() {
	fmt.Println(`| | ___ __ ___  _ __   ___  ___`)
	fmt.Println(`| |/ / '__/ _ \| '_ \ / _ \/ __|`)
	fmt.Println(`|   <| | | (_) | | | | (_) \__ \`)
	fmt.Println(`|_|\_\_|  \___/|_| |_|\___/|___/`)
	fmt.Println()
}

func main() {
	printLogo()

	if len(os.Args) < 2 {
		log.Fatal("no config provided")
		os.Exit(1)
	}

	conf, err := config.Parse(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	setupLogging(conf.Service.Logging)

	dbConn, err := db.Open(conf.Store)
	if err != nil {
		log.Fatal(err)
	}

	m := db.NewMigrator(dbConn)
	if err := m.Migrate(); err != nil {
		log.Fatal(err)
	}

	schedRepo := db.GetScheduleRepo(dbConn)
	schedSvc := service.NewScheduleService(dbConn, schedRepo, service.NewNotificationService())
	alertSvc := service.NewAlertService(conf.Service.Alert.Email)

	trigger := notification.NewScheduleTrigger(schedSvc)

	awakeTrigger := func(_ *model.Schedule) {
		trigger.WakeUp()
	}

	schedSvc.OnScheduleRegistered(awakeTrigger)
	schedSvc.OnScheduleResumed(awakeTrigger)

	schedSvc.OnSchedulePaused(func(sched *model.Schedule) {
		alertSvc.Send(sched)
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	trigger.Start(ctx)

	configureRouter(schedSvc)

	http.ListenAndServe(fmt.Sprintf(":%d", conf.Service.Port), nil)
}

func setupLogging(config config.Log) {
	log.SetReportCaller(true)
	log.SetLevel(getLogLevel(config.Level))
	log.SetFormatter(getFormatter(config.Format))
}

func getLogLevel(level string) log.Level {
	switch strings.ToUpper(level) {
	case "TRACE":
		return log.TraceLevel
	case "DEBUG":
		return log.DebugLevel
	case "INFO":
		return log.InfoLevel
	case "FATAL":
		return log.FatalLevel
	case "PANIC":
		return log.PanicLevel
	}
	return log.InfoLevel
}

func getFormatter(format string) log.Formatter {
	switch format {
	case "JSON":
		return &log.JSONFormatter{}
	case "TEXT":
		return &log.TextFormatter{}
	}
	return &log.JSONFormatter{}
}

func configureRouter(svc service.ScheduleService) {
	r := mux.NewRouter()

	api := api.NewScheduleApi(svc)

	r.HandleFunc("/schedules", api.ListSchedules).Methods("GET")
	r.HandleFunc("/schedules/{id}", api.GetSchedule).Methods("GET")
	r.HandleFunc("/schedules/{id}", api.DeleteSchedule).Methods("DELETE")

	r.HandleFunc("/schedules", api.RegisterSchedule).Methods("POST")
	r.HandleFunc("/schedules/{id}/pause", api.PauseSchedule).Methods("POST")
	r.HandleFunc("/schedules/{id}/resume", api.ResumeSchedule).Methods("POST")
	r.HandleFunc("/schedules/{id}/trigger", api.TriggerSchedule).Methods("POST")

	http.Handle("/", r)
}

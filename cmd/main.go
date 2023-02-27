package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/ostafen/kronos/internal/api"
	"github.com/ostafen/kronos/internal/config"
	"github.com/ostafen/kronos/internal/metrics"
	"github.com/ostafen/kronos/internal/model"
	"github.com/ostafen/kronos/internal/sched"
	"github.com/ostafen/kronos/internal/service"
	"github.com/ostafen/kronos/internal/store"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	log "github.com/sirupsen/logrus"
)

var (
	version   string
	commit    string = "none"
	buildTime string = time.Now().Format(time.UnixDate)
)

func getVersion() string {
	if version == "" {
		return pseudoVersion()
	}
	return version
}

func pseudoVersion() string {
	return fmt.Sprintf("v0.0.0-%s-%s", time.Now().Format("20060102150405"), commit)
}

func printLogo() {
	fmt.Println("| | ___ __ ___  _ __   ___  ___")
	fmt.Println("| |/ / '__/ _ \\| '_ \\ / _ \\/ __|")
	fmt.Println("|   <| | | (_) | | | | (_) \\__ \\")
	fmt.Println("|_|\\_\\_|  \\___/|_| |_|\\___/|___/")
	fmt.Printf("Version: %s\n", getVersion())
	fmt.Printf("Commit: %s\n", commit)
	fmt.Printf("Build.Time: %s\n\n", buildTime)
}

func main() {
	printLogo()

	conf, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	setupLogging(conf.Logging)

	store, err := store.New(conf.Store.Path)
	if err != nil {
		log.Fatal(err)
	}

	svc := service.NewScheduleService(store, service.NewNotificationService())
	manager := sched.NewScheduleManager(svc.OnTick)

	err = store.Iterate(func(sched *model.Schedule) error {
		if sched.IsActive() {
			log.Info("scheduling %s at %s", sched.ID, sched.NextTickAt())

			manager.Schedule(sched.ID, sched.NextTickAt())
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	awake := func(_ *model.Schedule) {
		manager.WakeUp()
	}

	submitSchedule := func(s *model.Schedule) {
		manager.Schedule(s.ID, s.FirstTick())
	}

	svc.OnScheduleRegistered(awake, submitSchedule)

	svc.OnSchedulePaused(func(s *model.Schedule) {
		manager.Remove(s.ID)
	})

	svc.OnScheduleResumed(awake, func(s *model.Schedule) {
		manager.Schedule(s.ID, s.NextTickAt())
		metrics.ResetScheduleFailures(s.ID)
	})

	svc.OnScheduleNotified(func(s *model.Schedule, code int) {
		if code < 200 || code >= 300 {
			metrics.IncScheduleFailures(s.ID)
		} else {
			metrics.ResetScheduleFailures(s.ID)
		}
		metrics.IncWebhookRequests(s.URL, code)
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	manager.Start(ctx)

	configureRouter(svc)

	http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), nil)
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

	scheduleApi := api.NewScheduleApi(svc)

	r.Handle("/metrics", promhttp.Handler()).Methods("GET")

	r.HandleFunc("/schedules", scheduleApi.ListSchedules).Methods("GET")
	r.HandleFunc("/schedules/{id}", scheduleApi.GetSchedule).Methods("GET")
	r.HandleFunc("/schedules/{id}", scheduleApi.DeleteSchedule).Methods("DELETE")

	r.HandleFunc("/schedules", scheduleApi.RegisterSchedule).Methods("POST")
	r.HandleFunc("/schedules/{id}/pause", scheduleApi.PauseSchedule).Methods("POST")
	r.HandleFunc("/schedules/{id}/resume", scheduleApi.ResumeSchedule).Methods("POST")
	r.HandleFunc("/schedules/{id}/trigger", scheduleApi.TriggerSchedule).Methods("POST")

	http.Handle("/", r)
}

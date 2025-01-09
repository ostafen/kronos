package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/ostafen/kronos/internal/api"
	"github.com/ostafen/kronos/internal/config"
	"github.com/ostafen/kronos/internal/service"
	"github.com/ostafen/kronos/internal/store"
	statichttp "github.com/ostafen/kronos/webbuild"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/rs/cors"
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
	fmt.Println()
	fmt.Printf("Version: %s\n", getVersion())
	fmt.Printf("Commit: %s\n", commit)
	fmt.Printf("Built At: %s\n\n", buildTime)
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

	svc := service.NewScheduleService(
		store,
		service.NewNotificationService(),
	)
	defer svc.Stop()

	configureRouter(svc)

	http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), nil)

	// TODO: soft shutdown
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
	fs := http.FileServer(http.FS(statichttp.Static))
	r.PathPrefix("/web").Handler(http.StripPrefix("/", fs))

	handler := api.NewScheduleApiHandler(svc)

	r.Handle("/metrics", promhttp.Handler()).Methods("GET")

	// TODO: health endpoint

	r.HandleFunc("/api/v1/schedules", handler.ListSchedules).Methods("GET")
	r.HandleFunc("/api/v1/schedules/{id}", handler.GetSchedule).Methods("GET")
	r.HandleFunc("/api/v1/schedules/{id}", handler.DeleteSchedule).Methods("DELETE")

	r.HandleFunc("/api/v1/schedules", handler.RegisterSchedule).Methods("POST")
	r.HandleFunc("/api/v1/schedules/{id}/pause", handler.PauseSchedule).Methods("POST")
	r.HandleFunc("/api/v1/schedules/{id}/resume", handler.ResumeSchedule).Methods("POST")
	r.HandleFunc("/api/v1/schedules/{id}/trigger", handler.TriggerSchedule).Methods("POST")

	r.HandleFunc("/api/v1/history", handler.GetHistory).Methods("GET")
	r.HandleFunc("/api/v1/history/{id}", handler.GetCronHistory).Methods("GET")

	http.Handle("/", withCors(r, cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	}))
}

func withCors(handler http.Handler, opts cors.Options) http.Handler {
	c := cors.New(opts)
	return c.Handler(handler)
}

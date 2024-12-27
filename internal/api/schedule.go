package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/ostafen/kronos/internal/model"
	"github.com/ostafen/kronos/internal/service"
)

type ScheduleApiHandler struct {
	svc service.ScheduleService
}

func NewScheduleApiHandler(svc service.ScheduleService) *ScheduleApiHandler {
	return &ScheduleApiHandler{
		svc: svc,
	}
}

func (api *ScheduleApiHandler) RegisterSchedule(w http.ResponseWriter, r *http.Request) {
	var input model.ScheduleRegisterInput

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	v := validator.New()
	if err := v.Struct(input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sched, err := api.svc.RegisterSchedule(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, sched)
}

func (api *ScheduleApiHandler) PauseSchedule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sched, err := api.svc.PauseSchedule(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, sched)
}

func (api *ScheduleApiHandler) ResumeSchedule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sched, err := api.svc.ResumeSchedule(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, sched)
}

func (api *ScheduleApiHandler) TriggerSchedule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sched, err := api.svc.TriggerSchedule(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, sched)
}

func (api *ScheduleApiHandler) GetSchedule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sched, err := api.svc.GetSchedule(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if sched == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	writeJSON(w, sched)
}

func (api *ScheduleApiHandler) DeleteSchedule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = api.svc.DeleteSchedule(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted) // TODO: write json error
}

func (api *ScheduleApiHandler) ListSchedules(w http.ResponseWriter, r *http.Request) {
	schedules := make([]*model.CronSchedule, 0)
	err := api.svc.IterSchedules(func(s *model.CronSchedule) error {
		schedules = append(schedules, s)
		return nil
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, schedules)
}

func (api *ScheduleApiHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	statuses, err := api.svc.GetHistory()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, statuses)
}

func (api *ScheduleApiHandler) GetCronHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	statuses, err := api.svc.GetCronHistory(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, statuses)
}

func writeJSON(w http.ResponseWriter, body any) {
	w.Header().Add("content-type", "application/json")

	data, err := json.Marshal(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/adhocore/gronx"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/ostafen/kronos/internal/service"
)

type ScheduleApi struct {
	svc service.ScheduleService
}

func NewScheduleApi(svc service.ScheduleService) *ScheduleApi {
	return &ScheduleApi{
		svc: svc,
	}
}

func (api *ScheduleApi) RegisterSchedule(w http.ResponseWriter, r *http.Request) {
	var schedReq ScheduleCreateRequest

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&schedReq); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	v := validator.New()
	if err := v.Struct(schedReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	gron := gronx.New()
	if !gron.IsValid(schedReq.CronExpr) {
		http.Error(w, fmt.Sprintf("invalid chron expr %s", schedReq.CronExpr), http.StatusBadRequest)
		return
	}

	sched := schedReq.ToSched()
	err := api.svc.RegisterSchedule(sched)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, sched)
}

func (api *ScheduleApi) PauseSchedule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	sched, err := api.svc.PauseSchedule(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, sched)
}

func (api *ScheduleApi) ResumeSchedule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	sched, err := api.svc.ResumeSchedule(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, sched)
}

func (api *ScheduleApi) TriggerSchedule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	sched, err := api.svc.TriggerSchedule(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, sched)
}

func (api *ScheduleApi) GetSchedule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	sched, err := api.svc.GetSchedule(vars["id"])
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

func (api *ScheduleApi) DeleteSchedule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	err := api.svc.DeleteSchedule(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted) // TODO: write json error
}

func (api *ScheduleApi) ListSchedules(w http.ResponseWriter, r *http.Request) {
	schedules, err := api.svc.ListSchedules(-1, -1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, schedules)
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

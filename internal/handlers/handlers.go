package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/deeprajsshetty/alarm-service/internal/models"
	"github.com/deeprajsshetty/alarm-service/internal/services"
)

// AlarmHandler handles HTTP requests for alarm-related operations.
type AlarmHandler struct {
	service *services.AlarmService
}

// NewAlarmHandler initializes and returns a new AlarmHandler instance.
func NewAlarmHandler(service *services.AlarmService) *AlarmHandler {
	return &AlarmHandler{service: service}
}

// CreateAlarm handles the creation of new alarms.
func (h *AlarmHandler) CreateAlarm(w http.ResponseWriter, r *http.Request) {
	var alarm models.Alarm
	if err := json.NewDecoder(r.Body).Decode(&alarm); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	createdAlarm, err := h.service.CreateAlarm(alarm)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusCreated, createdAlarm)
}

// BulkCreateAlarms handles bulk creation of multiple alarms.
func (h *AlarmHandler) BulkCreateAlarms(w http.ResponseWriter, r *http.Request) {
	var alarms []models.Alarm
	if err := json.NewDecoder(r.Body).Decode(&alarms); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request payload for bulk creation")
		return
	}

	createdAlarms, err := h.service.BulkCreateAlarms(alarms)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusCreated, createdAlarms)
}

// GetAllAlarms retrieves and returns all alarms.
func (h *AlarmHandler) GetAllAlarms(w http.ResponseWriter, r *http.Request) {
	alarms := h.service.GetAllAlarms()
	h.respondWithJSON(w, http.StatusOK, alarms)
}

// GetAlarmByID retrieves a specific alarm by its ID.
func (h *AlarmHandler) GetAlarmByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	alarm, err := h.service.GetAlarmByID(id)
	if err != nil {
		h.respondWithError(w, http.StatusNotFound, "Alarm not found")
		return
	}

	h.respondWithJSON(w, http.StatusOK, alarm)
}

// UpdateAlarmState updates an existing alarm's state by ID.
func (h *AlarmHandler) UpdateAlarmState(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	var request struct {
		State models.AlarmState `json:"state"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	alarm, err := h.service.UpdateAlarmState(id, request.State)
	if err != nil {
		h.respondWithError(w, http.StatusNotFound, "Alarm not found")
		return
	}

	h.respondWithJSON(w, http.StatusOK, alarm)
}

// DeleteAlarm deletes an alarm by ID and responds with a proper status and message.
func (h *AlarmHandler) DeleteAlarm(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		h.respondWithError(w, http.StatusBadRequest, "Alarm ID is required")
		return
	}

	msg, err := h.service.DeleteAlarm(id)
	if err != nil {
		h.respondWithError(w, http.StatusNotFound, "Alarm not found")
		return
	}

	response := map[string]string{"message": msg}
	h.respondWithJSON(w, http.StatusOK, response)
}

// respondWithJSON sends a JSON response with the given status code and payload.
func (h *AlarmHandler) respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}

// respondWithError sends an error response with the given status code and message.
func (h *AlarmHandler) respondWithError(w http.ResponseWriter, statusCode int, message string) {
	h.respondWithJSON(w, statusCode, map[string]string{"error": message})
}
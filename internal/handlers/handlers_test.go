package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/deeprajsshetty/alarm-service/internal/models"
	"github.com/deeprajsshetty/alarm-service/internal/services"
	"github.com/stretchr/testify/assert"
)

// TestCreateAlarm_Success tests successful alarm creation via HTTP handler.
func TestCreateAlarm_Success(t *testing.T) {
	service := services.NewAlarmService()
	handler := NewAlarmHandler(service)

	alarmPayload := `{"name": "Server Overload", "state": "Triggered"}`
	req := httptest.NewRequest(http.MethodPost, "/alarm", bytes.NewBuffer([]byte(alarmPayload)))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	handler.CreateAlarm(recorder, req)

	assert.Equal(t, http.StatusCreated, recorder.Code, "Expected HTTP 201 Created")
	var response models.Alarm
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Server Overload", response.Name, "Expected correct alarm name")
}

// TestCreateAlarm_InvalidPayload tests creating an alarm with invalid payload.
func TestCreateAlarm_InvalidPayload(t *testing.T) {
	service := services.NewAlarmService()
	handler := NewAlarmHandler(service)

	invalidPayload := `{"state": "Triggered"}`
	req := httptest.NewRequest(http.MethodPost, "/alarm", bytes.NewBuffer([]byte(invalidPayload)))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	handler.CreateAlarm(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code, "Expected HTTP 400 Bad Request")
}

// TestGetAllAlarms tests fetching all alarms.
func TestGetAllAlarms(t *testing.T) {
	service := services.NewAlarmService()
	handler := NewAlarmHandler(service)

	req := httptest.NewRequest(http.MethodGet, "/alarms", nil)
	recorder := httptest.NewRecorder()

	handler.GetAllAlarms(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code, "Expected HTTP 200 OK")
}

// TestGetAlarmByID_Success tests fetching an alarm by valid ID.
func TestGetAlarmByID_Success(t *testing.T) {
	service := services.NewAlarmService()
	alarm, _ := service.CreateAlarm(models.Alarm{Name: "Memory Alert", State: models.Triggered})
	handler := NewAlarmHandler(service)

	req := httptest.NewRequest(http.MethodGet, "/alarm?id="+alarm.ID, nil)
	recorder := httptest.NewRecorder()

	handler.GetAlarmByID(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code, "Expected HTTP 200 OK")
}

// TestGetAlarmByID_NotFound tests fetching an alarm with an invalid ID.
func TestGetAlarmByID_NotFound(t *testing.T) {
	service := services.NewAlarmService()
	handler := NewAlarmHandler(service)

	req := httptest.NewRequest(http.MethodGet, "/alarm?id=invalidID", nil)
	recorder := httptest.NewRecorder()

	handler.GetAlarmByID(recorder, req)

	assert.Equal(t, http.StatusNotFound, recorder.Code, "Expected HTTP 404 Not Found")
}

// TestUpdateAlarmState_Success tests successfully updating an alarm's state.
func TestUpdateAlarmState_Success(t *testing.T) {
	service := services.NewAlarmService()
	alarm, _ := service.CreateAlarm(models.Alarm{Name: "Disk Space Alert", State: models.Triggered})
	handler := NewAlarmHandler(service)

	updatePayload := `{"state":"Cleared"}`
	req := httptest.NewRequest(http.MethodPut, "/alarm?id="+alarm.ID, bytes.NewBuffer([]byte(updatePayload)))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	handler.UpdateAlarmState(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code, "Expected HTTP 200 OK")
}

// TestUpdateAlarmState_NotFound tests updating an alarm with an invalid ID.
func TestUpdateAlarmState_NotFound(t *testing.T) {
	service := services.NewAlarmService()
	handler := NewAlarmHandler(service)

	req := httptest.NewRequest(http.MethodPut, "/alarm?id=invalidID", bytes.NewBuffer([]byte(`{"state":"Cleared"}`)))
	recorder := httptest.NewRecorder()

	handler.UpdateAlarmState(recorder, req)

	assert.Equal(t, http.StatusNotFound, recorder.Code, "Expected HTTP 404 Not Found")
}

// TestDeleteAlarm_Success tests successfully deleting an alarm.
func TestDeleteAlarm_Success(t *testing.T) {
	service := services.NewAlarmService()
	alarm, _ := service.CreateAlarm(models.Alarm{Name: "CPU Overload", State: models.Triggered})
	handler := NewAlarmHandler(service)

	req := httptest.NewRequest(http.MethodDelete, "/alarm?id="+alarm.ID, nil)
	recorder := httptest.NewRecorder()

	handler.DeleteAlarm(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code, "Expected HTTP 200 OK")
}

// TestDeleteAlarm_NotFound tests deleting an alarm with an invalid ID.
func TestDeleteAlarm_NotFound(t *testing.T) {
	service := services.NewAlarmService()
	handler := NewAlarmHandler(service)

	req := httptest.NewRequest(http.MethodDelete, "/alarm?id=invalidID", nil)
	recorder := httptest.NewRecorder()

	handler.DeleteAlarm(recorder, req)

	assert.Equal(t, http.StatusNotFound, recorder.Code, "Expected HTTP 404 Not Found")
}

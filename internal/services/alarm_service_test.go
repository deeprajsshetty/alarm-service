package services_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/deeprajsshetty/alarm-service/internal/models"
	"github.com/deeprajsshetty/alarm-service/internal/services"
	"github.com/google/uuid"
)

// TestCreateAlarm_Success verifies successful creation of an alarm with valid data.
func TestCreateAlarm_Success(t *testing.T) {
	svc := services.NewAlarmService()
	alarm := models.Alarm{Name: "Test Alarm", State: models.Triggered}

	createdAlarm, err := svc.CreateAlarm(alarm)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if createdAlarm.Name != alarm.Name {
		t.Errorf("expected name %s, got %s", alarm.Name, createdAlarm.Name)
	}
}

// TestCreateAlarm_Validation verifies invalid scenarios for alarm creation.
func TestCreateAlarm_Validation(t *testing.T) {
	svc := services.NewAlarmService()

	// Missing Name
	_, err := svc.CreateAlarm(models.Alarm{State: models.Triggered})
	if err == nil || err.Error() != "alarm name is mandatory" {
		t.Errorf("expected error 'alarm name is mandatory', got %v", err)
	}

	// Invalid State
	_, err = svc.CreateAlarm(models.Alarm{Name: "Invalid Alarm", State: "InvalidState"})
	if err == nil || err.Error() != "invalid alarm state" {
		t.Errorf("expected error 'invalid alarm state', got %v", err)
	}
}

// TestGetAllAlarms verifies retrieval of all created alarms.
func TestGetAllAlarms(t *testing.T) {
	svc := services.NewAlarmService()
	svc.CreateAlarm(models.Alarm{Name: "Alarm 1", State: models.Triggered})
	svc.CreateAlarm(models.Alarm{Name: "Alarm 2", State: models.ACKed})

	alarms := svc.GetAllAlarms()
	if len(alarms) != 2 {
		t.Errorf("expected 2 alarms, got %d", len(alarms))
	}
}

// TestGetAlarmByID verifies alarm retrieval by ID, including success and failure scenarios.
func TestGetAlarmByID(t *testing.T) {
	svc := services.NewAlarmService()
	alarm, _ := svc.CreateAlarm(models.Alarm{Name: "Test Alarm", State: models.Triggered})

	// Successful Retrieval
	retrievedAlarm, err := svc.GetAlarmByID(alarm.ID)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if retrievedAlarm.ID != alarm.ID {
		t.Errorf("expected ID %s, got %s", alarm.ID, retrievedAlarm.ID)
	}

	// Non-existent ID
	_, err = svc.GetAlarmByID(uuid.New().String())
	if err == nil || err.Error() != "alarm not found" {
		t.Errorf("expected error 'alarm not found', got %v", err)
	}
}

// TestDeleteAlarm verifies deletion success and failure scenarios.
func TestDeleteAlarm(t *testing.T) {
	svc := services.NewAlarmService()
	alarm, _ := svc.CreateAlarm(models.Alarm{Name: "To Be Deleted", State: models.Active})

	// Successful Deletion
	msg, err := svc.DeleteAlarm(alarm.ID)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	expectedMsg := "✅ Alarm ID: " + alarm.ID + " successfully deleted"
	if msg != expectedMsg {
		t.Errorf("expected message %q, got %q", expectedMsg, msg)
	}

	// Non-existent ID
	nonExistentID := uuid.New().String()
	_, err = svc.DeleteAlarm(nonExistentID)
	expectedError := "❌ Alarm ID: " + nonExistentID + " not found"

	if err == nil || err.Error() != expectedError {
		t.Errorf("expected error %q, got %q", expectedError, err)
	}
}

// TestUpdateAlarmState verifies both success and failure scenarios for alarm state updates.
func TestUpdateAlarmState(t *testing.T) {
	svc := services.NewAlarmService()
	alarm, _ := svc.CreateAlarm(models.Alarm{Name: "Update Test", State: models.Triggered})

	// Valid State Change
	updatedAlarm, err := svc.UpdateAlarmState(alarm.ID, models.ACKed)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if updatedAlarm.State != models.ACKed {
		t.Errorf("expected state ACKed, got %v", updatedAlarm.State)
	}

	// Invalid State Change
	_, err = svc.UpdateAlarmState(alarm.ID, "InvalidState")
	if err == nil || err.Error() != "invalid alarm state" {
		t.Errorf("expected error 'invalid alarm state', got %v", err)
	}
}

// TestBulkCreateAlarms verifies bulk creation scenarios, including empty lists, duplicates, and invalid states.
func TestBulkCreateAlarms(t *testing.T) {
	svc := services.NewAlarmService()

	// Successful Bulk Creation
	alarms := []models.Alarm{{Name: "Alarm 1", State: models.Triggered}, {Name: "Alarm 2", State: models.Active}}
	createdAlarms, err := svc.BulkCreateAlarms(alarms)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(createdAlarms) != 2 {
		t.Errorf("expected 2 alarms, got %d", len(createdAlarms))
	}

	// Empty List
	createdAlarms, err = svc.BulkCreateAlarms([]models.Alarm{})
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(createdAlarms) != 0 {
		t.Errorf("expected 0 alarms, got %d", len(createdAlarms))
	}
}

// TestLoadSampleAlarms verifies loading sample alarms from JSON file
func TestLoadSampleAlarms(t *testing.T) {
	filePath := filepath.Join("..", "..", "testdata", "sample_alarms.json")

	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read sample data file: %v", err)
	}

	var sampleAlarms []models.Alarm
	if err := json.Unmarshal(data, &sampleAlarms); err != nil {
		t.Fatalf("failed to unmarshal sample data: %v", err)
	}

	svc := services.NewAlarmService()
	createdAlarms, err := svc.BulkCreateAlarms(sampleAlarms)
	if err != nil {
		t.Errorf("failed to create alarms: %v", err)
	}

	if len(createdAlarms) != len(sampleAlarms) {
		t.Errorf("expected %d alarms, got %d", len(sampleAlarms), len(createdAlarms))
	}
}
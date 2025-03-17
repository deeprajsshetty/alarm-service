package services

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/deeprajsshetty/alarm-service/internal/models"
	"github.com/google/uuid"
)

// AlarmService manages alarm operations with thread safety and notification support.
type AlarmService struct {
	alarms               map[string]models.Alarm
	lock                 sync.RWMutex
	notifyChan           chan models.Alarm
	notificationSchedule map[string]time.Time
}

// NewAlarmService initializes and returns a new AlarmService instance.
func NewAlarmService() *AlarmService {
	svc := &AlarmService{
		alarms:               make(map[string]models.Alarm),
		notifyChan:           make(chan models.Alarm, 100),
		notificationSchedule: make(map[string]time.Time),
	}

	go svc.startNotificationHandler()
	go svc.startScheduler()
	return svc
}

// NotificationInterval defines intervals for sending notifications based on alarm state.
type NotificationInterval struct {
	Interval time.Duration
}

// stateNotificationIntervals manages notification intervals per alarm state.
var stateNotificationIntervals = map[models.AlarmState]NotificationInterval{
	models.Triggered: {Interval: 2 * time.Hour},
	models.ACKed:     {Interval: 24 * time.Hour},
}

// startNotificationHandler continuously processes alarm notifications.
func (s *AlarmService) startNotificationHandler() {
	for alarm := range s.notifyChan {
		s.processNotification(alarm)
	}
}

// startScheduler continuously checks scheduled alarms and triggers them automatically.
func (s *AlarmService) startScheduler() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		s.checkAndTriggerNotifications()
	}
}

// checkAndTriggerNotifications identifies alarms that are due for notification.
func (s *AlarmService) checkAndTriggerNotifications() {
	s.lock.Lock()
	defer s.lock.Unlock()

	now := time.Now()
	for id, nextNotifyTime := range s.notificationSchedule {
		if now.After(nextNotifyTime) {
			if alarm, found := s.alarms[id]; found {
	/*	
				// Commented this code as it is not part of requirement.
				// This logic about, in case alarm manually not acknowledged 
				// also by default acknoledged in 24 Hours	
				if alarm.State == models.Triggered {
					createdAt, err := s.getCreatedAtTime(alarm)
					if err == nil && now.Sub(createdAt) >= stateNotificationIntervals[models.ACKed].Interval {
						alarm.State = models.ACKed
						alarm.ACKedAt = now.Format(time.RFC3339)
					}
				}
	*/
				intervalData, exists := stateNotificationIntervals[alarm.State]
				if exists {
					s.notificationSchedule[alarm.ID] = now.Add(intervalData.Interval)
				}

				s.notifyChan <- alarm 
			}
		}
	}
}

// processNotification handles sending notifications with appropriate intervals.
func (s *AlarmService) processNotification(alarm models.Alarm) {
	s.lock.Lock()
	defer s.lock.Unlock()

	intervalData, exists := stateNotificationIntervals[alarm.State]
	if !exists {
		return
	}

	fmt.Printf("🔔 Notification for Alarm ID: %s - State: %s\n", alarm.ID, alarm.State)
	s.notificationSchedule[alarm.ID] = time.Now().Add(intervalData.Interval)
}

// CreateAlarm creates a new alarm with default values and triggers notification if applicable.
func (s *AlarmService) CreateAlarm(alarm models.Alarm) (models.Alarm, error) {
	if err := s.validateAlarm(alarm); err != nil {
		return models.Alarm{}, err
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	s.initializeAlarm(&alarm)
	s.alarms[alarm.ID] = alarm
	s.notifyChan <- alarm // Notify immediately when created in 'Triggered' state

	return alarm, nil
}

// BulkCreateAlarms handles bulk alarm creation with concurrency safety.
func (s *AlarmService) BulkCreateAlarms(alarms []models.Alarm) ([]models.Alarm, error) {
	var createdAlarms []models.Alarm
	var errorList []string

	s.lock.Lock()
	defer s.lock.Unlock()

	for _, alarm := range alarms {
		if err := s.validateAlarm(alarm); err != nil {
			errorList = append(errorList, fmt.Sprintf("Alarm %s: %v", alarm.Name, err))
			continue
		}

		s.initializeAlarm(&alarm)
		s.alarms[alarm.ID] = alarm
		s.notifyChan <- alarm
		createdAlarms = append(createdAlarms, alarm)
	}

	if len(errorList) > 0 {
		return createdAlarms, fmt.Errorf("failed to create some alarms: %v", errorList)
	}

	return createdAlarms, nil
}

// initializeAlarm sets default values for a new alarm.
func (s *AlarmService) initializeAlarm(alarm *models.Alarm) {
	alarm.ID = uuid.New().String()
	alarm.CreatedAt = time.Now().Format(time.RFC3339)
	alarm.State = models.Triggered
	s.notificationSchedule[alarm.ID] = time.Now().Add(stateNotificationIntervals[models.Triggered].Interval)
}

// GetAllAlarms retrieves all stored alarms in memory.
func (s *AlarmService) GetAllAlarms() []models.Alarm {
	s.lock.RLock()
	defer s.lock.RUnlock()

	alarms := make([]models.Alarm, 0, len(s.alarms))
	for _, alarm := range s.alarms {
		alarms = append(alarms, alarm)
	}
	return alarms
}

// GetAlarmByID retrieves an alarm by its unique ID.
func (s *AlarmService) GetAlarmByID(id string) (models.Alarm, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if alarm, found := s.alarms[id]; found {
		return alarm, nil
	}
	return models.Alarm{}, errors.New("alarm not found")
}

// UpdateAlarmState updates the state of an alarm and triggers a notification if necessary.
func (s *AlarmService) UpdateAlarmState(id string, state models.AlarmState) (models.Alarm, error) {
	if !state.IsValid() {
		return models.Alarm{}, errors.New("invalid alarm state")
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	if alarm, found := s.alarms[id]; found {
		alarm.State = state
		alarm.UpdatedAt = time.Now().Format(time.RFC3339)

		if state == models.ACKed {
			alarm.ACKedAt = time.Now().Format(time.RFC3339)
		}

		s.alarms[id] = alarm
		s.notifyChan <- alarm
		return alarm, nil
	}

	return models.Alarm{}, errors.New("alarm not found")
}

// DeleteAlarm removes an alarm from the in-memory store by ID.
func (s *AlarmService) DeleteAlarm(id string) (string, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, found := s.alarms[id]; found {
		delete(s.alarms, id)
		delete(s.notificationSchedule, id)

		logMessage := fmt.Sprintf("✅ Alarm ID: %s successfully deleted", id)
		return logMessage, nil
	}
	errMessage := fmt.Sprintf("❌ Alarm ID: %s not found", id)
	return "", errors.New(errMessage)
}

// validateAlarm verifies alarm data to ensure valid state and non-empty name.
func (s *AlarmService) validateAlarm(alarm models.Alarm) error {
	if alarm.Name == "" {
		return errors.New("alarm name is mandatory")
	}
	if !alarm.State.IsValid() {
		return errors.New("invalid alarm state")
	}
	return nil
}

// getCreatedAtTime parses the CreatedAt field as time.Time
func (s *AlarmService) getCreatedAtTime(alarm models.Alarm) (time.Time, error) {
	return time.Parse(time.RFC3339, alarm.CreatedAt)
}
package usecase

import (
	"context"
	"encoding/json"
	"log"
	"notification_service/src/internal/core/notification"
)

type NotificationUseCase struct{}

func NewNotificationUseCase() *NotificationUseCase {
	return &NotificationUseCase{}
}

func (uc *NotificationUseCase) HandleNotification(ctx context.Context, data string) error {
	var note notification.TaskNotification
	if err := json.Unmarshal([]byte(data), &note); err != nil {
		return err
	}

	log.Printf("[NOTIFICATION RECEIVED] TaskID: %d, Event: %s, AssignedTo: %d, Message: %s",
		note.TaskID, note.Event, note.AssignedTo, note.Message)

	return nil
}

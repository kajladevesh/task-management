package session

type TaskNotification struct {
	Event      string `json:"event"` // "assigned" or "updated"
	TaskID     int64  `json:"task_id"`
	AssignedTo int64  `json:"assigned_to"`
	Message    string `json:"message"`
}

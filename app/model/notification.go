package model

import "time"

const (
	// QueueDefault adalah nama queue Asynq yang dipakai worker saat ini.
	QueueDefault = "default"

	// TaskSendNotification adalah jenis task pertama yang akan ditangani worker.
	TaskSendNotification = "notification:send"
)

type Notification struct {
	ID        int64      `json:"id"`
	UserID    int        `json:"user_id"`
	Title     string     `json:"title"`
	Body      string     `json:"body"`
	ReadAt    *time.Time `json:"read_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

package model

import "time"

type TicketStatus string

const (
	TicketStatusOpen       TicketStatus = "OPEN"
	TicketStatusAssigned   TicketStatus = "ASSIGNED"
	TicketStatusInProgress TicketStatus = "IN_PROGRESS"
	TicketStatusWaiting    TicketStatus = "WAITING_USER"
	TicketStatusResolved   TicketStatus = "RESOLVED"
	TicketStatusClosed     TicketStatus = "CLOSED"
)

type TicketPriority string

const (
	PriorityLow    TicketPriority = "LOW"
	PriorityNormal TicketPriority = "NORMAL"
	PriorityHigh   TicketPriority = "HIGH"
	PriorityUrgent TicketPriority = "URGENT"
)

type Ticket struct {
	ID          int            `json:"id"`
	Code        string         `json:"code"`
	CategoryID  int            `json:"category_id"`
	RequesterID int            `json:"requester_id"`
	AssigneeID  *int           `json:"assignee_id,omitempty"`
	Priority    TicketPriority `json:"priority"`
	Status      TicketStatus   `json:"status"`
	Subject     string         `json:"subject"`
	Description string         `json:"description"`
	DueAt       time.Time      `json:"due_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	ResolvedAt  *time.Time     `json:"resolved_at,omitempty"`
}

type CreateTicketRequest struct {
	CategoryID  int            `json:"category_id" validate:"required,min=1"`
	Subject     string         `json:"subject" validate:"required,min=5,max=150"`
	Description string         `json:"description" validate:"required,min=10,max=5000"`
	Priority    TicketPriority `json:"priority" validate:"omitempty,oneof=LOW NORMAL HIGH URGENT"`
}

type AssignTicketRequest struct {
	AssigneeID int `json:"assignee_id" validate:"required,min=1"`
}

type UpdateTicketStatusRequest struct {
	Status TicketStatus `json:"status" validate:"required,oneof=OPEN ASSIGNED IN_PROGRESS WAITING_USER RESOLVED CLOSED"`
}

package response

import "prakarsa-app/entity"

type CreateNotificationRes struct {
	ID string `json:"id"`
}

// Get List Response
type GetListNotificationRes struct {
	ID            string         `json:"id"`
	Type          string         `json:"type"`
	ReferenceType string         `json:"reference_type"`
	ReferenceID   string         `json:"reference_id"`
	Title         string         `json:"title"`
	Message       string         `json:"message"`
	ActionURL     *string        `json:"action_url"`
	Profile       entity.Profile `json:"profile"`
	Priority      string         `json:"priority"`
	IsRead        bool           `json:"is_read"`
	ReadAt        *int64         `json:"read_at"`
	IsActive      *bool          `json:"is_active"`
	CreatedAt     int64          `json:"created_at"`
}

// Get Detail Response
type GetDetailNotificationRes struct {
	ID            string         `json:"id"`
	Type          string         `json:"type"`
	ReferenceType string         `json:"reference_type"`
	ReferenceID   string         `json:"reference_id"`
	Title         string         `json:"title"`
	Message       string         `json:"message"`
	ActionURL     *string        `json:"action_url"`
	Profile       entity.Profile `json:"profile"`
	Priority      string         `json:"priority"`
	IsRead        bool           `json:"is_read"`
	ReadAt        *int64         `json:"read_at"`
	IsActive      *bool          `json:"is_active"`
	CreatedAt     int64          `json:"created_at"`
	UpdatedAt     *int64         `json:"updated_at"`
}

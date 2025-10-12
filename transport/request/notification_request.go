package request

import (
	"prakarsa-app/utils"

	validation "github.com/go-ozzo/ozzo-validation"
)

// CreateNotificationReq represent create request body
type CreateNotificationReq struct {
	UserID        string `json:"user_id"`
	Type          string `json:"type"`
	ReferenceType string `json:"reference_type"`
	ReferenceID   string `json:"reference_id"`
	SourceUserID  string
	Title         string `json:"title"`
	Message       string `json:"message"`
	ActionURL     string `json:"action_url"`
	Priority      string `json:"priority"`
}

func (request CreateNotificationReq) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.UserID, validation.Required),
		validation.Field(&request.Type, validation.Required),
		validation.Field(&request.ReferenceType, validation.Required),
		validation.Field(&request.ReferenceID, validation.Required),
		validation.Field(&request.SourceUserID, validation.Required),
		validation.Field(&request.Priority, validation.Required, validation.In(utils.AllowedNotificationPriority...)),
	)
}

// Update request body
type UpdateNotificationReq struct {
	ID          string `param:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    *bool  `json:"is_active"`
	UserID      string
}

func (request UpdateNotificationReq) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.ID, validation.Required),
		validation.Field(&request.UserID, validation.Required),
	)
}

// Delete request body
type DeleteNotificationReq struct {
	ID     string `param:"id"`
	UserID string
}

func (request DeleteNotificationReq) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.ID, validation.Required),
		validation.Field(&request.UserID, validation.Required),
	)
}

// GetList request body
type GetListNotificationReq struct {
	PerPage  int64  `query:"per_page"`
	Page     int64  `query:"page"`
	Type     string `query:"type"`
	IsRead   *bool  `query:"is_read"`
	IsActive *bool  `query:"is_active"`
	UserID   string
}

func (request GetListNotificationReq) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.UserID, validation.Required),
	)
}

// GetDetail request body
type GetDetailNotificationReq struct {
	ID     string `param:"id"`
	UserID string
}

func (request GetDetailNotificationReq) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.ID, validation.Required),
		validation.Field(&request.UserID, validation.Required),
	)
}

// Mark Read request body
type MarkReadNotificationReq struct {
	ID     string `param:"id"`
	UserID string
}

func (request MarkReadNotificationReq) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.ID, validation.Required),
		validation.Field(&request.UserID, validation.Required),
	)
}

// Mark Read All request body
type MarkReadAllNotificationReq struct {
	UserID string
}

func (request MarkReadAllNotificationReq) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.UserID, validation.Required),
	)
}

package request

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"prakarsa-app/utils"
)

// CreateNotificationReq represent create request body
type CreateNotificationReq struct {
	UserID        string `json:"user_id"`
	Type          string `json:"type"`
	ReferenceType string `json:"reference_type"`
	ReferenceID   string `json:"reference_id"`
	//SourceUserID  string `json:"source_user_id"`
	Title     string `json:"title"`
	Message   string `json:"message"`
	ActionURL string `json:"action_url"`
	Priority  string `json:"priority"`
}

func (request CreateNotificationReq) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.UserID, validation.Required),
		validation.Field(&request.Type, validation.Required, validation.In(utils.AllowedNotificationType...).Error(utils.ValidationOneOfMsg("type", utils.AllowedNotificationType))),
		validation.Field(&request.ReferenceType, validation.Required, validation.In(utils.AllowedNotificationReferenceType...)),
		validation.Field(&request.ReferenceID, validation.Required),
		//validation.Field(&request.SourceUserID, validation.Required),
		validation.Field(&request.Title, validation.Required),
		validation.Field(&request.Message, validation.Required),
		validation.Field(&request.Priority, validation.Required, validation.In(utils.AllowedNotificationPriority...)),
	)
}

// Update request body
type UpdateNotificationReq struct {
	ID          string `param:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    *bool  `json:"is_active"`
}

func (request UpdateNotificationReq) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.ID, validation.Required),
	)
}

// Delete request body
type DeleteNotificationReq struct {
	ID string `param:"id"`
}

func (request DeleteNotificationReq) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.ID, validation.Required),
	)
}

// GetList request body
type GetListNotificationReq struct {
	PerPage  int64  `query:"per_page"`
	Page     int64  `query:"page"`
	Type     string `query:"type"`
	IsRead   *bool  `query:"is_read"`
	IsActive *bool  `query:"is_active"`
	UserID   string `query:"user_id"`
}

func (request GetListNotificationReq) Validate() error {
	return validation.ValidateStruct(
		&request,
	)
}

// GetDetail request body
type GetDetailNotificationReq struct {
	ID string `param:"id"`
}

func (request GetDetailNotificationReq) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.ID, validation.Required),
	)
}

// Mark Read request body
type MarkReadNotificationReq struct {
	ID string `param:"id"`
}

func (request MarkReadNotificationReq) Validate() error {
	return validation.ValidateStruct(
		&request,
	)
}

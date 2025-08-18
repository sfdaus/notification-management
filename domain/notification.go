package domain

import (
	"context"
	"prakarsa-app/transport/request"
	"prakarsa-app/transport/response"
)

type Notification struct {
	ID            string  `json:"id"`
	UserID        string  `json:"user_id"`
	Type          string  `json:"type"`           // e.g. "MENTION", "COMMENT"
	ReferenceType string  `json:"reference_type"` // e.g. "comment", "thread"
	ReferenceID   string  `json:"reference_id"`
	SourceUserID  *string `json:"source_user_id"`
	Title         string  `json:"title"`      // varchar(100)
	Message       string  `json:"message"`    // text
	ActionURL     *string `json:"action_url"` // deep link / web url
	Priority      int32   `json:"priority"`   // 0=low,1=normal,2=high

	IsRead   bool   `json:"is_read"`
	ReadAt   *int64 `json:"read_at"`
	IsActive *bool  `json:"is_active"`

	CreatedBy string `json:"created_by"`
	CreatedAt int64  `json:"created_at"`
	UpdatedBy string `json:"updated_by"`
	UpdatedAt int64  `json:"updated_at"`
	DeletedAt int64  `json:"deleted_at"`
}

// // NotificationRepository represent the notification repository contract
type NotificationRepository interface {
	Create(ctx context.Context, notification *Notification) error
	Update(ctx context.Context, notification *Notification) error
	Delete(ctx context.Context, notification *Notification) (int64, error)
	GetList(ctx context.Context, request *request.GetListNotificationReq) ([]response.GetListNotificationRes, response.MetaRes, error)
	GetDetail(ctx context.Context, request *request.GetDetailNotificationReq) (response.GetDetailNotificationRes, error)
	MarkRead(ctx context.Context, notification *Notification) error
}

// NotificationUsecase represent the notification usecase contract
type NotificationUsecase interface {
	Create(ctx context.Context, request *request.CreateNotificationReq) (response.CreateNotificationRes, error)
	Update(ctx context.Context, request *request.UpdateNotificationReq) error
	Delete(ctx context.Context, request *request.DeleteNotificationReq) (int64, error)
	GetList(ctx context.Context, request *request.GetListNotificationReq) ([]response.GetListNotificationRes, response.MetaRes, error)
	GetDetail(ctx context.Context, request *request.GetDetailNotificationReq) (response.GetDetailNotificationRes, error)
	MarkRead(ctx context.Context, request *request.MarkReadNotificationReq) error
}

package domain

import (
	"context"
	"prakarsa-app/entity"
	"prakarsa-app/transport/request"
	"prakarsa-app/transport/response"
)

// // NotificationRepository represent the notification repository contract
type NotificationRepository interface {
	Create(ctx context.Context, notification *entity.Notification) error
	Update(ctx context.Context, notification *entity.Notification) error
	Delete(ctx context.Context, notification *entity.Notification) (int64, error)
	GetList(ctx context.Context, request *request.GetListNotificationReq) ([]response.GetListNotificationRes, response.MetaRes, error)
	GetDetail(ctx context.Context, request *request.GetDetailNotificationReq) (response.GetDetailNotificationRes, error)
	MarkRead(ctx context.Context, notification *entity.Notification) error
	MarkReadAll(ctx context.Context, notification *entity.Notification) error
}

// NotificationUsecase represent the notification usecase contract
type NotificationUsecase interface {
	Create(ctx context.Context, request *request.CreateNotificationReq) (response.CreateNotificationRes, error)
	Update(ctx context.Context, request *request.UpdateNotificationReq) error
	Delete(ctx context.Context, request *request.DeleteNotificationReq) (int64, error)
	GetList(ctx context.Context, request *request.GetListNotificationReq) ([]response.GetListNotificationRes, response.MetaRes, error)
	GetDetail(ctx context.Context, request *request.GetDetailNotificationReq) (response.GetDetailNotificationRes, error)
	MarkRead(ctx context.Context, request *request.MarkReadNotificationReq) error
	MarkReadAll(ctx context.Context, request *request.MarkReadAllNotificationReq) error
}

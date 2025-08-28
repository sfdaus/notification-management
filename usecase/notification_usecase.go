package usecase

import (
	"context"
	"prakarsa-app/transport/response"
	"prakarsa-app/utils"
	"time"

	"prakarsa-app/domain"
	"prakarsa-app/repository/redis"
	"prakarsa-app/transport/request"

	"github.com/google/uuid"
)

type NotificationUsecase struct {
	notificationRepo domain.NotificationRepository
	redisRepo        redis.RedisRepository
	ctxTimeout       time.Duration
}

// NewNotificationUsecase will create new an notificationUsecase object representation of ThreadUsecase interface
func NewNotificationUsecase(notificationRepo domain.NotificationRepository, redisRepo redis.RedisRepository, ctxTimeout time.Duration) *NotificationUsecase {
	return &NotificationUsecase{
		notificationRepo: notificationRepo,
		redisRepo:        redisRepo,
		ctxTimeout:       ctxTimeout,
	}
}

func (u *NotificationUsecase) Create(c context.Context, request *request.CreateNotificationReq) (res response.CreateNotificationRes, err error) {
	ctx, cancel := context.WithTimeout(c, u.ctxTimeout)
	defer cancel()

	// Create Payload
	notificationID := uuid.NewString()
	t := true
	srcID := request.SourceUserID
	notificationPayload := &domain.Notification{
		ID:            notificationID,
		UserID:        request.UserID,
		Type:          request.Type,
		ReferenceType: request.ReferenceType,
		ReferenceID:   request.ReferenceID,
		SourceUserID:  &srcID,
		Title:         request.Title,
		Message:       request.Message,
		ActionURL:     &request.ActionURL,
		Priority:      utils.NotificationPriority[request.Priority],
		IsRead:        false,
		IsActive:      &t,
		CreatedBy:     request.SourceUserID,
		CreatedAt:     time.Now().Unix(),
	}

	// Response Payload
	res.ID = notificationID

	err = u.notificationRepo.Create(ctx, notificationPayload)
	return
}

func (u *NotificationUsecase) Update(c context.Context, request *request.UpdateNotificationReq) (err error) {
	ctx, cancel := context.WithTimeout(c, u.ctxTimeout)
	defer cancel()

	// Update Payload
	notificationPayload := &domain.Notification{
		ID:        request.ID,
		UpdatedBy: request.UserID,
		UpdatedAt: time.Now().Unix(),
	}

	if request.IsActive != nil {
		notificationPayload.IsActive = request.IsActive
	}

	err = u.notificationRepo.Update(ctx, notificationPayload)
	return
}
func (u *NotificationUsecase) Delete(c context.Context, request *request.DeleteNotificationReq) (rowsAffected int64, err error) {
	ctx, cancel := context.WithTimeout(c, u.ctxTimeout)
	defer cancel()

	threadPayload := &domain.Notification{
		ID: request.ID,
	}

	rowsAffected, err = u.notificationRepo.Delete(ctx, threadPayload)
	return
}

func (u *NotificationUsecase) GetList(c context.Context, request *request.GetListNotificationReq) (res []response.GetListNotificationRes, meta response.MetaRes, err error) {
	ctx, cancel := context.WithTimeout(c, u.ctxTimeout)
	defer cancel()

	res, meta, err = u.notificationRepo.GetList(ctx, request)
	return
}

func (u *NotificationUsecase) GetDetail(c context.Context, request *request.GetDetailNotificationReq) (res response.GetDetailNotificationRes, err error) {
	ctx, cancel := context.WithTimeout(c, u.ctxTimeout)
	defer cancel()

	res, err = u.notificationRepo.GetDetail(ctx, request)
	return
}

func (u *NotificationUsecase) MarkRead(c context.Context, request *request.MarkReadNotificationReq) (err error) {
	ctx, cancel := context.WithTimeout(c, u.ctxTimeout)
	defer cancel()

	var notification = &domain.Notification{
		ID:        request.ID,
		UserID:    request.UserID,
		UpdatedAt: time.Now().Unix(),
		UpdatedBy: request.UserID,
	}

	err = u.notificationRepo.MarkRead(ctx, notification)
	return
}

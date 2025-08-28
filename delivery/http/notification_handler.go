package http

import (
	"net/http"
	"prakarsa-app/delivery/middleware"
	"prakarsa-app/domain"
	"prakarsa-app/transport/request"
	"prakarsa-app/utils"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/labstack/echo/v4"
)

type NotificationHandler struct {
	NotificationUC domain.NotificationUsecase
}

// NewNotificationHandler will initialize the todo resources endpoint
func NewNotificationHandler(e *echo.Echo, middleware *middleware.Middleware, notificationUC domain.NotificationUsecase) {
	handler := &NotificationHandler{
		NotificationUC: notificationUC,
	}

	apiV1 := e.Group("/api/v1")
	apiV1.POST("/notifications", handler.Create)
	apiV1.PATCH("/notifications/:id", handler.Update)
	apiV1.PATCH("/notifications/:id/read", handler.MarkRead)
	apiV1.DELETE("/notifications/:id", handler.Delete)
	apiV1.GET("/notifications", handler.GetList)
	apiV1.GET("/notifications/:id", handler.GetDetail)
}

func (h *NotificationHandler) Create(c echo.Context) error {
	ctx := c.Request().Context()
	var req request.CreateNotificationReq

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewUnprocessableEntityError(err.Error()))
	}

	req.SourceUserID = c.Request().Header.Get("x-user-id")
	req.Priority = strings.ToLower(req.Priority)

	if err := req.Validate(); err != nil {
		errVal := err.(validation.Errors)
		return c.JSON(http.StatusBadRequest, utils.NewInvalidInputError(errVal))
	}

	if res, err := h.NotificationUC.Create(ctx, &req); err != nil {
		return c.JSON(utils.ParseHttpError(err))
	} else {
		return c.JSON(http.StatusCreated, map[string]interface{}{
			"message": "Notification successfully created",
			"data":    res,
		})
	}

}

func (h *NotificationHandler) Update(c echo.Context) error {
	ctx := c.Request().Context()
	var req request.UpdateNotificationReq

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewUnprocessableEntityError(err.Error()))
	}

	req.UserID = c.Request().Header.Get("x-user-id")

	if err := req.Validate(); err != nil {
		errVal := err.(validation.Errors)
		return c.JSON(http.StatusBadRequest, utils.NewInvalidInputError(errVal))
	}

	if err := h.NotificationUC.Update(ctx, &req); err != nil {
		return c.JSON(utils.ParseHttpError(err))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Notification successfully updated",
	})
}

func (h *NotificationHandler) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	var req request.DeleteNotificationReq

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewUnprocessableEntityError(err.Error()))
	}

	req.UserID = c.Request().Header.Get("x-user-id")

	if err := req.Validate(); err != nil {
		errVal := err.(validation.Errors)
		return c.JSON(http.StatusBadRequest, utils.NewInvalidInputError(errVal))
	}

	if rowsAffected, err := h.NotificationUC.Delete(ctx, &req); err != nil {
		return c.JSON(utils.ParseHttpError(err))
	} else {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Notification successfully deleted",
			"data": map[string]int64{
				"rows_affected": rowsAffected,
			},
		})
	}
}

func (h *NotificationHandler) GetList(c echo.Context) error {
	ctx := c.Request().Context()
	var req request.GetListNotificationReq

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewUnprocessableEntityError(err.Error()))
	}

	if err := req.Validate(); err != nil {
		errVal := err.(validation.Errors)
		return c.JSON(http.StatusBadRequest, utils.NewInvalidInputError(errVal))
	}

	if res, meta, err := h.NotificationUC.GetList(ctx, &req); err != nil {
		return c.JSON(utils.ParseHttpError(err))
	} else {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Notification successfully retrieved",
			"data": map[string]interface{}{
				"data": res,
				"meta": meta,
			},
		})
	}
}

func (h *NotificationHandler) GetDetail(c echo.Context) error {
	ctx := c.Request().Context()
	var req request.GetDetailNotificationReq

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewUnprocessableEntityError(err.Error()))
	}

	if err := req.Validate(); err != nil {
		errVal := err.(validation.Errors)
		return c.JSON(http.StatusBadRequest, utils.NewInvalidInputError(errVal))
	}

	if res, err := h.NotificationUC.GetDetail(ctx, &req); err != nil {
		return c.JSON(utils.ParseHttpError(err))
	} else {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Notification successfully retrieved",
			"data":    res,
		})
	}
}

func (h *NotificationHandler) MarkRead(c echo.Context) error {
	ctx := c.Request().Context()
	var req request.MarkReadNotificationReq

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewUnprocessableEntityError(err.Error()))
	}

	req.UserID = c.Request().Header.Get("x-user-id")

	if err := req.Validate(); err != nil {
		errVal := err.(validation.Errors)
		return c.JSON(http.StatusBadRequest, utils.NewInvalidInputError(errVal))
	}

	if err := h.NotificationUC.MarkRead(ctx, &req); err != nil {
		return c.JSON(utils.ParseHttpError(err))
	} else {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Notification successfully marked as read",
		})
	}
}

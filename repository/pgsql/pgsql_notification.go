package pgsql

import (
	"context"
	"database/sql"
	"fmt"
	"prakarsa-app/domain"
	"prakarsa-app/transport/request"
	"prakarsa-app/transport/response"
	"prakarsa-app/utils"
	"strings"
)

type pgsqlNotificationRepository struct {
	db *sql.DB
}

// NewPgsqlNotificationRepository will create new an todoRepository object representation of NotificationRepository interface
func NewPgsqlNotificationRepository(db *sql.DB) *pgsqlNotificationRepository {
	return &pgsqlNotificationRepository{
		db: db,
	}
}

func (r *pgsqlNotificationRepository) Create(ctx context.Context, notification *domain.Notification) (err error) {
	query := `INSERT INTO notifications (
				id, user_id, type, reference_type, reference_id, source_user_id,
				title, message, action_url, priority,
				is_read, is_active,
				created_by, created_at
			) VALUES (
				$1, $2, $3, $4, $5, $6,
				$7, $8, $9, $10,
				$11, $12, $13,
				$14
			)`
	if _, err = r.db.ExecContext(
		ctx, query,
		notification.ID,
		notification.UserID,
		notification.Type,
		notification.ReferenceType,
		notification.ReferenceID,
		notification.SourceUserID,
		notification.Title,
		notification.Message,
		notification.ActionURL,
		notification.Priority,
		notification.IsRead,
		notification.IsActive,
		notification.CreatedBy,
		notification.CreatedAt,
	); err != nil {
		return err
	}

	return
}

func (r *pgsqlNotificationRepository) Update(ctx context.Context, notification *domain.Notification) (err error) {
	// Build dynamic SET clauses from Notification struct
	sets := []string{}
	args := []interface{}{}
	idx := 1

	//if notification.Name != "" {
	//	sets = append(sets, fmt.Sprintf("name = $%d", idx))
	//	args = append(args, notification.Name)
	//	idx++
	//}
	//
	//if notification.Description != "" {
	//	sets = append(sets, fmt.Sprintf("description = $%d", idx))
	//	args = append(args, notification.Description)
	//	idx++
	//}

	if notification.IsActive != nil {
		sets = append(sets, fmt.Sprintf("is_active = $%d", idx))
		args = append(args, notification.IsActive)
		idx++
	}

	// kalau ada sesuatu untuk di‐update, commit ke SQL
	if len(sets) > 0 {
		// Update stamp
		sets = append(sets, fmt.Sprintf("updated_at = $%d", idx))
		args = append(args, notification.UpdatedAt)
		idx++

		sets = append(sets, fmt.Sprintf("updated_by = $%d", idx))
		args = append(args, notification.UpdatedBy)
		idx++

		// tambahkan WHERE id = $idx
		args = append(args, notification.ID)
		query := fmt.Sprintf(
			"UPDATE notifications SET %s WHERE id = $%d",
			strings.Join(sets, ", "),
			idx,
		)

		if _, err = r.db.ExecContext(ctx, query, args...); err != nil {
			return
		}
	}

	return
}

func (r *pgsqlNotificationRepository) Delete(ctx context.Context, notification *domain.Notification) (rowsAffected int64, err error) {
	query := "DELETE FROM notifications WHERE id = $1"
	res, err := r.db.ExecContext(ctx, query, notification.ID)
	if err != nil {
		return
	}

	rowsAffected, err = res.RowsAffected()
	if err != nil {
		return
	}

	return
}

func (r *pgsqlNotificationRepository) GetList(ctx context.Context, request *request.GetListNotificationReq) (res []response.GetListNotificationRes, meta response.MetaRes, err error) {
	// 1. Build WHERE clauses
	wheres := []string{fmt.Sprintf("user_id = $%d", 1)}
	args := []interface{}{request.UserID}
	idx := 2

	if request.Type != "" {
		wheres = append(wheres, fmt.Sprintf("name = $%d", idx))
		args = append(args, request.Type)
		idx++
	}

	if request.IsActive != nil {
		wheres = append(wheres, fmt.Sprintf("is_active = $%d", idx))
		args = append(args, request.IsActive)
		idx++
	}

	if request.IsRead != nil {
		wheres = append(wheres, fmt.Sprintf("is_read = $%d", idx))
		args = append(args, request.IsRead)
		idx++
	}

	whereSQL := ""
	if len(wheres) > 0 {
		whereSQL = "WHERE " + strings.Join(wheres, " AND ")
	}

	// --- 2. Hitung totalCount dulu (tanpa LIMIT/OFFSET) ---
	countQuery := fmt.Sprintf(
		"SELECT COUNT(*) FROM notifications %s",
		whereSQL,
	)

	if err = r.db.QueryRowContext(ctx, countQuery, args...).Scan(&meta.TotalData); err != nil {
		return nil, meta, err
	}

	// 2. Calculate LIMIT & OFFSET
	perPage := request.PerPage
	if perPage <= 0 {
		perPage = 10
	}
	page := request.Page
	if page <= 0 {
		page = 1
	}

	// total pages = ceil(total / perPage)
	meta.Page = page
	meta.PerPage = perPage
	meta.TotalPages = (meta.TotalData + perPage - 1) / perPage

	offset := (page - 1) * perPage

	// add LIMIT & OFFSET to args
	args = append(args, perPage, offset)
	limitPos, offsetPos := idx, idx+1

	// 3. Final query
	query := fmt.Sprintf(`
        SELECT
            id, type, title, message, action_url, priority, is_read, read_At,
            is_active, created_at
        FROM notifications
        %s
        ORDER BY created_at DESC
        LIMIT $%d OFFSET $%d
    `, whereSQL, limitPos, offsetPos)

	// 4. Execute
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, meta, err
	}
	defer rows.Close()

	// 5. Scan results
	for rows.Next() {
		var item response.GetListNotificationRes

		if err := rows.Scan(
			&item.ID,
			&item.Type,
			&item.Title,
			&item.Message,
			&item.ActionURL,
			&item.Priority,
			&item.IsRead,
			&item.ReadAt,
			&item.IsActive,
			&item.CreatedAt,
		); err != nil {
			return nil, meta, err
		}

		res = append(res, item)
	}
	if errRow := rows.Err(); errRow != nil {
		return nil, meta, errRow
	}

	return
}

func (r *pgsqlNotificationRepository) GetDetail(ctx context.Context, request *request.GetDetailNotificationReq) (res response.GetDetailNotificationRes, err error) {

	const query = `
					SELECT
						id, type, title, message, action_url, priority, is_read, read_At,
						is_active, created_at
					FROM notifications
					WHERE id = $1
					LIMIT 1
					`

	// 1. QueryRowContext untuk ambil satu baris
	row := r.db.QueryRowContext(ctx, query, request.ID)

	// 2. Scan kolom ke field di domain.Notification
	// since created_at is NOT NULL int8:
	var createdAt int64
	// updated_at/deleted_at can be NULL, so use NullInt64:
	var updatedAt sql.NullInt64

	err = row.Scan(
		&res.ID,
		&res.Type,
		&res.Title,
		&res.Message,
		&res.ActionURL,
		&res.Priority,
		&res.IsRead,
		&res.ReadAt,
		&res.IsActive,
		&res.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return res, utils.NewNotFoundError("notification not found")
		}
		return res, err
	}

	// assign into your domain fields
	res.CreatedAt = createdAt
	if updatedAt.Valid {
		res.UpdatedAt = updatedAt.Int64
	}

	return
}

func (r *pgsqlNotificationRepository) MarkRead(ctx context.Context, notification *domain.Notification) (err error) {
	const query = `
		UPDATE notifications
		SET is_read = TRUE,
		    read_at = COALESCE(read_at, $3),
		    updated_at = $3,
		    updated_by = $2
		WHERE id = $1 AND is_active = TRUE AND user_id = $4;
	`
	_, err = r.db.ExecContext(ctx, query, notification.ID, notification.UpdatedBy, notification.UpdatedAt, notification.UserID)

	return
}

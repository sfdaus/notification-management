package pgsql

import (
	"context"
	"database/sql"
	"fmt"
	"prakarsa-app/entity"
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

func (r *pgsqlNotificationRepository) Create(ctx context.Context, notification *entity.Notification) (err error) {
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

func (r *pgsqlNotificationRepository) Update(ctx context.Context, notification *entity.Notification) (err error) {
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

func (r *pgsqlNotificationRepository) Delete(ctx context.Context, notification *entity.Notification) (rowsAffected int64, err error) {
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
	// Mulai transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return
	}

	// Pastikan rollback kalau ada error
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 1. Build WHERE clauses
	wheres := []string{fmt.Sprintf("n.user_id = $%d", 1)}
	args := []interface{}{request.UserID}
	idx := 2

	if request.Type != "" {
		wheres = append(wheres, fmt.Sprintf("n.type = $%d", idx))
		args = append(args, request.Type)
		idx++
	}

	if request.IsActive != nil {
		wheres = append(wheres, fmt.Sprintf("n.is_active = $%d", idx))
		args = append(args, request.IsActive)
		idx++
	} else {
		// default only active
		wheres = append(wheres, fmt.Sprintf("n.is_active = $%d", idx))
		args = append(args, true)
		idx++
	}

	if request.IsRead != nil {
		wheres = append(wheres, fmt.Sprintf("n.is_read = $%d", idx))
		args = append(args, request.IsRead)
		idx++
	}

	whereSQL := ""
	if len(wheres) > 0 {
		whereSQL = "WHERE " + strings.Join(wheres, " AND ")
	}

	// --- 2. Hitung totalCount dulu (tanpa LIMIT/OFFSET) ---
	countQuery := fmt.Sprintf(
		"SELECT COUNT(*) FROM notifications n %s",
		whereSQL,
	)

	if err = tx.QueryRowContext(ctx, countQuery, args...).Scan(&meta.TotalData); err != nil {
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
            n.id, n.type, n.reference_type, n.reference_id, n.title, 
			n.message, n.action_url, n.priority, n.is_read, n.read_at,
            n.is_active, n.created_at,

			-- profile (1-1)
			CASE 
                WHEN p.user_id = $1 THEN 'SYSTEM'
                ELSE COALESCE(p.name,'')
            END AS prof_name,
			CASE 
                WHEN p.user_id = $1 THEN 'SYSTEM'
                ELSE COALESCE(p.name_alias,'')
            END AS prof_name_alias,
			COALESCE(p.avatar,'') AS prof_avatar,
		
			-- institution inside profile
			COALESCE(i.name,'')  AS prof_inst_name,
			COALESCE(i.alias,'') AS prof_inst_alias,
			COALESCE(i.type,'')  AS prof_inst_type

        FROM notifications n
		LEFT JOIN profiles p ON n.source_user_id = p.user_id
		LEFT JOIN institutions i ON p.institution_id = i.id
        %s
        ORDER BY n.created_at DESC
        LIMIT $%d OFFSET $%d
    `, whereSQL, limitPos, offsetPos)

	// 4. Execute
	rows, err := tx.QueryContext(ctx, query, args...)
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
			&item.ReferenceType,
			&item.ReferenceID,
			&item.Title,
			&item.Message,
			&item.ActionURL,
			&item.Priority,
			&item.IsRead,
			&item.ReadAt,
			&item.IsActive,
			&item.CreatedAt,
			&item.Profile.Name,
			&item.Profile.NameAlias,
			&item.Profile.Avatar,
			&item.Profile.Institution.Name,
			&item.Profile.Institution.Alias,
			&item.Profile.Institution.Type,
		); err != nil {
			return nil, meta, err
		}

		res = append(res, item)
	}
	if errRow := rows.Err(); errRow != nil {
		return nil, meta, errRow
	}

	// Commit jika semua sukses
	if err = tx.Commit(); err != nil {
		return
	}

	return
}

func (r *pgsqlNotificationRepository) GetDetail(ctx context.Context, request *request.GetDetailNotificationReq) (res response.GetDetailNotificationRes, err error) {

	const query = `
					SELECT
						n.id, n.type, n.reference_type, n.reference_id, n.title, 
						n.message, n.action_url, n.priority, n.is_read, n.read_at,
						n.is_active, n.created_at, n.updated_at,

						-- profile (1-1)
						CASE 
							WHEN p.user_id = $1 THEN 'SYSTEM'
							ELSE COALESCE(p.name,'')
						END AS prof_name,
						CASE 
							WHEN p.user_id = $1 THEN 'SYSTEM'
							ELSE COALESCE(p.name_alias,'')
						END AS prof_name_alias,
						COALESCE(p.avatar,'') AS prof_avatar,
					
						-- institution inside profile
						COALESCE(i.name,'')  AS prof_inst_name,
						COALESCE(i.alias,'') AS prof_inst_alias,
						COALESCE(i.type,'')  AS prof_inst_type

					FROM notifications n
					LEFT JOIN profiles p ON n.source_user_id = p.user_id
					LEFT JOIN institutions i ON p.institution_id = i.id
					WHERE n.id = $1
					LIMIT 1
					`

	// 1. QueryRowContext untuk ambil satu baris
	row := r.db.QueryRowContext(ctx, query, request.ID)

	err = row.Scan(
		&res.ID,
		&res.Type,
		&res.ReferenceType,
		&res.ReferenceID,
		&res.Title,
		&res.Message,
		&res.ActionURL,
		&res.Priority,
		&res.IsRead,
		&res.ReadAt,
		&res.IsActive,
		&res.CreatedAt,
		&res.UpdatedAt,
		&res.Profile.Name,
		&res.Profile.NameAlias,
		&res.Profile.Avatar,
		&res.Profile.Institution.Name,
		&res.Profile.Institution.Alias,
		&res.Profile.Institution.Type,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return res, utils.NewNotFoundError("notification not found")
		}
		return res, err
	}

	return
}

func (r *pgsqlNotificationRepository) MarkRead(ctx context.Context, notification *entity.Notification) (err error) {
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

func (r *pgsqlNotificationRepository) MarkReadAll(ctx context.Context, notifications *entity.Notification) (err error) {
	const query = `
		UPDATE notifications
		SET is_read = TRUE,
		    read_at = COALESCE(read_at, $2),
		    updated_at = $2,
		    updated_by = $1
		WHERE is_active = TRUE AND user_id = $3;
	`
	_, err = r.db.ExecContext(ctx, query, notifications.UpdatedBy, notifications.UpdatedAt, notifications.UserID)

	return
}

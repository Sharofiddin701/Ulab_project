package postgres

import (
	"context"
	"database/sql"
	"e-commerce/models"
	"e-commerce/pkg/logger"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lib/pq"
)

type colorRepo struct {
	db  *pgxpool.Pool
	log logger.LoggerI
}

func NewColorRepo(db *pgxpool.Pool, log logger.LoggerI) *colorRepo {
	return &colorRepo{
		db:  db,
		log: log,
	}
}

func (u *colorRepo) Create(ctx context.Context, req *models.ColorCreate) (*models.Color, error) {
	var (
		id = uuid.New().String()
	)

	query := `
        INSERT INTO "color" (
            id,
            product_id,
            color_name,
            color_url,
			count,
            created_at
        )
        VALUES($1, $2, $3, $4, $5, CURRENT_TIMESTAMP)
        RETURNING id, product_id, color_name, color_url, count, created_at
    `

	var (
		idd        sql.NullString
		product_id sql.NullString
		name       sql.NullString
		color_url  pq.StringArray
		count      sql.NullInt32
		created_at sql.NullTime
	)

	err := u.db.QueryRow(ctx, query, id, req.ProductId, req.Name, req.Url, req.Count).Scan(
		&idd,
		&product_id,
		&name,
		&color_url,
		&count,
		&created_at,
	)
	if err != nil {
		u.log.Error("Error while inserting color: " + err.Error())
		return nil, err
	}

	return &models.Color{
		Id:        idd.String,
		ProductId: req.ProductId,
		Name:      name.String,
		Url:       req.Url,
		Count:     int(count.Int32),
		CreatedAt: created_at.Time.Format(time.RFC3339),
	}, nil
}

func (u *colorRepo) GetList(ctx context.Context, req *models.ColorGetListRequest) (*models.ColorGetListResponse, error) {
	resp := &models.ColorGetListResponse{}
	offset := " OFFSET 0"
	limit := " LIMIT 10"

	query := `
		SELECT
			COUNT(*) OVER(),
			id,
			product_id,
			color_name,
			color_url,
			count,
			created_at
		FROM "color"
	`

	if req.Offset > 0 {
		offset = fmt.Sprintf(" OFFSET %d", req.Offset)
	}

	if req.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", req.Limit)
	}

	query += offset + limit
	rows, err := u.db.Query(ctx, query)
	if err != nil {
		u.log.Error("Error while getting color list: " + err.Error())
		return nil, err
	}

	for rows.Next() {
		var (
			id         sql.NullString
			product_id sql.NullString
			color_name sql.NullString
			color_url  pq.StringArray
			count      sql.NullInt32
			created_at sql.NullString
		)

		err = rows.Scan(
			&resp.Count,
			&id,
			&product_id,
			&color_name,
			&color_url,
			&count,
			&created_at,
		)
		if err != nil {
			u.log.Error("Error while scanning color list data: " + err.Error())
			return nil, err
		}

		resp.Color = append(resp.Color, &models.Color{
			Id:        id.String,
			ProductId: product_id.String,
			Name:      color_name.String,
			Url:       color_url,
			Count:     int(count.Int32),
			CreatedAt: created_at.String,
		})
	}
	return resp, nil
}

func (u *colorRepo) Delete(ctx context.Context, req *models.ColorPrimaryKey) error {
	_, err := u.db.Exec(ctx, `DELETE FROM "color" WHERE id = $1`, req.Id)
	if err != nil {
		u.log.Error("Error while deleting color: " + err.Error())
		return err
	}

	return nil
}

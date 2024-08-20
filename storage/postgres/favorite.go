package postgres

import (
	"context"
	"database/sql"
	"e-commerce/models"
	"e-commerce/pkg/logger"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type favoriteRepo struct {
	db  *pgxpool.Pool
	log logger.LoggerI
}

func NewFavoriteRepo(db *pgxpool.Pool, log logger.LoggerI) *favoriteRepo {
	return &favoriteRepo{
		db:  db,
		log: log,
	}
}

func (u *favoriteRepo) GetList(ctx context.Context, req *models.FavoriteGetListRequest) (*models.FavoriteGetListResponse, error) {
	var (
		resp   = &models.FavoriteGetListResponse{}
		query  string
		offset = " OFFSET 0"
		limit  = " LIMIT 10"
	)

	query = `
    SELECT
        COUNT(*) OVER() AS total,
        id AS product_id,
        created_at,
        updated_at
    FROM "product"
    WHERE is_favourite = false
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
		u.log.Error("error while getting favorite list: " + err.Error())
		return nil, err
	}

	for rows.Next() {
		var (
			productID sql.NullString
			createdAt sql.NullString
			updatedAt sql.NullString
		)

		err = rows.Scan(
			&resp.Total,
			&productID,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			u.log.Error("error while scanning favorite list data", logger.Error(err))
			return nil, err
		}

		resp.Favorite = append(resp.Favorite, models.Favorite{
			ProductID: productID.String,
			CreatedAt: createdAt.String,
			UpdatedAt: updatedAt.String,
		})
	}
	return resp, nil
}

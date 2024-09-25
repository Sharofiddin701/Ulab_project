package postgres

import (
	"context"
	"database/sql"
	"e-commerce/models"
	"e-commerce/pkg/helper"
	"e-commerce/pkg/logger"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type bannerRepo struct {
	db  *pgxpool.Pool
	log logger.LoggerI
}

// NewBannerRepo initializes a new instance of bannerRepo
func NewBannerRepo(db *pgxpool.Pool, log logger.LoggerI) *bannerRepo {
	return &bannerRepo{
		db:  db,
		log: log,
	}
}

// Create inserts a new banner into the database
func (u *bannerRepo) Create(ctx context.Context, req *models.BannerCreate) (*models.Banner, error) {
	id := uuid.New().String()
	query := `
		INSERT INTO "banner" (
			id,
			banner_image,
			created_at
		)
		VALUES ($1, $2, CURRENT_TIMESTAMP)
		RETURNING id, banner_image, created_at, updated_at, deleted_at
	`

	var (
		idd          sql.NullString
		banner_image sql.NullString
		created_at   sql.NullString
		updated_at   sql.NullString
		deleted_at   sql.NullString
	)

	err := u.db.QueryRow(ctx, query, id, req.Banner_image).Scan(
		&idd,
		&banner_image,
		&created_at,
		&updated_at,
		&deleted_at,
	)
	if err != nil {
		u.log.Error("Error while creating banner: " + err.Error())
		return nil, err
	}

	return &models.Banner{
		Id:           idd.String,
		Banner_image: banner_image.String,
		CreatedAt:    created_at.String,
		UpdatedAt:    updated_at.String,
		DeletedAt:    deleted_at.String,
	}, nil
}

// GetByID retrieves a banner by its ID
func (u *bannerRepo) GetByID(ctx context.Context, req *models.BannerPrimaryKey) (*models.Banner, error) {
	query := `
		SELECT 
			id,
			banner_image,
			created_at
		FROM "banner" 
		WHERE id = $1
	`

	var (
		id           sql.NullString
		banner_image sql.NullString
		created_at   sql.NullString
	)

	err := u.db.QueryRow(ctx, query, req.Id).Scan(
		&id,
		&banner_image,
		&created_at,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			u.log.Warn("No banner found with the given ID")
			return nil, nil
		}
		u.log.Error("Error while scanning banner data: " + err.Error())
		return nil, err
	}

	return &models.Banner{
		Id:           id.String,
		Banner_image: banner_image.String,
		CreatedAt:    created_at.String,
	}, nil
}

// GetList retrieves a list of banners with pagination
func (u *bannerRepo) GetList(ctx context.Context, req *models.BannerGetListRequest) (*models.BannerGetListResponse, error) {
	resp := &models.BannerGetListResponse{}
	offset := " OFFSET 0"
	limit := " LIMIT 10"

	query := `
		SELECT
			COUNT(*) OVER(),
			id,
			banner_image,
			created_at
		FROM "banner" 
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
		u.log.Error("Error while getting banner list: " + err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id           sql.NullString
			banner_image sql.NullString
			created_at   sql.NullString
		)

		err = rows.Scan(
			&resp.Count,
			&id,
			&banner_image,
			&created_at,
		)
		if err != nil {
			u.log.Error("Error while scanning banner list data: " + err.Error())
			return nil, err
		}

		resp.Banner = append(resp.Banner, &models.Banner{
			Id:           id.String,
			Banner_image: banner_image.String,
			CreatedAt:    created_at.String,
		})
	}

	return resp, nil
}

// Delete removes a banner from the database by ID
func (u *bannerRepo) Delete(ctx context.Context, req *models.BannerPrimaryKey) error {
	_, err := u.db.Exec(ctx, `DELETE FROM "banner" WHERE id = $1`, req.Id)
	if err != nil {
		u.log.Error("Error while deleting banner: " + err.Error())
		return err
	}

	return nil
}

// Update modifies a banner's data in the database
func (u *bannerRepo) Update(ctx context.Context, req *models.BannerUpdate) (int64, error) {
	query := `
		UPDATE "banner"
		SET
			banner_image = :banner_image,
			updated_at = NOW()
		WHERE id = :id
	`

	params := map[string]interface{}{
		"id":           req.Id,
		"banner_image": req.Banner_image,
	}

	query, args := helper.ReplaceQueryParams(query, params)
	result, err := u.db.Exec(ctx, query, args...)
	if err != nil {
		u.log.Error("Error while updating banner data: " + err.Error())
		return 0, err
	}

	return result.RowsAffected(), nil
}

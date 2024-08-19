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

type brandRepo struct {
	db  *pgxpool.Pool
	log logger.LoggerI
}

// NewAdminRepo initializes a new instance of adminRepo
func NewBrandRepo(db *pgxpool.Pool, log logger.LoggerI) *brandRepo {
	return &brandRepo{
		db:  db,
		log: log,
	}
}

// Create inserts a new admin into the database
func (u *brandRepo) Create(ctx context.Context, req *models.BrandCreate) (*models.Brand, error) {
	id := uuid.New().String()
	query := `
		INSERT INTO "brand" (
			id,
			name,
			brand_image,
			created_at
		)
		VALUES ($1, $2, $3, CURRENT_TIMESTAMP)
		RETURNING id, name, brand_image, created_at, updated_at, deleted_at
	`

	var (
		idd         sql.NullString
		name        sql.NullString
		brand_image sql.NullString
		created_at  sql.NullString
		updated_at  sql.NullString
		delete_at   sql.NullString
	)

	err := u.db.QueryRow(ctx, query, id, req.Name, req.Brand_image).Scan(
		&idd,
		&name,
		&brand_image,
		&created_at,
		&updated_at,
		&delete_at,
	)
	if err != nil {
		u.log.Error("Error while creating brand: " + err.Error())
		return nil, err
	}

	return &models.Brand{
		Id:          idd.String,
		Name:        name.String,
		Brand_image: brand_image.String,
		CreatedAt:   created_at.String,
		UpdatedAt:   updated_at.String,
		DeletedAt:   delete_at.String,
	}, nil
}

func (u *brandRepo) GetByID(ctx context.Context, req *models.BrandPrimaryKey) (*models.Brand, error) {
	var (
		query       string
		id          sql.NullString
		name        sql.NullString
		brand_image sql.NullString
		created_at  sql.NullString
	)

	query = `
		SELECT 
			id,
			name,
			brand_image,
			TO_CHAR(created_at,'dd/mm/yyyy')
		FROM "brand" 
		WHERE id = $1

	`

	err := u.db.QueryRow(ctx, query, req.Id).Scan(
		&id,
		&name,
		&brand_image,
		&created_at,
	)

	if err != nil && err.Error() != "no rows in result set" {
		u.log.Error("error while scanning data" + err.Error())
		return nil, err
	}

	return &models.Brand{
		Id:          id.String,
		Name:        name.String,
		Brand_image: brand_image.String,
		CreatedAt:   created_at.String,
	}, nil
}

func (u *brandRepo) GetList(ctx context.Context, req *models.BrandGetListRequest) (*models.BrandGetListResponse, error) {
	var (
		resp   = &models.BrandGetListResponse{}
		query  string
		offset = " OFFSET 0"
		limit  = " LIMIT 10"
	)

	query = `
		SELECT
			COUNT(*) OVER(),
			id,
			name,
			brand_image,
			TO_CHAR(created_at,'dd/mm/yyyy')
		FROM "brand" 
		
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
		u.log.Error("error is while getting brand list" + err.Error())
		return nil, err
	}

	for rows.Next() {
		var (
			id          sql.NullString
			name        sql.NullString
			brand_image sql.NullString
			created_at  sql.NullString
		)

		err = rows.Scan(
			&resp.Count,
			&id,
			&name,
			&brand_image,
			&created_at,
		)
		if err != nil {
			u.log.Error("error is while getting user list (scanning data)", logger.Error(err))
			return nil, err
		}

		resp.Brand = append(resp.Brand, &models.Brand{
			Id:          id.String,
			Name:        name.String,
			Brand_image: brand_image.String,
			CreatedAt:   created_at.String,
		})
	}
	return resp, nil
}

func (u *brandRepo) Delete(ctx context.Context, req *models.BrandPrimaryKey) error {

	_, err := u.db.Exec(ctx, `DELETE from brand WHERE id = $1`, req.Id)
	if err != nil {
		u.log.Error("error is while deleting brand", logger.Error(err))
		return err
	}

	return nil
}

func (u *brandRepo) Update(ctx context.Context, req *models.BrandUpdate) (int64, error) {
	var (
		query  string
		params map[string]interface{}
	)

	query = `
		UPDATE
			"brand"
		SET
			name = :name,
			brand_image=:brand_image,
			updated_at = NOW()
		WHERE id = :id
	`

	params = map[string]interface{}{
		"id":          req.Id,
		"name":        req.Name,
		"brand_image": req.Brand_image,
	}

	query, args := helper.ReplaceQueryParams(query, params)
	result, err := u.db.Exec(ctx, query, args...)
	if err != nil {
		u.log.Error("error is while updating brand data", logger.Error(err))
		return 0, err
	}

	return result.RowsAffected(), nil
}

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

type categoryRepo struct {
	db  *pgxpool.Pool
	log logger.LoggerI
}

// NewAdminRepo initializes a new instance of adminRepo
func NewCategoryRepo(db *pgxpool.Pool, log logger.LoggerI) *categoryRepo {
	return &categoryRepo{
		db:  db,
		log: log,
	}
}

// Create inserts a new admin into the database
func (u *categoryRepo) Create(ctx context.Context, req *models.CategoryCreate) (*models.Category, error) {
	id := uuid.New().String()
	query := `
		INSERT INTO "category" (
			id,
			name,
			url,
			created_at
		)
		VALUES ($1, $2, $3, CURRENT_TIMESTAMP)
		RETURNING id, name, url, created_at, updated_at, deleted_at
	`

	var (
		idd        sql.NullString
		name       sql.NullString
		url        sql.NullString
		created_at sql.NullString
		updated_at sql.NullString
		delete_at  sql.NullString
	)

	err := u.db.QueryRow(ctx, query, id, req.Name, req.Url).Scan(
		&idd,
		&name,
		&url,
		&created_at,
		&updated_at,
		&delete_at,
	)
	if err != nil {
		u.log.Error("Error while creating category: " + err.Error())
		return nil, err
	}

	return &models.Category{
		Id:        idd.String,
		Name:      name.String,
		Url:       url.String,
		CreatedAt: created_at.String,
		UpdatedAt: updated_at.String,
		DeletedAt: delete_at.String,
	}, nil
}

func (u *categoryRepo) GetByID(ctx context.Context, req *models.CategoryPrimaryKey) (*models.Category, error) {
	var (
		query      string
		id         sql.NullString
		name       sql.NullString
		url        sql.NullString
		created_at sql.NullString
	)

	query = `
		SELECT 
			id,
			name,
			url,
			TO_CHAR(created_at,'dd/mm/yyyy')
		FROM "category" 
		WHERE id = $1

	`

	err := u.db.QueryRow(ctx, query, req.Id).Scan(
		&id,
		&name,
		&url,
		&created_at,
	)

	if err != nil && err.Error() != "no rows in result set" {
		u.log.Error("error while scanning data" + err.Error())
		return nil, err
	}

	return &models.Category{
		Id:        id.String,
		Name:      name.String,
		Url:       url.String,
		CreatedAt: created_at.String,
	}, nil
}

func (u *categoryRepo) GetList(ctx context.Context, req *models.CategoryGetListRequest) (*models.CategoryGetListResponse, error) {
	var (
		resp   = &models.CategoryGetListResponse{}
		query  string
		offset = " OFFSET 0"
		limit  = " LIMIT 10"
	)

	query = `
		SELECT
			COUNT(*) OVER(),
			id,
			name,
			url,
			TO_CHAR(created_at,'dd/mm/yyyy')
		FROM "category" 
		
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
		u.log.Error("error is while getting category list" + err.Error())
		return nil, err
	}

	for rows.Next() {
		var (
			id         sql.NullString
			name       sql.NullString
			url        sql.NullString
			created_at sql.NullString
		)

		err = rows.Scan(
			&resp.Count,
			&id,
			&name,
			&url,
			&created_at,
		)
		if err != nil {
			u.log.Error("error is while getting user list (scanning data)", logger.Error(err))
			return nil, err
		}

		resp.Category = append(resp.Category, &models.Category{
			Id:        id.String,
			Name:      name.String,
			Url:       url.String,
			CreatedAt: created_at.String,
		})
	}
	return resp, nil
}

func (u *categoryRepo) Delete(ctx context.Context, req *models.CategoryPrimaryKey) error {

	_, err := u.db.Exec(ctx, `DELETE from category WHERE id = $1`, req.Id)
	if err != nil {
		u.log.Error("error is while deleting category", logger.Error(err))
		return err
	}

	return nil
}

func (u *categoryRepo) Update(ctx context.Context, req *models.CategoryUpdate) (int64, error) {
	var (
		query  string
		params map[string]interface{}
	)

	query = `
		UPDATE
			"category"
		SET
			name = :name,
			url  = :url,
			updated_at = NOW()
		WHERE id = :id
	`

	params = map[string]interface{}{
		"id":   req.Id,
		"name": req.Name,
		"url":  req.Url,
	}

	query, args := helper.ReplaceQueryParams(query, params)
	result, err := u.db.Exec(ctx, query, args...)
	if err != nil {
		u.log.Error("error is while updating category data", logger.Error(err))
		return 0, err
	}

	return result.RowsAffected(), nil
}

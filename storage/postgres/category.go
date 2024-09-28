package postgres

import (
	"context"
	"database/sql"
	"e-commerce/models"
	"e-commerce/pkg/helper"
	"e-commerce/pkg/logger"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type categoryRepo struct {
	db  *pgxpool.Pool
	log logger.LoggerI
}

func NewCategoryRepo(db *pgxpool.Pool, log logger.LoggerI) *categoryRepo {
	return &categoryRepo{
		db:  db,
		log: log,
	}
}

func (u *categoryRepo) Create(ctx context.Context, req *models.CategoryCreate) (*models.Category, error) {
	id := uuid.New()

	var parentId sql.NullString
	if req.ParentId != "" {

		parentUUID, err := uuid.Parse(req.ParentId)
		if err != nil {
			u.log.Error("Error parsing parent ID: " + err.Error())
			return nil, err
		}
		parentId = sql.NullString{String: parentUUID.String(), Valid: true}
	} else {
		parentId = sql.NullString{Valid: false}
	}

	query := `
		INSERT INTO "category" (
			id,
			name,
			url,
			parent_id,
			created_at
		)
		VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP)
		RETURNING id, name, url, parent_id, created_at, updated_at
	`

	var (
		idd        sql.NullString
		name       sql.NullString
		url        sql.NullString
		parent     sql.NullString
		created_at sql.NullString
		updated_at sql.NullString
	)

	err := u.db.QueryRow(ctx, query, id, req.Name, req.Url, parentId).Scan(
		&idd,
		&name,
		&url,
		&parent,
		&created_at,
		&updated_at,
	)
	if err != nil {
		u.log.Error("Error while creating category: " + err.Error())
		return nil, err
	}

	return &models.Category{
		Id:        idd.String,
		Name:      name.String,
		Url:       url.String,
		ParentId:  parent.String,
		CreatedAt: created_at.String,
		UpdatedAt: updated_at.String,
	}, nil
}

func (u *categoryRepo) GetByID(ctx context.Context, req *models.CategoryPrimaryKey) (*models.Category, error) {
	var (
		query      string
		id         sql.NullString
		name       sql.NullString
		url        sql.NullString
		parent_id  sql.NullString
		created_at sql.NullString
		updated_at sql.NullString
	)

	query = `
		SELECT 
			id,
			name,
			url,
			parent_id,
			created_at
		FROM "category" 
		WHERE id = $1
	`

	err := u.db.QueryRow(ctx, query, req.Id).Scan(
		&id,
		&name,
		&url,
		&parent_id,
		&created_at,
		&updated_at,
	)

	if err != nil && err.Error() != "no rows in result set" {
		u.log.Error("Error while scanning data: " + err.Error())
		return nil, err
	}

	return &models.Category{
		Id:        id.String,
		Name:      name.String,
		Url:       url.String,
		ParentId:  parent_id.String,
		CreatedAt: created_at.String,
		UpdatedAt: updated_at.String,
	}, nil
}

func (u *categoryRepo) GetList(ctx context.Context, req *models.CategoryGetListRequest) (*models.CategoryGetListResponse, error) {
	var (
		resp   = &models.CategoryGetListResponse{}
		query  string
		offset = " OFFSET 0"
		limit  = " LIMIT 10"
		filter = " WHERE TRUE"
		args   = make([]interface{}, 0)
	)

	query = `
		SELECT
			COUNT(*) OVER(),
			id,
			name,
			url,
			parent_id,
			created_at
		FROM "category"
	`

	if req.Offset > 0 {
		offset = fmt.Sprintf(" OFFSET %d", req.Offset)
	}

	if req.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", req.Limit)
	}

	query += filter + limit + offset

	rows, err := u.db.Query(ctx, query, args...)
	if err != nil {
		u.log.Error("Error while getting category list: " + err.Error())
		return nil, err
	}
	defer rows.Close()

	var categories []*models.Category

	for rows.Next() {
		var (
			id         sql.NullString
			name       sql.NullString
			url        sql.NullString
			parent_id  sql.NullString
			created_at sql.NullString
		)

		err = rows.Scan(
			&resp.Count,
			&id,
			&name,
			&url,
			&parent_id,
			&created_at,
		)
		if err != nil {
			u.log.Error("Error while getting category list (scanning data): " + err.Error())
			return nil, err
		}

		category := &models.Category{
			Id:        id.String,
			Name:      name.String,
			Url:       url.String,
			ParentId:  parent_id.String,
			CreatedAt: created_at.String,
		}
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		u.log.Error("Error while iterating over category rows: " + err.Error())
		return nil, err
	}

	// Sort categories by created_at in descending order
	sort.Slice(categories, func(i, j int) bool {
		timeI, _ := time.Parse(time.RFC3339, categories[i].CreatedAt)
		timeJ, _ := time.Parse(time.RFC3339, categories[j].CreatedAt)
		return timeI.After(timeJ)
	})

	if req.Name != "" {
		var filteredCategories []*models.Category
		for _, category := range categories {
			if strings.Contains(strings.ToLower(category.Name), strings.ToLower(req.Name)) {
				filteredCategories = append(filteredCategories, category)
			}
		}
		categories = filteredCategories
	}

	resp.Category = categories

	return resp, nil
}

func (u *categoryRepo) Delete(ctx context.Context, req *models.CategoryPrimaryKey) error {
	_, err := u.db.Exec(ctx, `DELETE FROM "category" WHERE id = $1`, req.Id)
	if err != nil {
		u.log.Error("Error while deleting category: " + err.Error())
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
			url = :url,
			parent_id = :parent_id,
			updated_at = NOW()
		WHERE id = :id
	`

	params = map[string]interface{}{
		"id":   req.Id,
		"name": req.Name,
	}

	query, args := helper.ReplaceQueryParams(query, params)
	result, err := u.db.Exec(ctx, query, args...)
	if err != nil {
		u.log.Error("Error while updating category data: " + err.Error())
		return 0, err
	}

	return result.RowsAffected(), nil
}

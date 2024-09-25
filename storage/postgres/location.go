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

type locationRepo struct {
	db  *pgxpool.Pool
	log logger.LoggerI
}

func NewLocationRepo(db *pgxpool.Pool, log logger.LoggerI) *locationRepo {
	return &locationRepo{
		db:  db,
		log: log,
	}
}

// Create inserts a new admin into the database
func (u *locationRepo) Create(ctx context.Context, req *models.LocationCreate) (*models.Location, error) {

	id := uuid.New().String()
	query := `
		INSERT INTO "location"(
			id,
			name,
			info,
			latitude,
			longitude,
			image,
			created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP)
		RETURNING id, name, info, latitude, longitude, image, created_at, updated_at
	`

	var (
		idd        sql.NullString
		name       sql.NullString
		info       sql.NullString
		latitude   sql.NullFloat64
		longitude  sql.NullFloat64
		image      sql.NullString
		created_at sql.NullString
		updated_at sql.NullString
	)

	err := u.db.QueryRow(ctx, query, id, req.Name, req.Info, req.Latitude, req.Longitude, req.Image).Scan(
		&idd,
		&name,
		&info,
		&latitude,
		&longitude,
		&image,
		&created_at,
		&updated_at,
	)
	if err != nil {
		u.log.Error("Error while creating location: " + err.Error())
		return nil, err
	}

	return &models.Location{
		Id:        idd.String,
		Name:      name.String,
		Info:      info.String,
		Latitude:  latitude.Float64,
		Longitude: longitude.Float64,
		Image:     image.String,
		CreatedAt: created_at.String,
		UpdatedAt: updated_at.String,
	}, nil
}

func (u *locationRepo) GetByID(ctx context.Context, req *models.LacationPrimaryKey) (*models.Location, error) {
	var (
		query      string
		id         sql.NullString
		name       sql.NullString
		info       sql.NullString
		latitude   sql.NullFloat64
		longitude  sql.NullFloat64
		image      sql.NullString
		created_at sql.NullString
	)

	query = `
		SELECT 
			id,
			name,
			info,
			latitude,
			longitude,
			image,
			created_at
		FROM "location" 
		WHERE id = $1

	`

	err := u.db.QueryRow(ctx, query, req.Id).Scan(
		&id,
		&name,
		&info,
		&latitude,
		&longitude,
		&image,
		&created_at,
	)

	if err != nil && err.Error() != "no rows in result set" {
		u.log.Error("error while scanning data" + err.Error())
		return nil, err
	}

	return &models.Location{
		Id:        id.String,
		Name:      name.String,
		Info:      info.String,
		Latitude:  latitude.Float64,
		Longitude: longitude.Float64,
		Image:     image.String,
		CreatedAt: created_at.String,
	}, nil
}

func (u *locationRepo) GetList(ctx context.Context, req *models.LocationGetListRequest) (*models.LocationGetListResponse, error) {
	var (
		resp   = &models.LocationGetListResponse{}
		query  string
		offset = " OFFSET 0"
		limit  = " LIMIT 10"
	)

	query = `
		SELECT
			COUNT(*) OVER(),
			id,
			name,
			info,
			latitude,
			longitude,
			image,
			created_at
		FROM "location" 
		
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
		u.log.Error("error is while getting location list" + err.Error())
		return nil, err
	}

	for rows.Next() {
		var (
			id         sql.NullString
			name       sql.NullString
			info       sql.NullString
			latitude   sql.NullFloat64
			longitude  sql.NullFloat64
			image      sql.NullString
			created_at sql.NullString
		)

		err = rows.Scan(
			&resp.Count,
			&id,
			&name,
			&info,
			&latitude,
			&longitude,
			&image,
			&created_at,
		)
		if err != nil {
			u.log.Error("error is while getting user list (scanning data)", logger.Error(err))
			return nil, err
		}

		resp.Location = append(resp.Location, &models.Location{
			Id:        id.String,
			Name:      name.String,
			Info:      info.String,
			Latitude:  latitude.Float64,
			Longitude: longitude.Float64,
			Image:     image.String,
			CreatedAt: created_at.String,
		})
	}
	return resp, nil
}

func (u *locationRepo) Delete(ctx context.Context, req *models.LacationPrimaryKey) error {

	_, err := u.db.Exec(ctx, `DELETE from location WHERE id = $1`, req.Id)
	if err != nil {
		u.log.Error("error is while deleting location", logger.Error(err))
		return err
	}

	return nil
}

func (u *locationRepo) Update(ctx context.Context, req *models.LocationUpdate) (int64, error) {
	var (
		query  string
		params map[string]interface{}
	)

	query = `
		UPDATE
			"location"
		SET
			name = :name,
			info = :info,
			latitude = :latitude,
			longitude = :longitude,
			image = :image,
			updated_at = NOW()
		WHERE id = :id
	`

	params = map[string]interface{}{
		"id":        req.Id,
		"name":      req.Name,
		"info":      req.Info,
		"latitude":  req.Latitude,
		"longitude": req.Longitude,
		"image":     req.Image,
	}

	query, args := helper.ReplaceQueryParams(query, params)
	result, err := u.db.Exec(ctx, query, args...)
	if err != nil {
		u.log.Error("error is while updating location data", logger.Error(err))
		return 0, err
	}

	return result.RowsAffected(), nil
}

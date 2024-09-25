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

type adminRepo struct {
	db  *pgxpool.Pool
	log logger.LoggerI
}

// NewAdminRepo initializes a new instance of adminRepo
func NewAdminRepo(db *pgxpool.Pool, log logger.LoggerI) *adminRepo {
	return &adminRepo{
		db:  db,
		log: log,
	}
}

// Create inserts a new admin into the database
func (u *adminRepo) Create(ctx context.Context, req *models.AdminCreate) (*models.Admin, error) {

	if !helper.IsValidPhone(req.Phone_number) {
		u.log.Error("Invalid phone number format")
		return nil, fmt.Errorf("invalid phone number format")
	}

	if !helper.IsValidEmail(req.Email) {
		u.log.Error("Invalid email format")
		return nil, fmt.Errorf("invalid email format")
	}

	id := uuid.New().String()
	query := `
		INSERT INTO "admin" (
			id,
			name,
			phone_number,
			email,
			password,
			address,
			created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP)
		RETURNING id, name, phone_number, email, password, address, created_at, updated_at
	`

	var (
		idd          sql.NullString
		name         sql.NullString
		phone_number sql.NullString
		email        sql.NullString
		password     sql.NullString
		address      sql.NullString
		created_at   sql.NullString
		updated_at   sql.NullString
	)

	err := u.db.QueryRow(ctx, query, id, req.Name, req.Phone_number, req.Email, req.Password, req.Address).Scan(
		&idd,
		&name,
		&phone_number,
		&email,
		&password,
		&address,
		&created_at,
		&updated_at,
	)
	if err != nil {
		u.log.Error("Error while creating admin: " + err.Error())
		return nil, err
	}

	return &models.Admin{
		Id:           idd.String,
		Name:         name.String,
		Phone_number: phone_number.String,
		Email:        email.String,
		Password:     password.String,
		Address:      address.String,
		CreatedAt:    created_at.String,
		UpdatedAt:    updated_at.String,
	}, nil
}

func (u *adminRepo) GetByID(ctx context.Context, req *models.AdminPrimaryKey) (*models.Admin, error) {
	var (
		query        string
		id           sql.NullString
		name         sql.NullString
		phone_number sql.NullString
		email        sql.NullString
		password     sql.NullString
		address      sql.NullString
		created_at   sql.NullString
	)

	query = `
		SELECT 
			id,
			name,
			phone_number,
			email,
			password,
			address,
			created_at
		FROM "admin" 
		WHERE id = $1

	`

	err := u.db.QueryRow(ctx, query, req.Id).Scan(
		&id,
		&name,
		&phone_number,
		&email,
		&password,
		&address,
		&created_at,
	)

	if err != nil && err.Error() != "no rows in result set" {
		u.log.Error("error while scanning data" + err.Error())
		return nil, err
	}

	return &models.Admin{
		Id:           id.String,
		Name:         name.String,
		Phone_number: phone_number.String,
		Email:        email.String,
		Password:     password.String,
		Address:      address.String,
		CreatedAt:    created_at.String,
	}, nil
}

func (u *adminRepo) GetList(ctx context.Context, req *models.AdminGetListRequest) (*models.AdminGetListResponse, error) {
	var (
		resp   = &models.AdminGetListResponse{}
		query  string
		offset = " OFFSET 0"
		limit  = " LIMIT 10"
	)

	query = `
		SELECT
			COUNT(*) OVER(),
			id,
			name,
			phone_number,
			email,
			password,
			address,
			created_at
		FROM "admin" 
		
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
		u.log.Error("error is while getting admin list" + err.Error())
		return nil, err
	}

	for rows.Next() {
		var (
			id           sql.NullString
			name         sql.NullString
			phone_number sql.NullString
			email        sql.NullString
			password     sql.NullString
			address      sql.NullString
			created_at   sql.NullString
		)

		err = rows.Scan(
			&resp.Count,
			&id,
			&name,
			&phone_number,
			&email,
			&password,
			&address,
			&created_at,
		)
		if err != nil {
			u.log.Error("error is while getting user list (scanning data)", logger.Error(err))
			return nil, err
		}

		resp.Admin = append(resp.Admin, &models.Admin{
			Id:           id.String,
			Name:         name.String,
			Phone_number: phone_number.String,
			Email:        email.String,
			Password:     password.String,
			Address:      address.String,
			CreatedAt:    created_at.String,
		})
	}
	return resp, nil
}

func (u *adminRepo) Delete(ctx context.Context, req *models.AdminPrimaryKey) error {

	_, err := u.db.Exec(ctx, `DELETE from admin WHERE id = $1`, req.Id)
	if err != nil {
		u.log.Error("error is while deleting admin", logger.Error(err))
		return err
	}

	return nil
}

func (u *adminRepo) Update(ctx context.Context, req *models.AdminUpdate) (int64, error) {
	var (
		query  string
		params map[string]interface{}
	)

	if !helper.IsValidPhone(req.Phone_number) {
		u.log.Error("Invalid phone number format")
		return 0, fmt.Errorf("invalid phone number format")
	}

	if !helper.IsValidEmail(req.Email) {
		u.log.Error("Invalid email format")
		return 0, fmt.Errorf("invalid email format")
	}

	query = `
		UPDATE
			"admin"
		SET
			name = :name,
			phone_number = :phone_number,
			email = :email,
			password = :password,
			address = :address,
			updated_at = NOW()
		WHERE id = :id
	`

	params = map[string]interface{}{
		"id":           req.Id,
		"name":         req.Name,
		"phone_number": req.Phone_number,
		"email":        req.Email,
		"password":     req.Password,
		"address":      req.Address,
	}

	query, args := helper.ReplaceQueryParams(query, params)
	result, err := u.db.Exec(ctx, query, args...)
	if err != nil {
		u.log.Error("error is while updating admin data", logger.Error(err))
		return 0, err
	}

	return result.RowsAffected(), nil
}

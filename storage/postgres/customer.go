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

type customerRepo struct {
	db  *pgxpool.Pool
	log logger.LoggerI
}

func NewCustomerRepo(db *pgxpool.Pool, log logger.LoggerI) *customerRepo {
	return &customerRepo{
		db:  db,
		log: log,
	}
}

func (u *customerRepo) Create(ctx context.Context, req *models.CustomerCreate) (*models.Customer, error) {

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
	INSERT INTO "customer"(
		id,
		name,
		phone_number,
		address,
		email,
		password,
		created_at
)
	VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP)
	RETURNING id, name, phone_number, address, email, password, created_at, updated_at, deleted_at 
	`
	var (
		idd          sql.NullString
		name         sql.NullString
		phone_number sql.NullString
		address      sql.NullString
		email        sql.NullString
		password     sql.NullString
		created_at   sql.NullString
		updated_at   sql.NullString
		delete_at    sql.NullString
	)

	err := u.db.QueryRow(ctx, query, id, req.Name, req.Phone_number, req.Address, req.Email, req.Password).Scan(
		&idd,
		&name,
		&phone_number,
		&address,
		&email,
		&password,
		&created_at,
		&updated_at,
		&delete_at,
	)

	if err != nil {
		u.log.Error("Error while creating customer: " + err.Error())
		return nil, err
	}
	return &models.Customer{
		Id:           idd.String,
		Name:         name.String,
		Phone_number: phone_number.String,
		Address:      address.String,
		Email:        email.String,
		Password:     password.String,
		CreatedAt:    created_at.String,
		UpdatedAt:    updated_at.String,
		DeletedAt:    delete_at.String,
	}, nil
}

func (u *customerRepo) GetByID(ctx context.Context, req *models.CustomerPrimaryKey) (*models.Customer, error) {
	var (
		query        string
		id           sql.NullString
		name         sql.NullString
		phone_number sql.NullString
		address      sql.NullString
		email        sql.NullString
		password     sql.NullString
		created_at   sql.NullString
	)

	query = `
		SELECT 
			id,
			name,
			phone_number,
			address,
			email,
			password,
			TO_CHAR(created_at,'dd/mm/yyyy')
		FROM "customer" 
		WHERE id = $1

	`

	err := u.db.QueryRow(ctx, query, req.Id).Scan(
		&id,
		&name,
		&phone_number,
		&address,
		&email,
		&password,
		&created_at,
	)

	if err != nil && err.Error() != "no rows in result set" {
		u.log.Error("error while scanning data" + err.Error())
		return nil, err
	}

	return &models.Customer{
		Id:           id.String,
		Name:         name.String,
		Phone_number: phone_number.String,
		Address:      address.String,
		Email:        email.String,
		Password:     password.String,
		CreatedAt:    created_at.String,
	}, nil
}

func (u *customerRepo) GetList(ctx context.Context, req *models.CustomerGetListRequest) (*models.CustomerGetListResponse, error) {
	var (
		resp   = &models.CustomerGetListResponse{}
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
			address,
			email,
			password,
			TO_CHAR(created_at,'dd/mm/yyyy')
		FROM "customer" 
		s
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
			address      sql.NullString
			email        sql.NullString
			password     sql.NullString
			created_at   sql.NullString
		)

		err = rows.Scan(
			&resp.Count,
			&id,
			&name,
			&phone_number,
			&address,
			&email,
			&password,
			&created_at,
		)
		if err != nil {
			u.log.Error("error is while getting user list (scanning data)", logger.Error(err))
			return nil, err
		}

		resp.Customer = append(resp.Customer, &models.Customer{
			Id:           id.String,
			Name:         name.String,
			Phone_number: phone_number.String,
			Address:      address.String,
			Email:        email.String,
			Password:     password.String,
			CreatedAt:    created_at.String,
		})
	}
	return resp, nil
}

func (u *customerRepo) Delete(ctx context.Context, req *models.CustomerPrimaryKey) error {

	_, err := u.db.Exec(ctx, `DELETE from customer WHERE id = $1`, req.Id)
	if err != nil {
		u.log.Error("error is while deleting customer", logger.Error(err))
		return err
	}

	return nil
}

func (u *customerRepo) Update(ctx context.Context, req *models.CustomerUpdate) (int64, error) {

	if !helper.IsValidPhone(req.Phone_number) {
		u.log.Error("Invalid phone number format")
		return 0, fmt.Errorf("invalid phone number format")
	}

	if !helper.IsValidEmail(req.Email) {
		u.log.Error("Invalid email format")
		return 0, fmt.Errorf("invalid email format")
	}
	
	var (
		query  string
		params map[string]interface{}
	)

	query = `
		UPDATE
			"customer"
		SET
			name = :name,
			phone_number = :phone_number,
			address = :address,
			email = :email,
			password = :password,
			updated_at = NOW()
		WHERE id = :id
	`

	params = map[string]interface{}{
		"id":           req.Id,
		"name":         req.Name,
		"phone_number": req.Phone_number,
		"address":      req.Address,
		"email":        req.Email,
		"password":     req.Password,
	}

	query, args := helper.ReplaceQueryParams(query, params)
	result, err := u.db.Exec(ctx, query, args...)
	if err != nil {
		u.log.Error("error is while updating customer data", logger.Error(err))
		return 0, err
	}

	return result.RowsAffected(), nil
}

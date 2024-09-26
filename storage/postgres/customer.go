package postgres

import (
	"context"
	"database/sql"
	"e-commerce/models"
	"e-commerce/pkg/helper"
	"e-commerce/pkg/logger"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
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

func (c *customerRepo) GetByLogin(ctx context.Context, login string) (models.Customer, error) {
	var (
		name      sql.NullString
		surname   sql.NullString
		phone     sql.NullString
		birthday  sql.NullString
		gender    sql.NullString
		createdat sql.NullString
		updatedat sql.NullString
	)

	query := `SELECT 
	 id, 
	 name, 
	 surname, 
	 phone_number,
	 birthday,
	 gender,
	 created_at, 
	 updated_at,
	 password
	 FROM "customer" WHERE phone_number= $1 `

	row := c.db.QueryRow(ctx, query, login)

	user := models.Customer{}

	err := row.Scan(
		&user.Id,
		&name,
		&surname,
		&phone,
		&birthday,
		&gender,
		&createdat,
		&updatedat,
		&user.Password,
	)

	if err != nil {
		c.log.Error("failed to scan user by LOGIN from database", logger.Error(err))
		return models.Customer{}, err
	}

	user.Name = name.String
	user.Surname = surname.String // Yangi qo'shilgan ustun
	user.Phone_number = phone.String
	user.Birthday = birthday.String // Yangi qo'shilgan ustun
	user.Gender = gender.String     // Yangi qo'shilgan ustun
	user.CreatedAt = createdat.String
	user.UpdatedAt = updatedat.String

	return user, nil
}

func (c *customerRepo) Login(ctx context.Context, login models.Customer) (string, error) {
	var hashedPass string

	query := `SELECT password
	FROM "customer"
	WHERE phone_number = $1`

	err := c.db.QueryRow(ctx, query,
		login.Phone_number,
	).Scan(&hashedPass)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("incorrect login")
		}
		c.log.Error("failed to get user password from database", logger.Error(err))
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(login.Password))
	if err != nil {
		return "", errors.New("password mismatch")
	}

	return "Logged in successfully", nil
}

func (u *customerRepo) Create(ctx context.Context, req *models.CustomerCreate) (*models.Customer, error) {

	if !helper.IsValidPhone(req.Phone_number) {
		u.log.Error("Invalid phone number format")
		return nil, fmt.Errorf("invalid phone number format")
	}

	// if !helper.IsValidEmail(req.Email) {
	// 	u.log.Error("Invalid email format")
	// 	return nil, fmt.Errorf("invalid email format")
	// }

	id := uuid.New().String()
	query := `
	INSERT INTO "customer"(
		id,
		name,
		surname,
		phone_number,
		birthday,
		gender,
		password,
		created_at
)
	VALUES ($1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP)
	RETURNING id, name, surname, phone_number, birthday, gender, password, created_at, updated_at
	`
	var (
		idd          sql.NullString
		name         sql.NullString
		surname      sql.NullString
		phone_number sql.NullString
		birthday     sql.NullString
		gender       sql.NullString
		password     sql.NullString
		created_at   sql.NullString
		updated_at   sql.NullString
	)

	err := u.db.QueryRow(ctx, query, id, req.Name, req.Surname, req.Phone_number, req.Birthday, req.Gender, req.Password).Scan(
		&idd,
		&name,
		&surname,
		&phone_number,
		&birthday,
		&gender,
		&password,
		&created_at,
		&updated_at,
	)

	if err != nil {
		u.log.Error("Error while creating customer: " + err.Error())
		return nil, err
	}
	return &models.Customer{
		Id:           idd.String,
		Name:         name.String,
		Surname:      surname.String,
		Phone_number: phone_number.String,
		Birthday:     birthday.String,
		Gender:       gender.String,
		Password:     password.String,
		CreatedAt:    created_at.String,
		UpdatedAt:    updated_at.String,
	}, nil
}

func (u *customerRepo) GetByID(ctx context.Context, req *models.CustomerPrimaryKey) (*models.Customer, error) {
	var (
		query        string
		id           sql.NullString
		name         sql.NullString
		surname      sql.NullString
		phone_number sql.NullString
		birthday     sql.NullString
		gender       sql.NullString
		password     sql.NullString
		created_at   sql.NullString
	)

	query = `
		SELECT 
			id,
			name,
			surname,
			phone_number,
			birthday,
			gender,
			password,
			created_at
		FROM "customer" 
		WHERE id = $1

	`

	err := u.db.QueryRow(ctx, query, req.Id).Scan(
		&id,
		&name,
		&surname,
		&phone_number,
		&birthday,
		&gender,
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
		Surname:      surname.String,
		Phone_number: phone_number.String,
		Birthday:     birthday.String,
		Gender:       gender.String,
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
			surname,
			phone_number,
			birthday,
			gender,
			password,
			created_at
		FROM "customer" 
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
			surname      sql.NullString
			phone_number sql.NullString
			birthday     sql.NullString
			gender       sql.NullString
			password     sql.NullString
			created_at   sql.NullString
		)

		err = rows.Scan(
			&resp.Count,
			&id,
			&name,
			&surname,
			&phone_number,
			&birthday,
			&gender,
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
			Surname:      surname.String,
			Phone_number: phone_number.String,
			Birthday:     birthday.String,
			Gender:       gender.String,
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

	var (
		query  string
		params map[string]interface{}
	)

	query = `
		UPDATE
			"customer"
		SET
			name = :name,
			surname = :surname,
			phone_number = :phone_number,
			birthday = :birthday,
			gender = :gender,
			password = :password,
			updated_at = NOW()
		WHERE id = :id
	`

	params = map[string]interface{}{
		"id":           req.Id,
		"name":         req.Name,
		"surname":      req.Surname,
		"phone_number": req.Phone_number,
		"birthday":     req.Birthday,
		"gender":       req.Gender,
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

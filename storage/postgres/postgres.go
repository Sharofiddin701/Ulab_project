package postgres

import (
	"context"
	"e-commerce/config"
	"e-commerce/pkg/logger"
	"e-commerce/storage"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type store struct {
	db       *pgxpool.Pool
	log      logger.LoggerI
	admin    *adminRepo
	customer *customerRepo
	brand    *brandRepo
	category *categoryRepo
	order    *orderRepo
	product  *productRepo
	banner   *bannerRepo
	color    *colorRepo
	location *locationRepo
	// auth     *authRepo
}

func NewConnectionPostgres(cfg *config.Config) (storage.StorageI, error) {
	connect, err := pgxpool.ParseConfig(fmt.Sprintf(
		"host=%s user=%s dbname=%s password=%s port=%d ",
		cfg.PostgresHost,
		cfg.PostgresUser,
		cfg.PostgresDatabase,
		cfg.PostgresPassword,
		cfg.PostgresPort,
	))

	if err != nil {
		return nil, err
	}
	connect.MaxConns = 30

	pgxpool, err := pgxpool.ConnectConfig(context.Background(), connect)
	if err != nil {
		return nil, err
	}
	var loggerLevel = new(string)
	log := logger.NewLogger("app", *loggerLevel)
	defer func() {
		err := logger.Cleanup(log)
		if err != nil {
			return
		}
	}()
	return &store{
		db:  pgxpool,
		log: logger.NewLogger("app", *loggerLevel),
	}, nil
}

func (s *store) Close() {
	s.db.Close()
}

func (s *store) Admin() storage.AdminI {
	if s.admin == nil {
		s.admin = &adminRepo{
			db:  s.db,
			log: s.log,
		}
	}
	return s.admin
}

func (s *store) Customer() storage.CustomerI {
	if s.customer == nil {
		s.customer = &customerRepo{
			db:  s.db,
			log: s.log,
		}
	}
	return s.customer
}

func (s *store) Brand() storage.BrandI {
	if s.brand == nil {
		s.brand = &brandRepo{
			db:  s.db,
			log: s.log,
		}
	}
	return s.brand
}

func (s *store) Category() storage.CategoryI {
	if s.category == nil {
		s.category = &categoryRepo{
			db:  s.db,
			log: s.log,
		}
	}
	return s.category
}

func (s *store) Order() storage.OrderI {
	if s.order == nil {
		s.order = &orderRepo{
			db:  s.db,
			log: s.log,
		}
	}
	return s.order
}

func (s *store) Product() storage.ProductI {
	if s.product == nil {
		s.product = &productRepo{
			db:  s.db,
			log: s.log,
		}
	}
	return s.product
}

func (s *store) Banner() storage.BannerI {
	if s.banner == nil {
		s.banner = &bannerRepo{
			db:  s.db,
			log: s.log,
		}
	}
	return s.banner
}

func (s *store) Color() storage.ColorI {
	if s.color == nil {
		s.color = &colorRepo{
			db:  s.db,
			log: s.log,
		}
	}
	return s.color
}

func (s *store) Location() storage.LocationI {
	if s.location == nil {
		s.location = &locationRepo{
			db:  s.db,
			log: s.log,
		}
	}
	return s.location
}

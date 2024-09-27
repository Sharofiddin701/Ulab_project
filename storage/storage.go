package storage

import (
	"context"
	"e-commerce/models"
	"time"
)

type StorageI interface {
	Close()
	Admin() AdminI
	Redis() RedisI
	Customer() CustomerI
	Brand() BrandI
	Category() CategoryI
	Order() OrderI
	Product() ProductI
	Banner() BannerI
	Color() ColorI
	Location() LocationI
	// Register() AuthRepoI
}

type RedisI interface {
	SetX(ctx context.Context, key string, value interface{}, duration time.Duration) error
	Get(ctx context.Context, key string) (interface{}, error)
	Del(ctx context.Context, key string) error
}

type AdminI interface {
	Create(ctx context.Context, req *models.AdminCreate) (*models.Admin, error)
	GetByID(ctx context.Context, req *models.AdminPrimaryKey) (*models.Admin, error)
	GetList(ctx context.Context, req *models.AdminGetListRequest) (*models.AdminGetListResponse, error)
	Update(ctx context.Context, req *models.AdminUpdate) (int64, error)
	Delete(ctx context.Context, req *models.AdminPrimaryKey) error
}

type CustomerI interface {
	Create(ctx context.Context, req *models.CustomerCreate) (*models.Customer, error)
	GetByID(ctx context.Context, req *models.CustomerPrimaryKey) (*models.Customer, error)
	GetList(ctx context.Context, req *models.CustomerGetListRequest) (*models.CustomerGetListResponse, error)
	Update(ctx context.Context, req *models.CustomerUpdate) (int64, error)
	Delete(ctx context.Context, req *models.CustomerPrimaryKey) error
	// Login(ctx context.Context, login models.Customer) (string, error)
	GetByLogin(ctx context.Context, login string) (models.Customer, error)
	GetByPhoneNumber(ctx context.Context, req string) (models.Customer, error)
}

type BrandI interface {
	Create(ctx context.Context, req *models.BrandCreate) (*models.Brand, error)
	GetByID(ctx context.Context, req *models.BrandPrimaryKey) (*models.Brand, error)
	GetList(ctx context.Context, req *models.BrandGetListRequest) (*models.BrandGetListResponse, error)
	Update(ctx context.Context, req *models.BrandUpdate) (int64, error)
	Delete(ctx context.Context, req *models.BrandPrimaryKey) error
}

type CategoryI interface {
	Create(ctx context.Context, req *models.CategoryCreate) (*models.Category, error)
	GetByID(ctx context.Context, req *models.CategoryPrimaryKey) (*models.Category, error)
	GetList(ctx context.Context, req *models.CategoryGetListRequest) (*models.CategoryGetListResponse, error)
	Update(ctx context.Context, req *models.CategoryUpdate) (int64, error)
	Delete(ctx context.Context, req *models.CategoryPrimaryKey) error
}

// type OrderI interface {
// 	Create(ctx context.Context, req *models.OrderCreate) (*models.Order, error)
// 	GetByID(ctx context.Context, req *models.OrderPrimaryKey) (*models.Order, error)
// 	GetList(ctx context.Context, req *models.OrderGetListRequest) (*models.OrderGetListResponse, error)
// 	Update(ctx context.Context, req *models.OrderUpdate) (int64, error)
// 	Delete(ctx context.Context, req *models.OrderPrimaryKey) error
// }

type OrderI interface {
	CreateOrder(request *models.OrderCreateRequest) (*models.OrderCreateRequest, error)
	GetOrder(orderId string) (*models.Order, error)
	GetAll(ctx context.Context, request *models.OrderGetListRequest) (*[]models.OrderCreateRequest, error)
	UpdateOrder(order models.Order) error
	DeleteOrder(orderId string) error
}

type ProductI interface {
	Create(ctx context.Context, req *models.ProductCreate) (*models.Product, error)
	GetByID(ctx context.Context, req *models.ProductPrimaryKey) (*models.Product, error)
	GetList(ctx context.Context, req *models.ProductGetListRequest) (*models.ProductGetListResponse, error)
	Update(ctx context.Context, req *models.ProductUpdate) (int64, error)
	Delete(ctx context.Context, req *models.ProductPrimaryKey) error
}

type BannerI interface {
	Create(ctx context.Context, req *models.BannerCreate) (*models.Banner, error)
	GetByID(ctx context.Context, req *models.BannerPrimaryKey) (*models.Banner, error)
	GetList(ctx context.Context, req *models.BannerGetListRequest) (*models.BannerGetListResponse, error)
	Update(ctx context.Context, req *models.BannerUpdate) (int64, error)
	Delete(ctx context.Context, req *models.BannerPrimaryKey) error
}

type ColorI interface {
	Create(ctx context.Context, req *models.ColorCreate) (*models.Color, error)
	GetList(ctx context.Context, req *models.ColorGetListRequest) (*models.ColorGetListResponse, error)
	Delete(ctx context.Context, req *models.ColorPrimaryKey) error
}

type LocationI interface {
	Create(ctx context.Context, req *models.LocationCreate) (*models.Location, error)
	GetByID(ctx context.Context, req *models.LacationPrimaryKey) (*models.Location, error)
	GetList(ctx context.Context, req *models.LocationGetListRequest) (*models.LocationGetListResponse, error)
	Update(ctx context.Context, req *models.LocationUpdate) (int64, error)
	Delete(ctx context.Context, req *models.LacationPrimaryKey) error
}

// type AuthRepoI interface {
// 	VerifyCode(ctx context.Context, req *models.RegisterRequest) (string, error)
// 	DeleteVerifiedCode(ctx context.Context, req *models.RegisterRequest) error
// 	Create(ctx context.Context, reg *models.Register) error
// }

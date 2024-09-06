package storage

import (
	"context"
	"e-commerce/models"
)

type StorageI interface {
	Close()
	Admin() AdminI
	Customer() CustomerI
	Brand() BrandI
	Category() CategoryI
	Order() OrderI
	Product() ProductI
	OrderProduct() OrderProductI
	Banner() BannerI
	Color() ColorI
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

type OrderI interface {
	Create(ctx context.Context, req *models.OrderCreate) (*models.Order, error)
	GetByID(ctx context.Context, req *models.OrderPrimaryKey) (*models.Order, error)
	GetList(ctx context.Context, req *models.OrderGetListRequest) (*models.OrderGetListResponse, error)
	Update(ctx context.Context, req *models.OrderUpdate) (int64, error)
	Delete(ctx context.Context, req *models.OrderPrimaryKey) error
}

type ProductI interface {
	Create(ctx context.Context, req *models.ProductCreate) (*models.Product, error)
	GetByID(ctx context.Context, req *models.ProductPrimaryKey) (*models.Product, error)
	GetList(ctx context.Context, req *models.ProductGetListRequest) (*models.ProductGetListResponse, error)
	Update(ctx context.Context, req *models.ProductUpdate) (int64, error)
	Delete(ctx context.Context, req *models.ProductPrimaryKey) error
}

type OrderProductI interface {
	Create(ctx context.Context, req *models.OrderProductCreate) (*models.OrderProduct, error)
	GetByID(ctx context.Context, req *models.OrderProductPrimaryKey) (*models.OrderProduct, error)
	GetList(ctx context.Context, req *models.OrderProductGetListRequest) (*models.OrderProductGetListResponse, error)
	Update(ctx context.Context, req *models.OrderProductUpdate) (int64, error)
	Delete(ctx context.Context, req *models.OrderProductPrimaryKey) error
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

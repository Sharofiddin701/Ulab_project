package api

import (
	"e-commerce/api/handler"
	"e-commerce/config"
	"e-commerce/pkg/logger"
	"e-commerce/service"
	"e-commerce/storage"

	_ "e-commerce/api/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewApi(r *gin.Engine, cfg *config.Config, storage storage.StorageI, logger logger.LoggerI, service service.IServiceManager) {
	h := handler.NewHandler(cfg, storage, logger, service)
	r.Use(customCORSMiddleware())
	v1 := r.Group("/e_commerce/api/v1")

	v1.POST("/login", h.UserLogin)
	v1.POST("/sendcode", h.UserRegister)
	v1.POST("/verifycode", h.UserRegisterConfirm)
	v1.POST("/byphoneconfirm", h.UserLoginByPhoneConfirm)

	v1.POST("/admin", h.CreateAdmin)
	v1.GET("/admin/:id", h.GetByIdAdmin)
	v1.GET("/admin", h.GetListAdmin)
	v1.PUT("/admin/:id", h.UpdateAdmin)
	v1.DELETE("/admin/:id", h.DeleteAdmin)

	v1.POST("/color", h.CreateColor)
	v1.GET("/color", h.GetListColor)
	v1.DELETE("/color/:id", h.DeleteColor)

	v1.POST("/banner", h.CreateBanner)
	v1.GET("/banner/:id", h.GetByIdBanner)
	v1.GET("/banner", h.GetListBanner)
	v1.PUT("/banner/:id", h.UpdateBanner)
	v1.DELETE("/banner/:id", h.DeleteBanner)

	v1.POST("/customer", h.CreateCustomer)
	v1.GET("/customer/:id", h.GetByIdCustomer)
	v1.GET("/customer", h.GetListCustomer)
	v1.PUT("/customer/:id", h.UpdateCustomer)
	v1.DELETE("/customer/:id", h.DeleteCustomer)

	v1.POST("/brand", h.CreateBrand)
	v1.GET("/brand/:id", h.GetByIdBrand)
	v1.GET("/brand", h.GetListBrand)
	v1.PUT("/brand/:id", h.UpdateBrand)
	v1.DELETE("/brand/:id", h.DeleteBrand)

	v1.POST("/category", h.CreateCategory)
	v1.GET("/category/:id", h.GetByIdCategory)
	v1.GET("/category", h.GetListCategory)
	v1.PUT("/category/:id", h.UpdateCategory)
	v1.DELETE("/category/:id", h.DeleteCategory)

	v1.POST("/order", h.CreateOrder)
	v1.GET("/order/:id", h.GetByIdOrder)
	v1.GET("/order", h.GetAllOrders)
	v1.PUT("/order/:id", h.UpdateOrder)
	v1.DELETE("/order/:id", h.DeleteOrder)

	v1.POST("/product", h.CreateProduct)
	v1.GET("/product/:id", h.GetByIdProduct)
	v1.GET("/product", h.GetListProduct)
	v1.PUT("/product/:id", h.UpdateProduct)
	v1.DELETE("/product/:id", h.DeleteProduct)

	v1.POST("upload-files", h.UploadFiles)
	v1.DELETE("delete-file", h.DeleteFile)

	v1.POST("/location", h.CreateLocation)
	v1.GET("/location/:id", h.GetByIdLocation)
	v1.GET("/location", h.GetListLocation)
	v1.PUT("/location/:id", h.UpdateLocation)
	v1.DELETE("/location/:id", h.DeleteLocation)

	url := ginSwagger.URL("swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
}

func customCORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE, HEAD")
		c.Header("Access-Control-Allow-Headers", "Platform-Id, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Max-Age", "3600")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

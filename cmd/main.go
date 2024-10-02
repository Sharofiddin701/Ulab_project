package main

import (
	"e-commerce/api"
	"e-commerce/config"
	"e-commerce/service"
	"fmt"
	"net/http"
	"time"

	postgres "e-commerce/storage/postgres"
	"e-commerce/storage/redis"

	"github.com/gin-gonic/gin"
	"github.com/saidamir98/udevs_pkg/logger"
)

// Yangi qo'shilgan keepAlive funksiyasi
func keepAlive(cfg *config.Config) {
	for {
		_, err := http.Get(fmt.Sprintf("http://localhost%s/ping", cfg.HTTPPort))
		if err != nil {
			fmt.Println("Ping yuborishda xatolik:", err)
		} else {
			fmt.Println("Ping muvaffaqiyatli yuborildi")
		}
		time.Sleep(1 * time.Minute)
	}
}

func main() {
	cfg := config.Load()

	var loggerLevel = new(string)

	*loggerLevel = logger.LevelDebug

	switch cfg.Environment {
	case config.DebugMode:
		*loggerLevel = logger.LevelDebug
		gin.SetMode(gin.DebugMode)
	case config.TestMode:
		*loggerLevel = logger.LevelDebug
		gin.SetMode(gin.TestMode)
	default:
		*loggerLevel = logger.LevelInfo
		gin.SetMode(gin.ReleaseMode)
	}

	log := logger.NewLogger("app", *loggerLevel)
	defer func() {
		err := logger.Cleanup(log)
		if err != nil {
			return
		}
	}()

	pgconn, err := postgres.NewConnectionPostgres(&cfg)
	if err != nil {
		panic("postgres no connection: " + err.Error())
	}
	defer pgconn.Close()

	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger())

	newRedis := redis.New(cfg)
	services := service.New(pgconn, log, newRedis)

	api.NewApi(r, &cfg, pgconn, log, services)

	// Yangi qo'shilgan: Keep-alive funksiyasini ishga tushirish
	go keepAlive(&cfg)

	// Yangi qo'shilgan: Ping endpointini qo'shish
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	fmt.Println("Listening server", cfg.PostgresHost+cfg.HTTPPort)
	err = r.Run(cfg.HTTPPort)
	if err != nil {
		panic(err)
	}
}

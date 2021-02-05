package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"local/auth-svc/auth"
	handlers "local/auth-svc/handler"
	"local/auth-svc/middleware"
	accounts "local/auth-svc/services"
	"local/auth-svc/sqldb"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func setEncryptionKeyEnv() {
	bytes := make([]byte, 32) //generate a random 32 byte key for AES-256
	if _, err := rand.Read(bytes); err != nil {
		panic(err.Error())
	}

	key := hex.EncodeToString(bytes) //encode key in bytes to string
	os.Setenv("AESKEY", key)
	fmt.Println("key :", os.Getenv("AESKEY"))
}

func NewRedisDB(host, port, password string) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       0,
	})
	return redisClient
}

func main() {
	// setEncryptionKeyEnv()
	sqldb.InitializeDatabase()
	appAddr := ":" + os.Getenv("PORT")

	//redis details
	redis_host := os.Getenv("REDIS_HOST")
	redis_port := os.Getenv("REDIS_PORT")
	redis_password := os.Getenv("REDIS_PASSWORD")

	redisClient := NewRedisDB(redis_host, redis_port, redis_password)

	var rd = auth.NewAuth(redisClient)
	var tk = auth.NewToken()
	var service = handlers.NewProfile(rd, tk)

	var router = gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "HEAD", "OPTIONS", "GET", "PUT"},
		AllowHeaders:     []string{"Access-Control-Allow-Headers", "Origin", "Accept", "X-Requested-With", "Content-Type", "Authorization", "Access-Control-Request-Method", "Access-Control-Request-Headers"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.POST("/api/v1/login", service.SendLoginCookie)
	router.POST("/api/v1/logout", service.LogoutSession)
	router.POST("/api/v1/refresh", service.RefreshSession)
	router.POST("/api/v1/account/new", middleware.TokenAuthMiddleware(), accounts.AddAccount)
	router.POST("/api/v1/account/update", middleware.TokenAuthMiddleware(), accounts.ModifyAccount)
	router.POST("/api/v1/account/remove", middleware.TokenAuthMiddleware(), accounts.RemoveAccount)
	router.POST("/api/v1/account/list", middleware.TokenAuthMiddleware(), accounts.ListAccounts)
	router.POST("/api/v1/account/find", middleware.TokenAuthMiddleware(), accounts.FindUser)
	srv := &http.Server{
		Addr:    appAddr,
		Handler: router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	//Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}

package handler

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

var (
	ctx = context.Background()
	rdb *redis.Client
)

// Twilio credentials
const (
	AccountSID        = "ACc39725a654edb7264186eca22f221a47"
	AuthToken         = "0f95c1f4de45c8e8fed7a68302acea7d"
	TwilioPhoneNumber = "+13026045203"
	RedisAddr         = "localhost:6379" // Redis server address
)

// User model
type User struct {
	PhoneNumber string `json:"phone_number"`
	Code        string `json:"code"`
}

func main() {
	// Initialize Redis client
	rdb = redis.NewClient(&redis.Options{
		Addr: RedisAddr,
	})

	// Check Redis connection
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	// Initialize Gin router
	router := gin.Default()

	// Define routes
	router.POST("/register", registerUserHandler)
	router.GET("/send-code", sendCodeHandler)
	router.GET("/verify-code", verifyCodeHandler)

	// Start the server
	log.Println("Server started at :8080")
	router.Run(":8080")
}

func registerUserHandler(c *gin.Context) {
	var user User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	// Check if phone number already exists (in Redis)
	_, err := rdb.Get(ctx, user.PhoneNumber).Result()
	if err != redis.Nil {
		c.JSON(400, gin.H{"error": "Phone number already registered"})
		return
	}

	// Save user data (in a real app, save it in a database)
	err = rdb.Set(ctx, user.PhoneNumber, "", 0).Err()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(200, gin.H{"message": "User registered successfully"})
}

func sendCodeHandler(c *gin.Context) {
	phone := c.Query("phone")

	if phone == "" {
		c.JSON(400, gin.H{"error": "Phone number is required"})
		return
	}

	// Generate a random 6-digit code
	code := generateCode()

	// Save the code in Redis with an expiration time (e.g., 10 minutes)
	err := rdb.Set(ctx, phone, code, 10*time.Minute).Err()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to save code in Redis"})
		return
	}

	// Send the code via Twilio
	err = sendSMS(phone, code)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to send SMS"})
		return
	}

	c.JSON(200, gin.H{"message": fmt.Sprintf("Verification code sent to %s", phone)})
}

func verifyCodeHandler(c *gin.Context) {
	phone := c.Query("phone")
	code := c.Query("code")

	if phone == "" || code == "" {
		c.JSON(400, gin.H{"error": "Phone number and code are required"})
		return
	}

	// Retrieve the code from Redis
	storedCode, err := rdb.Get(ctx, phone).Result()
	if err == redis.Nil {
		c.JSON(404, gin.H{"error": "Phone number not found or code expired"})
		return
	} else if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve code from Redis"})
		return
	}

	if storedCode == code {
		c.JSON(200, gin.H{"message": fmt.Sprintf("Phone number %s verified successfully!", phone)})
	} else {
		c.JSON(401, gin.H{"error": "Invalid code"})
	}
}

func generateCode() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func sendSMS(to string, code string) error {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: AccountSID,
		Password: AuthToken,
	})

	params := &openapi.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom(TwilioPhoneNumber)
	params.SetBody(fmt.Sprintf("Your verification code is %s", code))

	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		log.Printf("Failed to send SMS: %v", err)
		return err
	}

	log.Printf("Message sent successfully: SID %s", *resp.Sid)
	return nil
}

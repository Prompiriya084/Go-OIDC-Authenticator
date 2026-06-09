package dataaccess

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func InitRedis(
	ctx context.Context,
) (*redis.Client, error) {
	err := godotenv.Load("../../../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	address := os.Getenv("REDIS_ADDRESS")
	password := os.Getenv("REDIS_PASSWORD")
	redisClient := redis.NewClient(&redis.Options{
		Addr:     address,  // Address ของ Redis Server
		Password: password, // ใส่รหัสผ่านหากมี
		DB:       0,        // เลือก Database Index (ปกติเป็น 0)
		PoolSize: 10,       // กำหนด Connection Pool สำหรับงาน Production
	})

	// ทดสอบการเชื่อมต่อ
	err = redisClient.Ping(ctx).Err()
	if err != nil {
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}

	log.Println("Connect Redis successfully.")

	return redisClient, nil
}

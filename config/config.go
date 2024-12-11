package config

import (
    "os"
    "github.com/joho/godotenv"
)

type Config struct {
    MongoURI        string
    DatabaseName    string
    JWTSecret       string
    AccountID       string // Thêm field cho Cloudflare account ID
    R2AccessKeyID   string
    R2AccessKeySecret string
    R2BucketName    string
    R2PublicURL     string // Đổi tên từ PublicURL thành R2PublicURL
}

func LoadConfig() (*Config, error) {
    err := godotenv.Load()
    if err != nil {
        return nil, err
    }

    return &Config{
        MongoURI:         os.Getenv("MONGO_URI"),
        DatabaseName:     os.Getenv("DATABASE_NAME"),
        JWTSecret:        os.Getenv("JWT_SECRET"),
        AccountID:        os.Getenv("R2_ACCOUNT_ID"),
        R2AccessKeyID:    os.Getenv("R2_ACCESS_KEY"),
        R2AccessKeySecret: os.Getenv("R2_SECRET_KEY"),
        R2BucketName:     os.Getenv("R2_BUCKET"),
        R2PublicURL:      os.Getenv("R2_PUBLIC_URL"),
    }, nil
}
package main

import (
    "context"
    "mime/multipart"
    "go-blog-backend/config"
    "go-blog-backend/handlers"
    "go-blog-backend/middleware"
    "go-blog-backend/services"
    "go-blog-backend/repositories"
    "go-blog-backend/pkg/cloudflare"
    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "log"
)

// UploadServiceAdapter adapts UploadService to handlers.UploadService interface
type UploadServiceAdapter struct {
    Service *services.UploadService
}

func (a *UploadServiceAdapter) UploadImage(file *multipart.FileHeader) (*handlers.FileUpload, error) {
    result, err := a.Service.UploadImage(file)
    if err != nil {
        return nil, err
    }
    return &handlers.FileUpload{URL: result.URL}, nil
}

func (a *UploadServiceAdapter) DeleteImage(filename string) error {
    return a.Service.DeleteImage(filename)
}

func main() {
    // Load configuration
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatal("Cannot load config:", err)
    }

    // Setup MongoDB connection
    mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.MongoURI))
    if err != nil {
        log.Fatal("Cannot connect to MongoDB:", err)
    }
    defer mongoClient.Disconnect(context.Background())

    // Test MongoDB connection
    err = mongoClient.Ping(context.Background(), nil)
    if err != nil {
        log.Fatal("Cannot ping MongoDB:", err)
    }
    if err == nil {
        log.Println("Successfully pinged MongoDB")
    }

    db := mongoClient.Database(cfg.DatabaseName)

    // Setup R2 client
    r2Client, err := cloudflare.NewR2Client(cloudflare.R2Config{
        AccountID:       cfg.AccountID,        // Thêm AccountID từ config
        AccessKeyID:     cfg.R2AccessKeyID,
        AccessKeySecret: cfg.R2AccessKeySecret,
        BucketName:      cfg.R2BucketName,
        PublicURL:       cfg.R2PublicURL,      // Sửa thành R2PublicURL từ config
    })
    if err != nil {
        log.Fatal("Cannot create R2 client:", err)
    }
    if err == nil {
        log.Println("Successfully created R2 client")
    }

    // Setup repositories
    userRepo := repositories.NewUserRepository(db)
    postRepo := repositories.NewPostRepository(db)

    // Setup services
    userService := services.NewUserService(userRepo, cfg.JWTSecret)
    postService := services.NewPostService(postRepo)
    uploadService := services.NewUploadService(r2Client)

    // Setup handlers
    userHandler := handlers.NewUserHandler(userService)
    postHandler := handlers.NewPostHandler(postService)
    uploadHandler := handlers.NewUploadHandler(&UploadServiceAdapter{
        Service: uploadService,
    })

    // Setup Gin router
    r := gin.Default()

    // CORS middleware
    r.Use(func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Authorization, Content-Type")
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        c.Next()
    })

    // API routes
    api := r.Group("/api")
    {
        // Public routes
        api.POST("/register", userHandler.Register)
        api.POST("/login", userHandler.Login)
        api.GET("/posts", postHandler.List)
        api.GET("/posts/:id", postHandler.Get)

        // Protected routes
        protected := api.Group("/")
        protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
        {
            // User routes
            protected.PUT("/user", userHandler.Update)
            protected.DELETE("/user", userHandler.Delete)

            // Post routes
            protected.POST("/posts", postHandler.Create)
            protected.PUT("/posts/:id", postHandler.Update)
            protected.DELETE("/posts/:id", postHandler.Delete)

            // Upload routes
            protected.POST("/upload", uploadHandler.UploadImage)
        }
    }

    // Start server
    log.Println("Server starting on :8080")
    if err := r.Run(":8080"); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}
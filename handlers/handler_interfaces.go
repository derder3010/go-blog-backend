package handlers

import "go-blog-backend/models"

type Response struct {
    Status  string      `json:"status"`
    Message string      `json:"message,omitempty"`
    Data    interface{} `json:"data,omitempty"`
}

type UserService interface {
    Register(username, email, password string) (*models.User, error)
    Login(email, password string) (string, error) // returns JWT token
    Update(userID string, updates map[string]interface{}) error
    Delete(userID string) error
    GetByID(userID string) (*models.User, error)
}

type PostService interface {
    Create(post *models.Post) error
    Update(postID string, updates map[string]interface{}) error
    Delete(postID string) error
    Get(postID string) (*models.Post, error)
    List(page, limit int) ([]*models.Post, error)
}
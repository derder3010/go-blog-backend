package services

import (
    "go-blog-backend/models"
    "time"
)

type PostRepository interface {
    Create(post *models.Post) error
    GetByID(id string) (*models.Post, error)
    Update(id string, updates map[string]interface{}) error
    Delete(id string) error
    List(page, limit int) ([]*models.Post, error)
    GetByAuthor(authorID string) ([]*models.Post, error)
}

type PostService struct {
    repo PostRepository
}

func NewPostService(repo PostRepository) *PostService {
    return &PostService{
        repo: repo,
    }
}

func (s *PostService) Create(post *models.Post) error {
    post.CreatedAt = time.Now()
    post.UpdatedAt = time.Now()
    return s.repo.Create(post)
}

func (s *PostService) Get(postID string) (*models.Post, error) {
    return s.repo.GetByID(postID)
}

func (s *PostService) Update(postID string, updates map[string]interface{}) error {
    updates["updated_at"] = time.Now()
    return s.repo.Update(postID, updates)
}

func (s *PostService) Delete(postID string) error {
    return s.repo.Delete(postID)
}

func (s *PostService) List(page, limit int) ([]*models.Post, error) {
    return s.repo.List(page, limit)
}
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

// NewPostService returns a new PostService instance, given a PostRepository.
func NewPostService(repo PostRepository) *PostService {
    return &PostService{
        repo: repo,
    }
}

// Create creates a new post in the "posts" collection in the MongoDB database.
//
// The created_at and updated_at fields are automatically set to the current time.
//
// The returned error will be non-nil if any error occurred during the create
// process.
func (s *PostService) Create(post *models.Post) error {
    post.CreatedAt = time.Now()
    post.UpdatedAt = time.Now()
    return s.repo.Create(post)
}

// Get returns a post by the given ID.
//
// The returned error will be non-nil if any error occurred during the get
// process.
func (s *PostService) Get(postID string) (*models.Post, error) {
    return s.repo.GetByID(postID)
}

// Update updates the fields of the post with the given ID in the "posts" collection.
//
// The updates parameter is a map of key-value pairs where the key is the field name
// and the value is the new value for that field. The updated_at field is automatically
// set to the current time.
//
// The returned error will be non-nil if any error occurred during the update process.
func (s *PostService) Update(postID string, updates map[string]interface{}) error {
    updates["updated_at"] = time.Now()
    return s.repo.Update(postID, updates)
}

// Delete deletes the post with the given ID from the "posts" collection in the
// MongoDB database.
//
// The returned error will be non-nil if any error occurred during the delete
// process.
func (s *PostService) Delete(postID string) error {
    return s.repo.Delete(postID)
}

// List returns a slice of posts, sorted by created_at in descending order,
// limited to the given number of items, and starting from the given page.
//
// The page and limit parameters are 1-indexed, so the first page should have
// page = 1 and limit = 10 to get the first 10 results.
//
// The returned error will be non-nil if any error occurred during the find
// process.
func (s *PostService) List(page, limit int) ([]*models.Post, error) {
    return s.repo.List(page, limit)
}
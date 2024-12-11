package handlers

import (
    "github.com/gin-gonic/gin"
    "net/http"
    "go-blog-backend/models"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "strconv"
)

type PostHandler struct {
    postService PostService
}

// NewPostHandler returns a new PostHandler instance, given a PostService.
func NewPostHandler(postService PostService) *PostHandler {
    return &PostHandler{
        postService: postService,
    }
}

type CreatePostRequest struct {
    Title    string `json:"title" binding:"required"`
    Content  string `json:"content" binding:"required"`
    ImageURL string `json:"image_url,omitempty"`
}

// Create creates a new post in the "posts" collection in the MongoDB database.
//
// The request body should contain a JSON object with the following fields:
//   - title: The title of the post.
//   - content: The content of the post.
//   - image_url: An optional URL to an image associated with the post.
//
// The response will be a JSON object with the following fields:
//   - status: The status of the request. Will be "success" on success, or "error" on error.
//   - message: A human-readable message describing the result of the request.
//   - data: The newly created Post instance, or nil if an error occurred.
func (h *PostHandler) Create(c *gin.Context) {
    var req CreatePostRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, Response{
            Status:  "error",
            Message: "Invalid request data",
        })
        return
    }

    userID, _ := c.Get("user_id")
    post := &models.Post{
        Title:    req.Title,
        Content:  req.Content,
        ImageURL: req.ImageURL,
        AuthorID: userID.(primitive.ObjectID),
    }

    if err := h.postService.Create(post); err != nil {
        c.JSON(http.StatusInternalServerError, Response{
            Status:  "error",
            Message: "Failed to create post",
        })
        return
    }

    c.JSON(http.StatusCreated, Response{
        Status: "success",
        Data:   post,
    })
}

// Get retrieves a post by its ID from the "posts" collection.
//
// The ID should be provided as a URL parameter.
//
// The response will be a JSON object with the following fields:
//   - status: The status of the request. Will be "success" on success, or "error" on error.
//   - message: A human-readable message describing the result of the request, if an error occurs.
//   - data: The requested Post instance on success, or nil if not found.
func (h *PostHandler) Get(c *gin.Context) {
    postID := c.Param("id")

    post, err := h.postService.Get(postID)
    if err != nil {
        c.JSON(http.StatusNotFound, Response{
            Status:  "error",
            Message: "Post not found",
        })
        return
    }

    c.JSON(http.StatusOK, Response{
        Status: "success",
        Data:   post,
    })
}

type UpdatePostRequest struct {
    Title    string `json:"title,omitempty"`
    Content  string `json:"content,omitempty"`
    ImageURL string `json:"image_url,omitempty"`
}

// Update updates the fields of the post with the given ID in the "posts" collection.
//
// The request body should contain a JSON object with any of the following fields:
//   - title: The new title for the post.
//   - content: The new content for the post.
//   - image_url: The new image URL for the post.
//
// The response will be a JSON object with the following fields:
//   - status: The status of the request. Will be "success" on success, or "error" on error.
//   - message: A human-readable message describing the result of the request.
func (h *PostHandler) Update(c *gin.Context) {
    postID := c.Param("id")
    _, _ = c.Get("user_id")

    var req UpdatePostRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, Response{
            Status:  "error",
            Message: "Invalid request data",
        })
        return
    }

    updates := make(map[string]interface{})
    if req.Title != "" {
        updates["title"] = req.Title
    }
    if req.Content != "" {
        updates["content"] = req.Content
    }
    if req.ImageURL != "" {
        updates["image_url"] = req.ImageURL
    }

    if err := h.postService.Update(postID, updates); err != nil {
        c.JSON(http.StatusInternalServerError, Response{
            Status:  "error",
            Message: "Failed to update post",
        })
        return
    }

    c.JSON(http.StatusOK, Response{
        Status:  "success",
        Message: "Post updated successfully",
    })
}

// Delete deletes the post with the given ID from the "posts" collection in the
// MongoDB database.
//
// The request body should contain no data.
//
// The response will be a JSON object with the following fields:
//   - status: The status of the request. Will be "success" on success, or "error" on error.
//   - message: A human-readable message describing the result of the request.
func (h *PostHandler) Delete(c *gin.Context) {
    postID := c.Param("id")
    _, _ = c.Get("user_id")

    if err := h.postService.Delete(postID); err != nil {
        c.JSON(http.StatusInternalServerError, Response{
            Status:  "error",
            Message: "Failed to delete post",
        })
        return
    }

    c.JSON(http.StatusOK, Response{
        Status:  "success",
        Message: "Post deleted successfully",
    })
}

// List retrieves a list of posts from the "posts" collection.
//
// The request parameters should include:
//   - page: The page number to retrieve. Defaults to 1 if not specified.
//   - limit: The number of posts per page. Defaults to 10 if not specified.
//
// The response will be a JSON object with the following fields:
//   - status: The status of the request. Will be "success" on success, or "error" on error.
//   - message: A human-readable message describing the result of the request, if an error occurs.
//   - data: A slice of Post instances on success.
func (h *PostHandler) List(c *gin.Context) {
    page := 1
    limit := 10

    if pageStr := c.Query("page"); pageStr != "" {
        if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
            page = p
        }
    }

    if limitStr := c.Query("limit"); limitStr != "" {
        if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
            limit = l
        }
    }

    posts, err := h.postService.List(page, limit)
    if err != nil {
        c.JSON(http.StatusInternalServerError, Response{
            Status:  "error",
            Message: "Failed to fetch posts",
        })
        return
    }

    c.JSON(http.StatusOK, Response{
        Status: "success",
        Data:   posts,
    })
}
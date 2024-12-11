package handlers

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

type UserHandler struct {
    userService UserService
}

// NewUserHandler creates a new UserHandler instance with the provided UserService.
//
// Parameters:
//   - userService: The UserService interface used for interacting with user-related operations.
//
// Returns a pointer to a UserHandler instance.
func NewUserHandler(userService UserService) *UserHandler {
    return &UserHandler{
        userService: userService,
    }
}

type RegisterRequest struct {
    Username string `json:"username" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}

// Register creates a new user in the "users" collection in the MongoDB database.
//
// If the email address is already registered, an error will be returned.
//
// Parameters:
//   - c: The Gin Context object for the current request.
//
// The request body should contain a JSON object with the following fields:
//   - username: The desired username for the new user.
//   - email: The email address for the new user.
//   - password: The desired password for the new user.
//
// The response will be a JSON object with the following fields:
//   - status: The status of the request. Will be "success" on success, or "error" on error.
//   - message: A human-readable message describing the result of the request.
//   - data: The newly created User instance, or nil if an error occurred.
func (h *UserHandler) Register(c *gin.Context) {
    var req RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, Response{
            Status:  "error",
            Message: "Invalid request data",
        })
        return
    }

    user, err := h.userService.Register(req.Username, req.Email, req.Password)
    if err != nil {
        c.JSON(http.StatusInternalServerError, Response{
            Status:  "error",
            Message: "Failed to register user",
        })
        return
    }

    c.JSON(http.StatusCreated, Response{
        Status: "success",
        Data:   user,
    })
}

type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

// Login logs in a user and returns a JWT token in the response.
//
// The request body should contain a JSON object with the following fields:
//   - email: The email address of the user to log in.
//   - password: The password of the user to log in.
//
// The response will be a JSON object with the following fields:
//   - status: The status of the request. Will be "success" on success, or "error" on error.
//   - message: A human-readable message describing the result of the request.
//   - data: An object containing a single field, "token", which is the JWT token for the logged-in user.
func (h *UserHandler) Login(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, Response{
            Status:  "error",
            Message: "Invalid request data",
        })
        return
    }

    token, err := h.userService.Login(req.Email, req.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, Response{
            Status:  "error",
            Message: "Invalid credentials",
        })
        return
    }

    c.JSON(http.StatusOK, Response{
        Status: "success",
        Data: map[string]string{
            "token": token,
        },
    })
}

type UpdateUserRequest struct {
    Username string `json:"username,omitempty"`
    Email    string `json:"email,omitempty" binding:"omitempty,email"`
    Password string `json:"password,omitempty" binding:"omitempty,min=6"`
}

// Update updates the user with the given ID in the "users" collection.
//
// The request body should contain a JSON object with the following fields:
//   - username: The new username for the user, or null if no change is desired.
//   - email: The new email address for the user, or null if no change is desired.
//   - password: The new password for the user, or null if no change is desired.
//
// The response will be a JSON object with the following fields:
//   - status: The status of the request. Will be "success" on success, or "error" on error.
//   - message: A human-readable message describing the result of the request.
func (h *UserHandler) Update(c *gin.Context) {
    userID, _ := c.Get("user_id")
    
    var req UpdateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, Response{
            Status:  "error",
            Message: "Invalid request data",
        })
        return
    }

    updates := make(map[string]interface{})
    if req.Username != "" {
        updates["username"] = req.Username
    }
    if req.Email != "" {
        updates["email"] = req.Email
    }
    if req.Password != "" {
        updates["password"] = req.Password
    }

    if err := h.userService.Update(userID.(string), updates); err != nil {
        c.JSON(http.StatusInternalServerError, Response{
            Status:  "error",
            Message: "Failed to update user",
        })
        return
    }

    c.JSON(http.StatusOK, Response{
        Status:  "success",
        Message: "User updated successfully",
    })
}

// Delete deletes the user with the given ID from the "users" collection in the
// MongoDB database.
//
// The request body should contain no data.
//
// The response will be a JSON object with the following fields:
//   - status: The status of the request. Will be "success" on success, or "error" on error.
//   - message: A human-readable message describing the result of the request.
func (h *UserHandler) Delete(c *gin.Context) {
    userID, _ := c.Get("user_id")

    if err := h.userService.Delete(userID.(string)); err != nil {
        c.JSON(http.StatusInternalServerError, Response{
            Status:  "error",
            Message: "Failed to delete user",
        })
        return
    }

    c.JSON(http.StatusOK, Response{
        Status:  "success",
        Message: "User deleted successfully",
    })
}

func (h *UserHandler) GetMe(c *gin.Context) {
    userID, _ := c.Get("user_id")
    
    user, err := h.userService.GetByID(userID.(string))
    if err != nil {
        c.JSON(http.StatusInternalServerError, Response{
            Status:  "error",
            Message: "Failed to get user information",
        })
        return
    }

    c.JSON(http.StatusOK, Response{
        Status: "success",
        Data:   user,
    })
}
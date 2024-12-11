package handlers

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

type UserHandler struct {
    userService UserService
}

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
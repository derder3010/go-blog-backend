package handlers

import (
	"mime/multipart"
    "github.com/gin-gonic/gin"
    "net/http"
)

// Sửa lại interface UploadService
type UploadService interface {
    UploadImage(file *multipart.FileHeader) (*FileUpload, error)
    DeleteImage(filename string) error
}

type FileUpload struct {
    URL string `json:"url"`
}

type UploadHandler struct {
    uploadService UploadService
}

func NewUploadHandler(uploadService UploadService) *UploadHandler {
    return &UploadHandler{
        uploadService: uploadService,
    }
}

func (h *UploadHandler) UploadImage(c *gin.Context) {
    file, err := c.FormFile("image")
    if err != nil {
        c.JSON(http.StatusBadRequest, Response{
            Status:  "error",
            Message: "No file uploaded",
        })
        return
    }

    validTypes := map[string]bool{
        "image/jpeg": true,
        "image/png":  true,
        "image/gif":  true,
    }
    if !validTypes[file.Header.Get("Content-Type")] {
        c.JSON(http.StatusBadRequest, Response{
            Status:  "error",
            Message: "Invalid file type. Only images are allowed",
        })
        return
    }

    // Pass the multipart.FileHeader to the service
    result, err := h.uploadService.UploadImage(file)
    if err != nil {
        c.JSON(http.StatusInternalServerError, Response{
            Status:  "error",
            Message: "Failed to upload image",
        })
        return
    }

    c.JSON(http.StatusOK, Response{
        Status: "success",
        Data:   result,
    })
}
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

// NewUploadHandler creates a new instance of UploadHandler with the provided UploadService.
// This function takes an UploadService as a parameter and returns a pointer to an UploadHandler.
// The created UploadHandler instance can be used to handle image upload requests.
func NewUploadHandler(uploadService UploadService) *UploadHandler {
    return &UploadHandler{
        uploadService: uploadService,
    }
}

// UploadImage handles image upload requests.
//
// The request body should contain a file under the field name "image".
// The response will be a JSON object with the following fields:
//   - status: The status of the request. Will be "success" on success, or "error" on error.
//   - message: A human-readable message describing the result of the request.
//   - data: A JSON object with a single field "url", which is the URL of the uploaded image.
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
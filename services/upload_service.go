package services

import (
	"fmt"
    "mime/multipart"
    "go-blog-backend/pkg/cloudflare"
    "go-blog-backend/pkg/utils"
)

type UploadService struct {
    r2Client       *cloudflare.R2Client
    imageProcessor *utils.ImageProcessor
}

// NewUploadService creates a new UploadService instance.
//
// The created instance will be configured to use the provided R2 client
// for image uploads and the image processor will be configured with a
// maximum width of 1920px, a maximum height of 1080px, and a quality of 85.
func NewUploadService(r2Client *cloudflare.R2Client) *UploadService {
    return &UploadService{
        r2Client: r2Client,
        imageProcessor: utils.NewImageProcessor(1920, 1080, 85),
    }
}

// UploadImage validates and processes the given image file, then uploads it
// to the Cloudflare R2 bucket. If the image type is invalid, it returns an
// error. If the image could not be processed, it returns an error. If the
// image could not be uploaded, it returns an error.
//
// Note that the processed image is uploaded with the same filename as the
// original image, but with the processed image contents.
func (s *UploadService) UploadImage(file *multipart.FileHeader) (*cloudflare.FileUpload, error) {
    if !s.imageProcessor.ValidateImage(file) {
        return nil, fmt.Errorf("invalid image type")
    }

    processedImage, err := s.imageProcessor.ProcessImage(file)
    if err != nil {
        return nil, err
    }

    newFile := &multipart.FileHeader{
        Filename: file.Filename,
        Header:   file.Header,
        Size:     int64(len(processedImage)),
    }

    return s.r2Client.UploadFile(newFile)
}

// DeleteImage deletes the image with the given filename from the Cloudflare R2
// bucket. The returned error will be non-nil if any error occurred during the
// delete process.
func (s *UploadService) DeleteImage(filename string) error {
    return s.r2Client.DeleteFile(filename)
}
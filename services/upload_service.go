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

func NewUploadService(r2Client *cloudflare.R2Client) *UploadService {
    return &UploadService{
        r2Client: r2Client,
        imageProcessor: utils.NewImageProcessor(1920, 1080, 85),
    }
}

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

func (s *UploadService) DeleteImage(filename string) error {
    return s.r2Client.DeleteFile(filename)
}
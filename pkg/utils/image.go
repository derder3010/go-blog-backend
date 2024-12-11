package utils

import (
    "bytes"
    "image"
    "image/jpeg"
    "image/png"
    "mime/multipart"

    "github.com/nfnt/resize"
)

type ImageProcessor struct {
    maxWidth  uint
    maxHeight uint
    quality   int
}

// NewImageProcessor returns a new instance of ImageProcessor.
//
// The returned ImageProcessor can be used to resize and optimize images.
// The maxWidth and maxHeight parameters are the maximum width and height
// of the processed image, and the quality parameter is the quality of
// the compressed image (0-100).
func NewImageProcessor(maxWidth, maxHeight uint, quality int) *ImageProcessor {
    return &ImageProcessor{
        maxWidth:  maxWidth,
        maxHeight: maxHeight,
        quality:   quality,
    }
}

// ProcessImage resizes and optimizes the image
func (p *ImageProcessor) ProcessImage(file *multipart.FileHeader) ([]byte, error) {
    src, err := file.Open()
    if err != nil {
        return nil, err
    }
    defer src.Close()

    // Decode image
    img, format, err := image.Decode(src)
    if err != nil {
        return nil, err
    }

    // Resize image if needed
    bounds := img.Bounds()
    if uint(bounds.Dx()) > p.maxWidth || uint(bounds.Dy()) > p.maxHeight {
        img = resize.Resize(p.maxWidth, p.maxHeight, img, resize.Lanczos3)
    }

    // Encode image
    buf := new(bytes.Buffer)
    switch format {
    case "jpeg":
        err = jpeg.Encode(buf, img, &jpeg.Options{Quality: p.quality})
    case "png":
        err = png.Encode(buf, img)
    }

    if err != nil {
        return nil, err
    }

    return buf.Bytes(), nil
}

// ValidateImage checks if the file is a valid image
func (p *ImageProcessor) ValidateImage(file *multipart.FileHeader) bool {
    validTypes := map[string]bool{
        "image/jpeg": true,
        "image/png":  true,
        "image/gif":  true,
    }
    return validTypes[file.Header.Get("Content-Type")]
}
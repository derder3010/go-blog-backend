package cloudflare

import (
    "bytes"
    "context"
    "fmt"
    "mime/multipart"
    "time"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/credentials"
    "github.com/aws/aws-sdk-go-v2/service/s3"
)

type R2Client struct {
    client     *s3.Client
    bucketName string
    publicURL  string
}

type R2Config struct {
    AccountID       string
    AccessKeyID     string
    AccessKeySecret string
    BucketName      string
    PublicURL       string
}

// NewR2Client creates a new R2Client instance with the given config.
//
// This function uses the given config to create an AWS client, which is then
// used to interact with the R2 bucket. The client is configured to use the
// given account ID, access key ID, and access key secret to authenticate
// requests. The given bucket name and public URL are used when uploading and
// deleting files.
//
// The returned R2Client is not thread-safe.
func NewR2Client(cfg R2Config) (*R2Client, error) {
    r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
        return aws.Endpoint{
            URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cfg.AccountID),
        }, nil
    })

    awsCfg, err := config.LoadDefaultConfig(context.Background(),
        config.WithEndpointResolverWithOptions(r2Resolver),
        config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
            cfg.AccessKeyID,
            cfg.AccessKeySecret,
            "",
        )),
        config.WithRegion("auto"),
    )
    if err != nil {
        return nil, err
    }

    client := s3.NewFromConfig(awsCfg)

    return &R2Client{
        client:     client,
        bucketName: cfg.BucketName,
        publicURL:  cfg.PublicURL,
    }, nil
}

// FileUpload represents the uploaded file metadata
type FileUpload struct {
    Filename    string
    ContentType string
    Size        int64
    URL         string
}

// UploadFile uploads the provided file to the Cloudflare R2 bucket.
// It reads the file content, generates a unique filename, and sets
// the appropriate content type and cache control headers. The uploaded
// file's metadata is returned upon successful upload. Returns an error
// if the file could not be read or uploaded.
func (c *R2Client) UploadFile(file *multipart.FileHeader) (*FileUpload, error) {
    src, err := file.Open()
    if err != nil {
        return nil, err
    }
    defer src.Close()

    buffer := make([]byte, file.Size)
    if _, err = src.Read(buffer); err != nil {
        return nil, err
    }

    filename := fmt.Sprintf("%d-%s", time.Now().UnixNano(), file.Filename)
    contentType := file.Header.Get("Content-Type")

    input := &s3.PutObjectInput{
        Bucket:       aws.String(c.bucketName),
        Key:          aws.String(filename),
        Body:         bytes.NewReader(buffer),
        ContentType:  aws.String(contentType),
        CacheControl: aws.String("max-age=31536000"),
    }

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    _, err = c.client.PutObject(ctx, input)
    if err != nil {
        return nil, err
    }

    return &FileUpload{
        Filename:    filename,
        ContentType: contentType,
        Size:        file.Size,
        URL:         fmt.Sprintf("%s/%s", c.publicURL, filename),
    }, nil
}

// DeleteFile deletes the file with the given filename from the Cloudflare R2
// bucket. The returned error will be non-nil if any error occurred during the
// delete process.
func (c *R2Client) DeleteFile(filename string) error {
    input := &s3.DeleteObjectInput{
        Bucket: aws.String(c.bucketName),
        Key:    aws.String(filename),
    }

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    _, err := c.client.DeleteObject(ctx, input)
    return err
}
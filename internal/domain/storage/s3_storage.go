package storage

import (
	"mime/multipart"
	"time"

	"github.com/Vanak/Backend/internal/domain/enums"
)

type S3Storage interface {
	UploadFile(bucketType enums.BucketType, key string, file *multipart.FileHeader) error
	DeleteObject(bucketType enums.BucketType, key string) error
	GetPresignedURL(bucketType enums.BucketType, key string, expiration time.Duration) (string, error)
}

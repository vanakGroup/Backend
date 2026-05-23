package storage

import (
	"context"
	"fmt"
	"mime/multipart"
	"slices"
	"time"

	"github.com/Vanak/Backend/bootstrap"
	"github.com/Vanak/Backend/internal/domain/enums"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Storage struct {
	storage *bootstrap.S3
	client  *s3.Client
	buckets map[enums.BucketType]string
}

func NewS3Storage(storage *bootstrap.S3) *S3Storage {
	buckets := make(map[enums.BucketType]string)
	buckets[enums.LogoPic] = storage.Buckets.LogoPic
	buckets[enums.ParkingDoc] = storage.Buckets.ParkingDoc
	buckets[enums.ParkingPic] = storage.Buckets.ParkingPic

	return &S3Storage{
		storage: storage,
		buckets: buckets,
	}
}

func (s3s *S3Storage) setS3Client(bucketType enums.BucketType) error {
	bucketTypes := enums.GetAllBucketTypes()
	if !slices.Contains(bucketTypes, bucketType) {
		return fmt.Errorf("bucket type %v is not valid", bucketType)
	}

	if s3s.client != nil {
		return nil
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(s3s.storage.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			s3s.storage.AccessKey,
			s3s.storage.SecretKey,
			"",
		)),
	)

	if err != nil {
		return err
	}

	s3s.client = s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(s3s.storage.Endpoint)
	})

	return nil
}

func (s3s *S3Storage) UploadFile(bucketType enums.BucketType, key string, file *multipart.FileHeader) error {
	ctx := context.Background()

	err := s3s.setS3Client(bucketType)
	if err != nil {
		return err
	}

	bucket, exists := s3s.buckets[bucketType]

	if !exists {
		return fmt.Errorf("bucket name not configured for type %v", bucketType)
	}

	fileReader, err := file.Open()

	if err != nil {
		return err
	}
	defer fileReader.Close()

	if err := s3s.ensureBucketExists(ctx, bucket); err != nil {
		return fmt.Errorf("bucket error: %w", err)
	}

	_, err = s3s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   fileReader,
	})

	if err != nil {
		return fmt.Errorf("failed to upload: %w", err)
	}

	return nil
}

func (s3s *S3Storage) DeleteObject(bucketType enums.BucketType, key string) error {
	ctx := context.Background()
	err := s3s.setS3Client(bucketType)
	if err != nil {
		return err
	}

	bucket, exists := s3s.buckets[bucketType]

	if !exists {
		return fmt.Errorf("bucket name not configured for type %v", bucketType)
	}

	_, err = s3s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return err
	}

	waiter := s3.NewObjectNotExistsWaiter(s3s.client)
	err = waiter.Wait(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	},
		time.Second*30,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s3s *S3Storage) GetPresignedURL(bucketType enums.BucketType, key string, expiration time.Duration) (string, error) {
	ctx := context.Background()
	err := s3s.setS3Client(bucketType)
	if err != nil {
		return "", err
	}

	bucket, exists := s3s.buckets[bucketType]

	if !exists {
		return "", fmt.Errorf("bucket name not configured for type %v", bucketType)
	}

	presigner := s3.NewPresignClient(s3s.client)
	presignedUrl, err := presigner.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(expiration))

	if err != nil {
		return "", err
	}

	return presignedUrl.URL, nil
}

func (s3s *S3Storage) ensureBucketExists(ctx context.Context, bucket string) error {
	_, err := s3s.client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})

	if err == nil {
		return nil
	}

	_, err = s3s.client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	})

	if err != nil {
		return err
	}

	waiter := s3.NewBucketExistsWaiter(s3s.client)
	err = waiter.Wait(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	},
		time.Second*30,
	)

	if err != nil {
		return err
	}

	return err
}

package service

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOService struct {
	client *minio.Client
	bucket string
}

func NewMinIOService(endpoint, accessKey, secretKey, bucket string, ssl bool) (*MinIOService, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: ssl,
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка создания MinIO клиента: %w", err)
	}

	// Проверяем подключение
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Проверяем существование бакета, если нет - создаем
	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		return nil, fmt.Errorf("ошибка проверки бакета: %w", err)
	}

	if !exists {
		err = client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("ошибка создания бакета: %w", err)
		}

		// Устанавливаем политику доступа для бакета
		policy := `{
            "Version": "2012-10-17",
            "Statement": [
                {
                    "Effect": "Allow",
                    "Principal": {"AWS": "*"},
                    "Action": ["s3:GetObject"],
                    "Resource": ["arn:aws:s3:::` + bucket + `/*"]
                }
            ]
        }`
		err = client.SetBucketPolicy(ctx, bucket, policy)
		if err != nil {
			return nil, fmt.Errorf("ошибка установки политики бакета: %w", err)
		}
	}

	return &MinIOService{
		client: client,
		bucket: bucket,
	}, nil
}

func (m *MinIOService) UploadFile(ctx context.Context, objectName string, file io.Reader, fileSize int64, contentType string) (string, error) {
	// Загружаем файл
	info, err := m.client.PutObject(ctx, m.bucket, objectName, file, fileSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("ошибка загрузки файла в MinIO: %w", err)
	}

	fmt.Printf("Файл успешно загружен: %s, размер: %d байт\n", objectName, info.Size)

	// Генерируем URL для доступа (публичный URL)
	url := fmt.Sprintf("http://%s/%s/%s", m.client.EndpointURL().Host, m.bucket, objectName)

	return url, nil
}

func (m *MinIOService) DeleteFile(ctx context.Context, objectName string) error {
	err := m.client.RemoveObject(ctx, m.bucket, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("ошибка удаления файла из MinIO: %w", err)
	}
	return nil
}

func (m *MinIOService) FileExists(ctx context.Context, objectName string) (bool, error) {
	_, err := m.client.StatObject(ctx, m.bucket, objectName, minio.StatObjectOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return false, nil
		}
		return false, fmt.Errorf("ошибка проверки существования файла: %w", err)
	}
	return true, nil
}

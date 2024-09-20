package file_storage

import (
	"bytes"
	"context"
	"log"
	"net/url"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

const bucket = "pinder"
const userPhotoDir = "/userphoto"

type FileStorage struct {
	minioClient *minio.Client
}

func New(minioClient *minio.Client) *FileStorage {
	return &FileStorage{
		minioClient: minioClient,
	}
}

func (fs *FileStorage) SavePhoto(ctx context.Context, photo []byte) (string, error) {
	objName := uuid.New().String()
	contentType := "application/octet-stream"
	_, err := fs.minioClient.PutObject(
		ctx,
		bucket,
		filepath.Join(userPhotoDir, objName),
		bytes.NewBuffer(photo),
		int64(len(photo)),
		minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		log.Fatalln(err)
	}
	return objName, nil
}

func (fs *FileStorage) DelPhoto(ctx context.Context, photoKey string) error {
	return fs.minioClient.RemoveObject(
		ctx, bucket,
		filepath.Join(userPhotoDir, photoKey),
		minio.RemoveObjectOptions{ForceDelete: true})
}

func (fs *FileStorage) MakeLink(ctx context.Context, photoKey string) (string, error) {
	url, err := fs.minioClient.PresignedGetObject(ctx, bucket, filepath.Join(userPhotoDir, photoKey), time.Hour*24, url.Values{})
	if err != nil {
		return "", err
	}
	return url.String(), nil
}

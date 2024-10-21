package file_storage

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/url"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

const (
	bucket       = "pinder"
	userPhotoDir = "/userphoto"
	chatPhotoDir = "/chatphoto"
	chatVoiceDir = "/chatvoice"
)

type FileStorage struct {
	minioClient *minio.Client
}

func New(minioClient *minio.Client) *FileStorage {
	return &FileStorage{
		minioClient: minioClient,
	}
}

func (fs *FileStorage) saveObj(ctx context.Context, name string, payload []byte) error {
	contentType := "application/octet-stream"
	_, err := fs.minioClient.PutObject(
		ctx,
		bucket,
		name,
		bytes.NewBuffer(payload),
		int64(len(payload)),
		minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		log.Fatalln(err)
	}
	return nil
}

func (fs *FileStorage) getObj(ctx context.Context, name string) (string, error) {
	obj, err := fs.minioClient.GetObject(
		ctx, bucket,
		name,
		minio.GetObjectOptions{})
	if err != nil {
		return "", err
	}
	defer obj.Close()
	b, err := io.ReadAll(obj)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (fs *FileStorage) delObj(ctx context.Context, name string) error {
	return fs.minioClient.RemoveObject(
		ctx, bucket,
		name,
		minio.RemoveObjectOptions{ForceDelete: true})
}

func (fs *FileStorage) shareObj(ctx context.Context, name string) (string, error) {
	url, err := fs.minioClient.PresignedGetObject(ctx, bucket, name, time.Hour*24, url.Values{})
	if err != nil {
		return "", err
	}
	return url.String(), nil
}

func (fs *FileStorage) SaveProfilePhoto(ctx context.Context, body []byte) (string, error) {
	key := uuid.New().String()
	err := fs.saveObj(ctx, filepath.Join(userPhotoDir, key), body)
	if err != nil {
		return "", err
	}
	return key, nil
}

func (fs *FileStorage) DelProfilePhoto(ctx context.Context, key string) error {
	return fs.delObj(ctx, filepath.Join(userPhotoDir, key))
}

func (fs *FileStorage) MakeProfilePhotoLink(ctx context.Context, key string) (string, error) {
	return fs.shareObj(ctx, filepath.Join(userPhotoDir, key))
}

func (fs *FileStorage) SaveChatPhoto(ctx context.Context, body []byte) (string, error) {
	key := uuid.New().String()
	err := fs.saveObj(ctx, filepath.Join(chatPhotoDir, key), body)
	if err != nil {
		return "", err
	}
	return key, nil
}

func (fs *FileStorage) MakeChatPhotoLink(ctx context.Context, key string) (string, error) {
	return fs.shareObj(ctx, filepath.Join(chatPhotoDir, key))
}

func (fs *FileStorage) SaveChatVoice(ctx context.Context, body []byte) (string, error) {
	key := uuid.New().String()
	err := fs.saveObj(ctx, filepath.Join(chatVoiceDir, key), body)
	if err != nil {
		return "", err
	}
	return key, nil
}

func (fs *FileStorage) MakeChatVoiceLink(ctx context.Context, key string) (string, error) {
	return fs.shareObj(ctx, filepath.Join(chatVoiceDir, key))
}

func (fs *FileStorage) GetChatVoice(ctx context.Context, key string) (string, error) {
	return fs.getObj(ctx, filepath.Join(chatVoiceDir, key))
}

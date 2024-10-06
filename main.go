package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/mayye4ka/pinder/authenticator"
	"github.com/mayye4ka/pinder/file_storage"
	grpc_server "github.com/mayye4ka/pinder/grpc-server"
	"github.com/mayye4ka/pinder/repository"
	"github.com/mayye4ka/pinder/service"
	"github.com/mayye4ka/pinder/stt"
	ws_server "github.com/mayye4ka/pinder/ws-server"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	amqp "github.com/rabbitmq/amqp091-go"
	migrate "github.com/rubenv/sql-migrate"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	MinioEndpoint  string `env:"MINIO_ENDPOINT"`
	MinioAccessKey string `env:"MINIO_AK"`
	MinioSecretKey string `env:"MINIO_SK"`
	DbDsn          string `env:"DB_DSN"`
	RabbitMqDsn    string `env:"RABBIT_MQ_DSN"`
	GrpcPort       int    `env:"GRPC_PORT"`
	WsPort         int    `env:"WS_PORT"`
}

func getMinio(config Config) (*minio.Client, error) {
	minioClient, err := minio.New(config.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.MinioAccessKey, config.MinioSecretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, fmt.Errorf("can't get minio: %w", err)
	}
	return minioClient, nil
}

func getDb(config Config) (*gorm.DB, error) {
	migrations := &migrate.FileMigrationSource{
		Dir: "migrations",
	}

	mgrDb, err := sql.Open("mysql", config.DbDsn)
	if err != nil {
		return nil, fmt.Errorf("can't get mysql: %w", err)
	}

	_, err = migrate.Exec(mgrDb, "mysql", migrations, migrate.Up)
	if err != nil {
		return nil, fmt.Errorf("can't get mysql: %w", err)
	}

	db, err := gorm.Open(mysql.Open(config.DbDsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("can't get mysql: %w", err)
	}
	return db, nil
}

func getRabbitMq(config Config) (*amqp.Connection, error) {
	conn, err := amqp.Dial(config.RabbitMqDsn)
	if err != nil {
		return nil, fmt.Errorf("can't get rabbit: %w", err)
	}
	return conn, nil
}

func main() {
	skipEnvLoad := false
	_, err := os.Open(".env")
	if err != nil && errors.Is(err, os.ErrNotExist) {
		skipEnvLoad = true
	}
	if skipEnvLoad {
		err = godotenv.Load()
		if err != nil {
			log.Fatal(err)
		}
	}
	var config Config
	err = env.Parse(&config)
	if err != nil {
		log.Fatal(err)
	}

	db, err := getDb(config)
	if err != nil {
		log.Fatal(err)
	}
	minio, err := getMinio(config)
	if err != nil {
		log.Fatal(err)
	}
	rabbit, err := getRabbitMq(config)
	if err != nil {
		log.Fatal(err)
	}

	fileStorage := file_storage.New(minio)
	repository := repository.New(db)
	stt := stt.New(rabbit)

	auth := authenticator.New(repository)
	notifier := ws_server.NewUserWsNotifier(auth)
	svc := service.New(repository, fileStorage, notifier, stt)

	server := grpc_server.New(svc, auth)

	go func() {
		if err := stt.Start(); err != nil {
			log.Fatal(err)
		}
	}()
	go func() {
		if err := notifier.Start(config.WsPort); err != nil {
			log.Fatal(err)
		}
	}()
	go func() {
		if err := svc.Start(); err != nil {
			log.Fatal(err)
		}
	}()
	go func() {
		if err := server.Start(config.GrpcPort); err != nil {
			log.Fatal(err)
		}
	}()

	select {}
}

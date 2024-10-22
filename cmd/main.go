package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	ntfc_receive "github.com/mayye4ka/pinder/internal/notifications/receive"
	ntfc_send "github.com/mayye4ka/pinder/internal/notifications/send"
	repository "github.com/mayye4ka/pinder/internal/repository/db"
	"github.com/mayye4ka/pinder/internal/repository/file_storage"
	grpc_server "github.com/mayye4ka/pinder/internal/server/grpc-server"
	ws_server "github.com/mayye4ka/pinder/internal/server/ws-server"
	stt_result "github.com/mayye4ka/pinder/internal/stt/result"
	stt_task "github.com/mayye4ka/pinder/internal/stt/task"
	"github.com/mayye4ka/pinder/internal/usecase/authenticator"
	"github.com/mayye4ka/pinder/internal/usecase/service"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	migrate "github.com/rubenv/sql-migrate"
	"golang.org/x/sync/errgroup"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const stopTimeout = time.Minute

type Starter interface {
	Start(context.Context) error
}

type Stopper interface {
	Stop(context.Context) error
}

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
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	skipEnvLoad := false
	_, err := os.Open(".env")
	if err != nil && errors.Is(err, os.ErrNotExist) {
		skipEnvLoad = true
	}
	if !skipEnvLoad {
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
	logger := zerolog.New(os.Stdout)

	fileStorage := file_storage.New(minio, &logger)
	repository := repository.New(db, &logger)
	sttTaskCreator, err := stt_task.NewTaskCreator(rabbit, &logger)
	if err != nil {
		log.Fatal(err)
	}
	ntfcSender, err := ntfc_send.NewNotificationSender(rabbit, &logger)
	if err != nil {
		log.Fatal(err)
	}
	ntfcReceiver := ntfc_receive.NewNotificationReceiver(rabbit, &logger)

	auth := authenticator.New(repository, &logger)
	wsServer := ws_server.NewWsServer(auth, ntfcReceiver, config.WsPort)
	svc := service.New(repository, fileStorage, ntfcSender, sttTaskCreator)
	sttResultReceiver := stt_result.NewResultReceiver(rabbit, svc, &logger)

	server := grpc_server.New(svc, auth, config.GrpcPort)

	eg, egCtx := errgroup.WithContext(ctx)

	for _, s := range []Starter{
		wsServer,
		ntfcReceiver,
		sttResultReceiver,
		server,
	} {
		eg.Go(func() error {
			return s.Start(egCtx)
		})
	}

	err = eg.Wait()
	if err != nil {
		logger.Err(err).Msg("finished with error")
	}

	select {
	case <-termChan:
	case <-egCtx.Done():
	}

	stopCtx, stopCancel := context.WithTimeout(context.Background(), stopTimeout)
	defer stopCancel()
	eg, egCtx = errgroup.WithContext(stopCtx)

	for _, s := range []Stopper{
		wsServer,
		ntfcReceiver,
		sttResultReceiver,
		server,
	} {
		eg.Go(func() error {
			return s.Stop(stopCtx)
		})
	}
	err = eg.Wait()
	if err != nil {
		logger.Err(err).Msg("shutdown error")
	}
}

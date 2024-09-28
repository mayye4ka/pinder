package main

import (
	"database/sql"
	"log"
	"os"
	"pinder/file_storage"
	"pinder/repository"
	"pinder/server"
	"pinder/service"

	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	migrate "github.com/rubenv/sql-migrate"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	minioClient, err := minio.New(os.Getenv("MINIO_ENDPOINT"), &minio.Options{
		Creds:  credentials.NewStaticV4(os.Getenv("MINIO_AK"), os.Getenv("MINIO_SK"), ""),
		Secure: false,
	})
	if err != nil {
		log.Fatal(err)
	}

	filestorage := file_storage.New(minioClient)

	dbDsn := os.Getenv("DB_DSN")
	migrations := &migrate.FileMigrationSource{
		Dir: "migrations",
	}

	mgrDb, err := sql.Open("mysql", dbDsn)
	if err != nil {
		log.Fatal(err)
	}

	_, err = migrate.Exec(mgrDb, "mysql", migrations, migrate.Up)
	if err != nil {
		log.Fatal(err)
	}

	db, err := gorm.Open(mysql.Open(dbDsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	repo := repository.New(db)
	service := service.New(repo, filestorage, nil) // TODO:
	server := server.New(service)
	if err = server.Start(); err != nil {
		log.Fatal(err)
	}
}

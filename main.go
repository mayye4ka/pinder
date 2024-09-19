package main

import (
	"database/sql"
	"log"
	"pinder/repository"
	"pinder/server"
	"pinder/service"

	migrate "github.com/rubenv/sql-migrate"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dbDsn := "pinder:Pinder_1234@tcp(192.168.3.42:3306)/pinder?charset=utf8mb4&parseTime=True&loc=Local"
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
	service := service.New(repo)
	server := server.New(service)
	if err = server.Start(); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/mateoferrari97/AnitiMonono-StudentAPI/cmd/server/internal"
	"github.com/mateoferrari97/AnitiMonono-StudentAPI/cmd/server/internal/service"
	"github.com/mateoferrari97/AnitiMonono-StudentAPI/cmd/server/internal/service/storage"
	"github.com/mateoferrari97/Kit/web/server"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	stg, err := newStorage()
	if err != nil {
		return err
	}

	svc := service.NewService(stg)
	sv := server.NewServer()
	handler := internal.NewHandler(sv, svc)

	handler.GetStudentSubjects()
	handler.GetSubjectDetails()
	handler.GetProfessorships()

	return sv.Run(getPort())
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		return ":8080"
	}

	return ":" + port
}

func newStorage() (*storage.Storage, error) {
	source := os.Getenv("DATABASE_CONFIG")
	if source == "" {
		source = "root:root@tcp(localhost:3306)/university"
	}

	db, err := sqlx.Open("mysql", source)
	if err != nil {
		return nil, fmt.Errorf("could not connect to db: %v", err)
	}

	return storage.NewStorage(db), nil
}

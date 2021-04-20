package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/mateoferrari97/my-path/cmd/server/internal"
	"github.com/mateoferrari97/my-path/cmd/server/internal/service"
	"github.com/mateoferrari97/my-path/cmd/server/internal/service/storage"
	"github.com/mateoferrari97/my-path/internal/server"
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

	return sv.Run(":8080")
}

func newStorage() (*storage.Storage, error) {
	db, err := sqlx.Connect("mysql", "root:root@tcp(localhost:3306)/university")
	if err != nil {
		return nil, fmt.Errorf("could not connect to db: %v", err)
	}

	return storage.NewStorage(db), nil
}

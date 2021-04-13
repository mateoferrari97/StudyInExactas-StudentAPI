package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"

	"github.com/mateoferrari97/my-path/src/server/internal"
	"github.com/mateoferrari97/my-path/src/server/internal/service"
	"github.com/mateoferrari97/my-path/src/server/internal/service/storage"
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
	handler := internal.NewHandler(svc)

	router := mux.NewRouter()
	router.HandleFunc("/students/{studentEmail}/careers/{careerID}/subjects", handler.GetStudentSubjects)
	router.HandleFunc("/careers/{careerID}/subjects/{subjectID}", handler.GetSubjectDetails)
	/*
		router.HandleFunc("/careers/{careerID}/subjects/{subjectID}/professorships", handler.GetSubjectProfessorships)
		router.HandleFunc("/careers/{careerID}/subjects/{subjectID}/materials", handler.GetSubjectMaterials)
	*/

	log.Println("starting server on port :8080")
	return http.ListenAndServe(":8080", router)
}

func newStorage() (*storage.Storage, error) {
	db, err := sqlx.Connect("mysql", "root:root@tcp(localhost:3306)/university")
	if err != nil {
		return nil, fmt.Errorf("could not connect to db: %v", err)
	}

	return storage.NewStorage(db), nil
}

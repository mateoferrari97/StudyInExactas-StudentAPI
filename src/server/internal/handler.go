package internal

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mateoferrari97/my-path/src/server/internal/service"
)

type Service interface {
	GetStudentSubjects(studentEmail, careerID string) ([]byte, error)
	GetSubjectDetails(subjectID, careerID string) ([]byte, error)
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetStudentSubjects(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	studentEmail, exist := params["studentEmail"]
	if !exist || studentEmail == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("student email is required"))
		return
	}

	careerID, exist := params["careerID"]
	if !exist || careerID == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("student career is required"))
		return
	}

	studentSubjects, err := h.service.GetStudentSubjects(studentEmail, careerID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if errors.Is(err, service.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
		}

		_, _ = w.Write([]byte(err.Error()))
		return
	}

	_, _ = w.Write(studentSubjects)
}

func (h *Handler) GetSubjectDetails(w http.ResponseWriter, r *http.Request) {
	/*
		"/career/{careerID}/subject/{id}"
		{
			"id": "123",
			"hours": 123,
			"type": OBLIGATORIA/ELECTIVA,
			"points": 5,
			"meet",
			"uri":
		}
	*/
	params := mux.Vars(r)
	careerID, exist := params["careerID"]
	if !exist || careerID == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("career id is required"))
		return
	}

	subjectID, exist := params["subjectID"]
	if !exist || subjectID == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("subject id is required"))
		return
	}

	subjectDetails, err := h.service.GetSubjectDetails(subjectID, careerID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if errors.Is(err, service.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
		}

		_, _ = w.Write([]byte(err.Error()))
		return
	}

	_, _ = w.Write(subjectDetails)
}

package internal

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/mateoferrari97/Kit/web/server"
	"github.com/mateoferrari97/my-path/cmd/server/internal/service"
)

type Wrapper interface {
	Wrap(method, pattern string, f server.HandlerFunc, mws ...server.Middleware)
}

type Service interface {
	GetStudentSubjects(studentEmail, careerID string) ([]byte, error)
	GetSubjectDetails(subjectID, careerID string) ([]byte, error)
	GetProfessorships(subjectID, careerID string) ([]byte, error)
}

type Handler struct {
	wrapper Wrapper
	service Service
}

func NewHandler(wrapper Wrapper, service Service) *Handler {
	return &Handler{
		wrapper: wrapper,
		service: service,
	}
}

func (h *Handler) GetStudentSubjects() {
	wrapH := func(w http.ResponseWriter, r *http.Request) error {
		params := mux.Vars(r)
		studentEmail, exist := params["studentEmail"]
		if !exist || studentEmail == "" {
			return server.NewError("student email is required", http.StatusBadRequest)
		}

		careerID, exist := params["careerID"]
		if !exist || careerID == "" {
			return server.NewError("career id is required", http.StatusBadRequest)
		}

		studentSubjects, err := h.service.GetStudentSubjects(studentEmail, careerID)
		if err != nil {
			if errors.Is(err, service.ErrNotFound) {
				return server.NewError(err.Error(), http.StatusNotFound)
			}

			return err
		}

		return server.RespondJSON(w, studentSubjects, http.StatusOK)
	}

	h.wrapper.Wrap(http.MethodGet, "/students/{studentEmail}/careers/{careerID}/subjects", wrapH)
}

func (h *Handler) GetSubjectDetails() {
	wrapH := func(w http.ResponseWriter, r *http.Request) error {
		params := mux.Vars(r)
		careerID, exist := params["careerID"]
		if !exist || careerID == "" {
			return server.NewError("career id is required", http.StatusBadRequest)
		}

		subjectID, exist := params["subjectID"]
		if !exist || subjectID == "" {
			return server.NewError("subject id is required", http.StatusBadRequest)
		}

		subjectDetails, err := h.service.GetSubjectDetails(subjectID, careerID)
		if err != nil {
			if errors.Is(err, service.ErrNotFound) {
				return server.NewError(err.Error(), http.StatusNotFound)
			}

			return err
		}

		return server.RespondJSON(w, subjectDetails, http.StatusOK)
	}

	h.wrapper.Wrap(http.MethodGet, "/careers/{careerID}/subjects/{subjectID}", wrapH)
}

func (h *Handler) GetProfessorships() {
	wrapH := func(w http.ResponseWriter, r *http.Request) error {
		params := mux.Vars(r)
		careerID, exist := params["careerID"]
		if !exist || careerID == "" {
			return server.NewError("career id is required", http.StatusBadRequest)
		}

		subjectID, exist := params["subjectID"]
		if !exist || subjectID == "" {
			return server.NewError("subject id is required", http.StatusBadRequest)
		}

		professorships, err := h.service.GetProfessorships(subjectID, careerID)
		if err != nil {
			if errors.Is(err, service.ErrNotFound) {
				return server.NewError(err.Error(), http.StatusNotFound)
			}

			return err
		}

		return server.RespondJSON(w, professorships, http.StatusOK)
	}

	h.wrapper.Wrap(http.MethodGet, "/careers/{careerID}/subjects/{subjectID}/professorships", wrapH)
}
package internal

import (
	"encoding/json"
	"errors"
	"gopkg.in/go-playground/validator.v9"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/mateoferrari97/AnitiMonono-StudentAPI/cmd/server/internal/service"
	"github.com/mateoferrari97/Kit/web/server"
)

type Wrapper interface {
	Wrap(method, pattern string, f server.HandlerFunc, mws ...server.Middleware)
}

type Service interface {
	CreateStudent(name, studentEmail string) error
	AssignStudentToCareer(studentEmail, careerID string) error
	GetStudentSubjects(studentEmail, careerID string) ([]byte, error)
	UpdateStudentSubject(req service.UpdateStudentSubjectRequest) error
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

func (h *Handler) CreateStudent() {
	wrapH := func(w http.ResponseWriter, r *http.Request) error {
		var studentInformation struct {
			Name         string `json:"name" validate:"required"`
			StudentEmail string `json:"student_email" validate:"required"`
		}

		if err := json.NewDecoder(r.Body).Decode(&studentInformation); err != nil {
			return server.NewError(err.Error(), http.StatusUnprocessableEntity)
		}

		if err := validate.Struct(studentInformation); err != nil {
			return server.NewError(err.Error(), http.StatusBadRequest)
		}

		if err := h.service.CreateStudent(studentInformation.Name, studentInformation.StudentEmail); err != nil {
			if errors.Is(err, service.ErrStudentAlreadyExist) {
				return server.NewError(err.Error(), http.StatusConflict)
			}

			return err
		}

		return server.RespondJSON(w, nil, http.StatusOK)
	}

	h.wrapper.Wrap(http.MethodPost, "/students", wrapH)
}

func (h *Handler) AssignStudentToCareer() {
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

		if err := h.service.AssignStudentToCareer(studentEmail, careerID); err != nil {
			switch {
			case errors.Is(err, service.ErrNotFound):
				return server.NewError(err.Error(), http.StatusNotFound)
			case errors.Is(err, service.ErrCareerAlreadyAssigned):
				return server.NewError(err.Error(), http.StatusConflict)
			case errors.Is(err, service.ErrMaxCareerReached):
				return server.NewError(err.Error(), http.StatusConflict)
			default:
				return err
			}
		}

		return nil
	}

	h.wrapper.Wrap(http.MethodPost, "/students/{studentEmail}/careers/{careerID}", wrapH)
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

var validate = validator.New()

func (h *Handler) UpdateStudentSubject() {
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

		subjectID, exist := params["subjectID"]
		if !exist || subjectID == "" {
			return server.NewError("subject id is required", http.StatusBadRequest)
		}

		var subjectInformation struct {
			Status      string `json:"status" validate:"required,oneof=PENDIENTE APROBADA"`
			Description string `json:"description" validate:"omitempty,min=1,max=128"`
		}

		if err := json.NewDecoder(r.Body).Decode(&subjectInformation); err != nil {
			return server.NewError(err.Error(), http.StatusUnprocessableEntity)
		}

		if err := validate.Struct(subjectInformation); err != nil {
			return server.NewError(err.Error(), http.StatusBadRequest)
		}

		if err := h.service.UpdateStudentSubject(service.UpdateStudentSubjectRequest{
			StudentEmail: studentEmail,
			CareerID:     careerID,
			SubjectID:    subjectID,
			Status:       subjectInformation.Status,
			Description:  subjectInformation.Description,
		}); err != nil {
			if errors.Is(err, service.ErrNotFound) {
				return server.NewError(err.Error(), http.StatusNotFound)
			}

			return err
		}

		return server.RespondJSON(w, nil, http.StatusOK)
	}

	h.wrapper.Wrap(http.MethodPut, "/students/{studentEmail}/careers/{careerID}/subjects/{subjectID}", wrapH)
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

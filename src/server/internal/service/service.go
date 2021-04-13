package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/mateoferrari97/my-path/src/server/internal/service/storage"
)

var ErrNotFound = errors.New("service: resource not found")

type Storage interface {
	GetStudentSubjects(studentEmail, careerID string) ([]storage.StudentSubject, error)
	GetSubjectDetails(subjectID, careerID string) (storage.SubjectDetails, error)
}

type Service struct {
	storage Storage
}

func NewService(storage Storage) *Service {
	return &Service{storage: storage}
}

func (s *Service) GetStudentSubjects(studentEmail, careerID string) ([]byte, error) {
	type (
		studentSubject struct {
			ID          int     `json:"id"`
			Name        string  `json:"name"`
			Type        string  `json:"type"`
			Status      string  `json:"status"`
			Description *string `json:"description"`
		}

		getStudentSubjectsResponse struct {
			SubjectsCorrelatives map[string][]int          `json:"subjects_correlatives"`
			Subjects             map[string]studentSubject `json:"subjects"`
		}
	)

	studentSubjects, err := s.storage.GetStudentSubjects(studentEmail, careerID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, fmt.Errorf("could not get student subjects from [email: %s]: %w", studentEmail, ErrNotFound)
		}

		return nil, fmt.Errorf("could not get student subjects from [email: %s]: %v", studentEmail, err)
	}

	subjects := make(map[string]studentSubject, len(studentSubjects))
	subjectsCorrelatives := make(map[string][]int, len(studentSubjects))
	for _, subject := range studentSubjects {
		from := subject.CorrelativeID
		to := strconv.Itoa(subject.ID)

		if subjectsCorrelatives[to] == nil {
			subjectsCorrelatives[to] = []int{}
		}

		if hasCorrelative(from) {
			subjectsCorrelatives[to] = append(subjectsCorrelatives[to], subject.CorrelativeID)
		}

		subjects[to] = studentSubject{
			ID:          subject.ID,
			Name:        subject.Name,
			Type:        subject.Type,
			Status:      subject.Status,
			Description: subject.Description,
		}
	}

	response, err := json.Marshal(getStudentSubjectsResponse{
		SubjectsCorrelatives: subjectsCorrelatives,
		Subjects:             subjects,
	})

	if err != nil {
		return nil, fmt.Errorf("could not marshal response: %v", err)
	}

	return response, nil
}

func hasCorrelative(correlativeID int) bool {
	return correlativeID != 0
}

func (s *Service) GetSubjectDetails(subjectID, careerID string) ([]byte, error) {
	type (
		subjectDetailsResponse struct {
			ID     int     `json:"id"`
			Hours  int     `json:"hours"`
			Points int     `json:"points"`
			Name   string  `json:"name"`
			Type   string  `json:"type"`
			URI    *string `json:"uri"`
			Meet   *string `json:"meet"`
		}
	)

	subjectDetails, err := s.storage.GetSubjectDetails(subjectID, careerID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, fmt.Errorf("could not get subject details from [id: %s]: %w", subjectID, ErrNotFound)
		}

		return nil, fmt.Errorf("could not get subject details from [id: %s]: %v", subjectID, err)
	}

	response, err := json.Marshal(subjectDetailsResponse{
		ID:     subjectDetails.ID,
		Hours:  subjectDetails.Hours,
		Points: subjectDetails.Points,
		Type:   subjectDetails.Type,
		Name:   subjectDetails.Name,
		URI:    subjectDetails.URI,
		Meet:   subjectDetails.Meet,
	})

	if err != nil {
		return nil, fmt.Errorf("could not marshal response: %v", err)
	}

	return response, nil
}

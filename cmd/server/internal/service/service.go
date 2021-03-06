package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"

	"github.com/mateoferrari97/AnitiMonono-StudentAPI/cmd/server/internal/service/storage"
)

const maxCareersPerStudent = 2

var (
	ErrNotFound              = errors.New("service: resource not found")
	ErrCareerAlreadyAssigned = errors.New("service: career already assigned")
	ErrMaxCareerReached      = errors.New("service: student already has maximum careers assigned")
	ErrStudentAlreadyExist   = errors.New("service: student already exist")
)

var (
	dayToDayNumber = map[string]int{
		"Lunes":     1,
		"Martes":    2,
		"Miércoles": 3,
		"Jueves":    4,
		"Viernes":   5,
		"Sábado":    6,
		"Domingo":   7,
	}

	dayNumberToDay = map[int]string{
		1: "Lunes",
		2: "Martes",
		3: "Miércoles",
		4: "Jueves",
		5: "Viernes",
		6: "Sábado",
		7: "Domingo",
	}
)

type Storage interface {
	CreateStudent(name, studentEmail string) error
	GetStudentSubjects(studentEmail, careerID string) ([]storage.StudentSubject, error)
	GetSubjectDetails(subjectID, careerID string) (storage.SubjectDetails, error)
	GetProfessorships(subjectID, careerID string) ([]storage.Professorship, error)
	GetStudentCareerIDs(studentEmail string) ([]int, error)
	AssignStudentToCareer(studentEmail, careerID string) error
	UpdateStudentSubject(req storage.UpdateStudentSubjectRequest) error
}

type Service struct {
	storage Storage
}

func NewService(storage Storage) *Service {
	return &Service{storage: storage}
}

func (s *Service) CreateStudent(name, studentEmail string) error {
	if err := s.storage.CreateStudent(name, studentEmail); err != nil {
		if errors.Is(err, storage.ErrResourceAlreadyExist) {
			return ErrStudentAlreadyExist
		}

		return err
	}

	return nil
}

func (s *Service) AssignStudentToCareer(studentEmail, careerID string) error {
	careersIDs, err := s.storage.GetStudentCareerIDs(studentEmail)
	if err != nil && !errors.Is(err, storage.ErrNotFound) {
		return fmt.Errorf("could not get student careers [student_email: %s]: %v", studentEmail, err)
	}

	if careersIDs != nil {
		for _, id := range careersIDs {
			if strconv.Itoa(id) == careerID {
				return fmt.Errorf("student [student_email: %s] %w", studentEmail, ErrCareerAlreadyAssigned)
			}
		}

		if len(careersIDs) >= maxCareersPerStudent {
			return fmt.Errorf("student [student_email: %s] %w", studentEmail, ErrMaxCareerReached)
		}
	}

	if err := s.storage.AssignStudentToCareer(studentEmail, careerID); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return fmt.Errorf("could not assign student [student_email: %s] to career: %w", studentEmail, ErrNotFound)
		}

		return fmt.Errorf("could not assign student [student_email: %s] to career: %v", studentEmail, err)
	}

	return nil
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
			Correlatives map[string][]int          `json:"correlatives"`
			Subjects     map[string]studentSubject `json:"subjects"`
		}
	)

	studentSubjects, err := s.storage.GetStudentSubjects(studentEmail, careerID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, fmt.Errorf("could not get subjects: %w", ErrNotFound)
		}

		return nil, fmt.Errorf("could not get subjects: %v", err)
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
		Correlatives: subjectsCorrelatives,
		Subjects:     subjects,
	})

	if err != nil {
		return nil, fmt.Errorf("could not marshal response: %v", err)
	}

	return response, nil
}

type UpdateStudentSubjectRequest struct {
	StudentEmail string
	CareerID     string
	SubjectID    string
	Status       string
	Description  string
}

func (s *Service) UpdateStudentSubject(req UpdateStudentSubjectRequest) error {
	storageReq := storage.UpdateStudentSubjectRequest{
		StudentEmail: req.StudentEmail,
		CareerID:     req.CareerID,
		SubjectID:    req.SubjectID,
		Status:       req.Status,
	}

	if req.Description != "" {
		storageReq.Description = &req.Description
	}

	if err := s.storage.UpdateStudentSubject(storageReq); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return fmt.Errorf("could not update subject: %w: %v", ErrNotFound, err)
		}

		return fmt.Errorf("could not update subject: %v", err)
	}

	return nil
}

func hasCorrelative(correlativeID int) bool {
	return correlativeID != 0
}

func (s *Service) GetSubjectDetails(subjectID, careerID string) ([]byte, error) {
	type subjectDetailsResponse struct {
		ID     int     `json:"id"`
		Name   string  `json:"name"`
		Type   string  `json:"type"`
		URI    *string `json:"uri"`
		Meet   *string `json:"meet"`
		Hours  *int    `json:"hours"`
		Points *int    `json:"points"`
	}

	subjectDetails, err := s.storage.GetSubjectDetails(subjectID, careerID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, fmt.Errorf("could not get subject details: %w", ErrNotFound)
		}

		return nil, fmt.Errorf("could not get subject details: %v", err)
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

func (s *Service) GetProfessorships(subjectID, careerID string) ([]byte, error) {
	type schedule struct {
		Day   string `json:"day"`
		Start string `json:"start"`
		End   string `json:"end"`
	}

	professorships, err := s.storage.GetProfessorships(subjectID, careerID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, fmt.Errorf("could not get professorships: %w", ErrNotFound)
		}

		return nil, fmt.Errorf("could not get professorships: %v", err)
	}

	professorshipInformation := make(map[string][]schedule, len(professorships))
	for _, professorship := range professorships {
		day, err := convertDayNumberToDay(professorship.Day)
		if err != nil {
			return nil, err
		}

		start, err := trimSecondsFromTime(professorship.Start)
		if err != nil {
			return nil, err
		}

		end, err := trimSecondsFromTime(professorship.End)
		if err != nil {
			return nil, err
		}

		if _, exist := professorshipInformation[professorship.Name]; !exist {
			professorshipInformation[professorship.Name] = []schedule{}
		}

		professorshipInformation[professorship.Name] = append(professorshipInformation[professorship.Name], schedule{
			Day:   day,
			Start: start,
			End:   end,
		})
	}

	for _, schedule := range professorshipInformation {
		sort.Slice(schedule, func(i, j int) bool {
			return isTargetLessThanCandidate(schedule[i].Day, schedule[j].Day)
		})
	}

	response, err := json.Marshal(professorshipInformation)
	if err != nil {
		return nil, fmt.Errorf("could not marshal response: %v", err)
	}

	return response, nil
}

func convertDayNumberToDay(dayNumber int) (string, error) {
	day, exist := dayNumberToDay[dayNumber]
	if !exist {
		return "", fmt.Errorf("could not convert day number [dayNumber: %d] to day", dayNumber)
	}

	return day, nil
}

func trimSecondsFromTime(time string) (string, error) {
	length := len(time)
	if length < 6 { // Expected: HH:MM:SS
		return "", errors.New("could not trim seconds from time")
	}

	return time[:length-3], nil
}

func isTargetLessThanCandidate(target, candidate string) bool {
	t := dayToDayNumber[target]
	c := dayToDayNumber[candidate]

	return t < c
}

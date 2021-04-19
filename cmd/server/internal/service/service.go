package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mateoferrari97/my-path/cmd/server/internal/service/storage"
	"sort"
	"strconv"
)

var ErrNotFound = errors.New("service: resource not found")

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
	GetStudentSubjects(studentEmail, careerID string) ([]storage.StudentSubject, error)
	GetSubjectDetails(subjectID, careerID string) (storage.SubjectDetails, error)
	GetProfessorshipSchedules(subjectID, careerID string) ([]storage.ProfessorshipSchedule, error)
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
	type subjectDetailsResponse struct {
		ID     int     `json:"id"`
		Hours  int     `json:"hours"`
		Points int     `json:"points"`
		Name   string  `json:"name"`
		Type   string  `json:"type"`
		URI    *string `json:"uri"`
		Meet   *string `json:"meet"`
	}

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

func (s *Service) GetProfessorshipSchedules(subjectID, careerID string) ([]byte, error) {
	type schedule struct {
		Day   string `json:"day"`
		Start string `json:"start"`
		End   string `json:"end"`
	}

	professorshipSchedules, err := s.storage.GetProfessorshipSchedules(subjectID, careerID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, fmt.Errorf("could not get professorship schedules from [subject_id: %s]: %w", subjectID, ErrNotFound)
		}

		return nil, fmt.Errorf("could not get professorship schedules from [subject_id: %s]: %v", subjectID, err)
	}

	schedules := make(map[string][]schedule, len(professorshipSchedules))
	for _, professorship := range professorshipSchedules {
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

		if _, exist := schedules[professorship.Name]; !exist {
			schedules[professorship.Name] = []schedule{}
		}

		schedules[professorship.Name] = append(schedules[professorship.Name], schedule{
			Day:   day,
			Start: start,
			End:   end,
		})
	}

	for _, schedule := range schedules {
		sort.Slice(schedule, func(i, j int) bool {
			return isTargetLessThanCandidate(schedule[i].Day, schedule[j].Day)
		})
	}

	response, err := json.Marshal(schedules)
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

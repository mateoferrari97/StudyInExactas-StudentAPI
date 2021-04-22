package service

import (
	"errors"
	"github.com/mateoferrari97/my-path/cmd/server/internal/service/storage"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

type storageMock struct {
	mock.Mock
}

func (s *storageMock) GetStudentSubjects(studentEmail, careerID string) ([]storage.StudentSubject, error) {
	args := s.Called(studentEmail, careerID)
	return args.Get(0).([]storage.StudentSubject), args.Error(1)
}

func (s *storageMock) GetSubjectDetails(subjectID, careerID string) (storage.SubjectDetails, error) {
	args := s.Called(subjectID, careerID)
	return args.Get(0).(storage.SubjectDetails), args.Error(1)
}

func (s *storageMock) GetProfessorships(subjectID, careerID string) ([]storage.Professorship, error) {
	args := s.Called(subjectID, careerID)
	return args.Get(0).([]storage.Professorship), args.Error(1)
}

func TestService_GetStudentSubjects(t *testing.T) {
	// Given
	storage_ := storageMock{}
	storage_.On("GetStudentSubjects", "example@gmail.com", "1").Return([]storage.StudentSubject{
		{
			ID:            1,
			CorrelativeID: 0,
			Status:        "PENDING",
			Name:          "Subject 1",
			Type:          "REQUIRED",
			Description:   nil,
		},
		{
			ID:            2,
			CorrelativeID: 1,
			Status:        "PENDING",
			Name:          "Subject 2",
			Type:          "REQUIRED",
			Description:   nil,
		},
	}, nil)

	s := NewService(&storage_)

	// When
	subjects, err := s.GetStudentSubjects("example@gmail.com", "1")
	if err != nil {
		t.Fatal(err)
	}

	// Then
	require.Equal(t, []byte(`{"correlatives":{"1":[],"2":[1]},"subjects":{"1":{"id":1,"name":"Subject 1","type":"REQUIRED","status":"PENDING","description":null},"2":{"id":2,"name":"Subject 2","type":"REQUIRED","status":"PENDING","description":null}}}`), subjects)
}

func TestService_GetStudentSubjects_StorageError(t *testing.T) {
	// Given
	storage_ := storageMock{}
	storage_.On("GetStudentSubjects", "example@gmail.com", "1").Return([]storage.StudentSubject{}, errors.New("error"))

	s := NewService(&storage_)

	// When
	_, err := s.GetStudentSubjects("example@gmail.com", "1")
	if err == nil {
		t.Fatal("test must fail")
	}

	// Then
	require.EqualError(t, err, "could not get student subjects from [student_email: example@gmail.com and career_id: 1]: error")
}

func TestService_GetStudentSubjects_StorageNotFoundError(t *testing.T) {
	// Given
	storage_ := storageMock{}
	storage_.On("GetStudentSubjects", "example@gmail.com", "1").Return([]storage.StudentSubject{}, storage.ErrNotFound)

	s := NewService(&storage_)

	// When
	_, err := s.GetStudentSubjects("example@gmail.com", "1")
	if err == nil {
		t.Fatal("test must fail")
	}

	// Then
	require.EqualError(t, err, "could not get student subjects from [email: example@gmail.com]: service: resource not found")
}

func TestService_GetSubjectDetails(t *testing.T) {
	// Given
	storage_ := storageMock{}
	storage_.On("GetSubjectDetails", "1", "2").Return(storage.SubjectDetails{
		ID:     1,
		Hours:  240,
		Points: 8,
		Name:   "Algebra",
		Type:   "REQUIRED",
		URI:    nil,
		Meet:   nil,
	}, nil)

	s := NewService(&storage_)

	// When
	subject, err := s.GetSubjectDetails("1", "2")
	if err != nil {
		t.Fatal(err)
	}

	// Then
	require.Equal(t, []byte(`{"id":1,"hours":240,"points":8,"name":"Algebra","type":"REQUIRED","uri":null,"meet":null}`), subject)
}

func TestService_GetSubjectDetails_StorageError(t *testing.T) {
	// Given
	storage_ := storageMock{}
	storage_.On("GetSubjectDetails", "1", "2").Return(storage.SubjectDetails{}, errors.New("error"))

	s := NewService(&storage_)

	// When
	_, err := s.GetSubjectDetails("1", "2")
	if err == nil {
		t.Fatal("test must fail")
	}

	// Then
	require.EqualError(t, err, "could not get subject details from [subect_id: 1 and career_id: 2]: error")
}

func TestService_GetSubjectDetails_StorageNotFoundError(t *testing.T) {
	// Given
	storage_ := storageMock{}
	storage_.On("GetSubjectDetails", "1", "2").Return(storage.SubjectDetails{}, storage.ErrNotFound)

	s := NewService(&storage_)

	// When
	_, err := s.GetSubjectDetails("1", "2")
	if err == nil {
		t.Fatal("test must fail")
	}

	// Then
	require.EqualError(t, err, "could not get subject details from [subect_id: 1 and career_id: 2]: service: resource not found")
}

func TestService_GetProfessorships(t *testing.T) {
	// Given
	storage_ := storageMock{}
	storage_.On("GetProfessorships", "1", "2").Return([]storage.Professorship{
		{
			Day:   1,
			Name:  "CATEDRA 1",
			Start: "17:00:00",
			End:   "21:00:00",
		},
		{
			Day:   1,
			Name:  "CATEDRA 2",
			Start: "9:00:00",
			End:   "12:00:00",
		},
		{
			Day:   2,
			Name:  "CATEDRA 2",
			Start: "9:00:00",
			End:   "12:00:00",
		},
	}, nil)

	s := NewService(&storage_)

	// When
	professorships, err := s.GetProfessorships("1", "2")
	if err != nil {
		t.Fatal(err)
	}

	// Then
	require.Equal(t, []byte(`{"CATEDRA 1":[{"day":"Lunes","start":"17:00","end":"21:00"}],"CATEDRA 2":[{"day":"Lunes","start":"9:00","end":"12:00"},{"day":"Martes","start":"9:00","end":"12:00"}]}`), professorships)
}

func TestService_GetProfessorships_StorageError(t *testing.T) {
	// Given
	storage_ := storageMock{}
	storage_.On("GetProfessorships", "1", "2").Return([]storage.Professorship{}, errors.New("error"))

	s := NewService(&storage_)

	// When
	_, err := s.GetProfessorships("1", "2")
	if err == nil {
		t.Fatal("test must fail")
	}

	// Then
	require.EqualError(t, err, "could not get professorships from [subject_id: 1 and career_id: 2]: error")
}

func TestService_GetProfessorships_StorageNotFoundError(t *testing.T) {
	// Given
	storage_ := storageMock{}
	storage_.On("GetProfessorships", "1", "2").Return([]storage.Professorship{}, storage.ErrNotFound)

	s := NewService(&storage_)

	// When
	_, err := s.GetProfessorships("1", "2")
	if err == nil {
		t.Fatal("test must fail")
	}

	// Then
	require.EqualError(t, err, "could not get professorship professorships from [subject_id: 1 and career_id: 2]: service: resource not found")
}

func TestService_GetProfessorships_DayNotExist(t *testing.T) {
	// Given
	storage_ := storageMock{}
	storage_.On("GetProfessorships", "1", "2").Return([]storage.Professorship{
		{
			Day:   9,
			Name:  "CATEDRA 1",
			Start: "17:00:00",
			End:   "21:00:00",
		},
	}, nil)

	s := NewService(&storage_)

	// When
	_, err := s.GetProfessorships("1", "2")
	if err == nil {
		t.Fatal("test must fail")
	}

	// Then
	require.EqualError(t, err, "could not convert day number [dayNumber: 9] to day")
}

func TestService_GetProfessorships_InvalidTimeLength(t *testing.T) {
	tt := []struct {
		name  string
		start string
		end   string
	}{
		{
			name:  "invalid start time",
			start: "17",
			end:   "20:00:00",
		},
		{
			name:  "invalid end time",
			start: "17:00:00",
			end:   "20",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// Given
			storage_ := storageMock{}
			storage_.On("GetProfessorships", "1", "2").Return([]storage.Professorship{
				{
					Day:   1,
					Name:  "CATEDRA 1",
					Start: tc.start,
					End:   tc.end,
				},
			}, nil)

			s := NewService(&storage_)

			// When
			_, err := s.GetProfessorships("1", "2")
			if err == nil {
				t.Fatal("test must fail")
			}

			// Then
			require.EqualError(t, err, "could not trim seconds from time")
		})
	}
}

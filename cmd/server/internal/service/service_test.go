package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/mateoferrari97/AnitiMonono-StudentAPI/cmd/server/internal/service/storage"
)

type storageMock struct {
	mock.Mock
}

func (s *storageMock) GetStudentCareerIDs(studentEmail string) ([]int, error) {
	args := s.Called(studentEmail)
	return args.Get(0).([]int), args.Error(1)
}

func (s *storageMock) AssignStudentToCareer(studentEmail, careerID string) error {
	args := s.Called(studentEmail, careerID)
	return args.Error(0)
}

func (s *storageMock) UpdateStudentSubject(req storage.UpdateStudentSubjectRequest) error {
	args := s.Called(req)
	return args.Error(0)
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

func TestService_AssignStudentToCareer(t *testing.T) {
	// Given
	storage_ := storageMock{}
	storage_.On("GetStudentCareerIDs", "example@gmail.com").Return([]int{}, nil)
	storage_.On("AssignStudentToCareer", "example@gmail.com", "1").Return(nil)

	s := NewService(&storage_)

	// When
	err := s.AssignStudentToCareer("example@gmail.com", "1")
	if err != nil {
		t.Fatal(err)
	}

	// Then
	require.Nil(t, err)
}

func TestService_AssignStudentToCareer_GetStudentCareerIDsError(t *testing.T) {
	// Given
	storage_ := storageMock{}
	storage_.On("GetStudentCareerIDs", "example@gmail.com").Return([]int{}, errors.New("error"))

	s := NewService(&storage_)

	// When
	err := s.AssignStudentToCareer("example@gmail.com", "1")
	if err == nil {
		t.Fatal("test must fail")
	}

	// Then
	require.EqualError(t, err, "could not get student careers [student_email: example@gmail.com]: error")
}

func TestService_AssignStudentToCareer_GetStudentCareerIDsNotFoundError(t *testing.T) {
	// Given
	storage_ := storageMock{}
	storage_.On("GetStudentCareerIDs", "example@gmail.com").Return([]int{}, storage.ErrNotFound)
	storage_.On("AssignStudentToCareer", "example@gmail.com", "1").Return(nil)

	s := NewService(&storage_)

	// When
	err := s.AssignStudentToCareer("example@gmail.com", "1")
	if err != nil {
		t.Fatal(err)
	}

	// Then
	require.Nil(t, err)
}

func TestService_AssignStudentToCareer_CareerAlreadyAssignedError(t *testing.T) {
	// Given
	storage_ := storageMock{}
	storage_.On("GetStudentCareerIDs", "example@gmail.com").Return([]int{1}, nil)

	s := NewService(&storage_)

	// When
	err := s.AssignStudentToCareer("example@gmail.com", "1")
	if err == nil {
		t.Fatal("test must fail")
	}

	// Then
	require.EqualError(t, err, "student [student_email: example@gmail.com] service: career already assigned")
}

func TestService_AssignStudentToCareer_MaxCareerReachedError(t *testing.T) {
	// Given
	storage_ := storageMock{}
	storage_.On("GetStudentCareerIDs", "example@gmail.com").Return([]int{2, 3, 4}, nil)

	s := NewService(&storage_)

	// When
	err := s.AssignStudentToCareer("example@gmail.com", "1")
	if err == nil {
		t.Fatal("test must fail")
	}

	// Then
	require.EqualError(t, err, "student [student_email: example@gmail.com] service: student already has maximum careers assigned")
}

func TestService_AssignStudentToCareer_StorageAssignStudentToCareerError(t *testing.T) {
	// Given
	storage_ := storageMock{}
	storage_.On("GetStudentCareerIDs", "example@gmail.com").Return([]int{}, nil)
	storage_.On("AssignStudentToCareer", "example@gmail.com", "1").Return(errors.New("error"))

	s := NewService(&storage_)

	// When
	err := s.AssignStudentToCareer("example@gmail.com", "1")
	if err == nil {
		t.Fatal("test must fail")
	}

	// Then
	require.EqualError(t, err, "could not assign student [student_email: example@gmail.com] to career: error")
}

func TestService_AssignStudentToCareer_StorageAssignStudentToCareerNotFoundError(t *testing.T) {
	// Given
	storage_ := storageMock{}
	storage_.On("GetStudentCareerIDs", "example@gmail.com").Return([]int{}, nil)
	storage_.On("AssignStudentToCareer", "example@gmail.com", "1").Return(storage.ErrNotFound)

	s := NewService(&storage_)

	// When
	err := s.AssignStudentToCareer("example@gmail.com", "1")
	if err == nil {
		t.Fatal("test must fail")
	}

	// Then
	require.EqualError(t, err, "could not assign student [student_email: example@gmail.com] to career: service: resource not found")
}

func stringToPtr(s string) *string {
	if s == "" {
		return nil
	}

	return &s
}

func TestService_UpdateStudentSubject(t *testing.T) {
	// Given
	storage_ := storageMock{}
	storage_.On("UpdateStudentSubject", storage.UpdateStudentSubjectRequest{
		StudentEmail: "test@gmail.com",
		CareerID:     "1",
		SubjectID:    "2",
		Status:       "APROBADA",
		Description:  stringToPtr("Aprobe!"),
	}).Return(nil)

	s := NewService(&storage_)

	// When
	err := s.UpdateStudentSubject(UpdateStudentSubjectRequest{
		StudentEmail: "test@gmail.com",
		CareerID:     "1",
		SubjectID:    "2",
		Status:       "APROBADA",
		Description:  "Aprobe!",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Then
	require.Nil(t, err)
}

func TestService_UpdateStudentSubject_StorageError(t *testing.T) {
	// Given
	storage_ := storageMock{}
	storage_.On("UpdateStudentSubject", storage.UpdateStudentSubjectRequest{
		StudentEmail: "test@gmail.com",
		CareerID:     "1",
		SubjectID:    "2",
		Status:       "APROBADA",
		Description:  stringToPtr("Aprobe!"),
	}).Return(errors.New("error"))

	s := NewService(&storage_)

	// When
	err := s.UpdateStudentSubject(UpdateStudentSubjectRequest{
		StudentEmail: "test@gmail.com",
		CareerID:     "1",
		SubjectID:    "2",
		Status:       "APROBADA",
		Description:  "Aprobe!",
	})
	if err == nil {
		t.Fatal("test must fail")
	}

	// Then
	require.EqualError(t, err, "could not update subject: error")
}

func TestService_UpdateStudentSubject_StorageNotFoundError(t *testing.T) {
	// Given
	storage_ := storageMock{}
	storage_.On("UpdateStudentSubject", storage.UpdateStudentSubjectRequest{
		StudentEmail: "test@gmail.com",
		CareerID:     "1",
		SubjectID:    "2",
		Status:       "APROBADA",
		Description:  stringToPtr("Aprobe!"),
	}).Return(storage.ErrNotFound)

	s := NewService(&storage_)

	// When
	err := s.UpdateStudentSubject(UpdateStudentSubjectRequest{
		StudentEmail: "test@gmail.com",
		CareerID:     "1",
		SubjectID:    "2",
		Status:       "APROBADA",
		Description:  "Aprobe!",
	})
	if err == nil {
		t.Fatal("test must fail")
	}

	// Then
	require.EqualError(t, err, "could not update subject: service: resource not found")
}

func TestService_UpdateStudentSubject_NilDescription(t *testing.T) {
	// Given
	storage_ := storageMock{}
	storage_.On("UpdateStudentSubject", storage.UpdateStudentSubjectRequest{
		StudentEmail: "test@gmail.com",
		CareerID:     "1",
		SubjectID:    "2",
		Status:       "APROBADA",
	}).Return(nil)

	s := NewService(&storage_)

	// When
	err := s.UpdateStudentSubject(UpdateStudentSubjectRequest{
		StudentEmail: "test@gmail.com",
		CareerID:     "1",
		SubjectID:    "2",
		Status:       "APROBADA",
		Description:  "",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Then
	require.Nil(t, err)
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
	require.EqualError(t, err, "could not get subjects: error")
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
	require.EqualError(t, err, "could not get subjects: service: resource not found")
}

func TestService_GetSubjectDetails(t *testing.T) {
	// Given
	storage_ := storageMock{}
	storage_.On("GetSubjectDetails", "1", "2").Return(storage.SubjectDetails{
		ID:     1,
		Hours:  nil,
		Points: nil,
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
	require.Equal(t, []byte(`{"id":1,"name":"Algebra","type":"REQUIRED","uri":null,"meet":null,"hours":null,"points":null}`), subject)
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
	require.EqualError(t, err, "could not get subject details: error")
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
	require.EqualError(t, err, "could not get subject details: service: resource not found")
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
	require.EqualError(t, err, "could not get professorships: error")
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
	require.EqualError(t, err, "could not get professorships: service: resource not found")
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

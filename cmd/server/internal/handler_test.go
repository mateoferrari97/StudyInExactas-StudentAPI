package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/mateoferrari97/AnitiMonono-StudentAPI/cmd/server/internal/service"
	"github.com/mateoferrari97/Kit/web/server"
)

type wrapperMock struct {
	f server.HandlerFunc
}

func (w *wrapperMock) Wrap(_, _ string, f server.HandlerFunc, _ ...server.Middleware) {
	w.f = f
}

type serviceMock struct {
	mock.Mock
}

func (s *serviceMock) CreateStudent(name, studentEmail string) error {
	return s.Called(name, studentEmail).Error(0)
}

func (s *serviceMock) AssignStudentToCareer(studentEmail, careerID string) error {
	args := s.Called(studentEmail, careerID)
	return args.Error(0)
}

func (s *serviceMock) GetStudentSubjects(studentEmail, careerID string) ([]byte, error) {
	args := s.Called(studentEmail, careerID)
	return args.Get(0).([]byte), args.Error(1)
}

func (s *serviceMock) UpdateStudentSubject(req service.UpdateStudentSubjectRequest) error {
	args := s.Called(req)
	return args.Error(0)
}

func (s *serviceMock) GetSubjectDetails(subjectID, careerID string) ([]byte, error) {
	args := s.Called(subjectID, careerID)
	return args.Get(0).([]byte), args.Error(1)
}

func (s *serviceMock) GetProfessorships(subjectID, careerID string) ([]byte, error) {
	args := s.Called(subjectID, careerID)
	return args.Get(0).([]byte), args.Error(1)
}

func TestHandler_CreateStudent(t *testing.T) {
	// Given
	wrapper := wrapperMock{}
	service_ := serviceMock{}
	service_.On("CreateStudent", "example", "example@gmail.com").Return(nil)

	h := NewHandler(&wrapper, &service_)
	h.CreateStudent()

	b := bytes.NewReader([]byte(`{
		"name": "example",
		"student_email": "example@gmail.com"
	}`))
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "whocares", b)

	// When
	err := wrapper.f(w, r)
	if err != nil {
		t.Fatal(err)
	}

	// Then
	require.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_CreateStudent_BodyValidationError(t *testing.T) {
	tt := []struct {
		name          string
		b             io.Reader
		expectedError string
	}{
		{
			name: "name is empty",
			b: bytes.NewReader([]byte(`{
				"name": "",
				"student_email": "example@gmail.com"
			}`)),
			expectedError: "400 bad_request: Key: 'Name' Error:Field validation for 'Name' failed on the 'required' tag",
		},
		{
			name: "student email is empty",
			b: bytes.NewReader([]byte(`{
				"name": "example",
				"student_email": ""
			}`)),
			expectedError: "400 bad_request: Key: 'StudentEmail' Error:Field validation for 'StudentEmail' failed on the 'required' tag",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// Given
			wrapper := wrapperMock{}
			service_ := serviceMock{}

			h := NewHandler(&wrapper, &service_)
			h.CreateStudent()

			w := httptest.NewRecorder()
			r, _ := http.NewRequest("POST", "whocares", tc.b)

			// When
			err := wrapper.f(w, r)
			if err == nil {
				t.Fatal("test must fail")
			}

			// Then
			require.EqualError(t, err, tc.expectedError)
		})
	}
}

func TestHandler_CreateStudent_BodyError(t *testing.T) {
	// Given
	wrapper := wrapperMock{}
	service_ := serviceMock{}

	h := NewHandler(&wrapper, &service_)
	h.CreateStudent()

	b := bytes.NewReader([]byte(`{
			"name": 1,
			"student_email": "example@gmail.com"
		}`))
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "whocares", b)

	// When
	err := wrapper.f(w, r)
	if err == nil {
		t.Fatal("test must fail")
	}

	// Then
	require.EqualError(t, err, "422 unprocessable_entity: json: cannot unmarshal number into Go struct field .name of type string")
}

func TestHandler_CreateStudent_StudentAlreadyExist(t *testing.T) {
	// Given
	wrapper := wrapperMock{}
	service_ := serviceMock{}
	service_.On("CreateStudent", "example", "example@gmail.com").Return(service.ErrStudentAlreadyExist)

	h := NewHandler(&wrapper, &service_)
	h.CreateStudent()

	b := bytes.NewReader([]byte(`{
		"name": "example",
		"student_email": "example@gmail.com"
	}`))
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "whocares", b)

	// When
	err := wrapper.f(w, r)
	if err == nil {
		t.Fatal("test must fail")
	}

	// Then
	require.EqualError(t, err, "409 conflict: service: student already exist")
}

func TestHandler_CreateStudent_ServiceError(t *testing.T) {

	// Given
	wrapper := wrapperMock{}
	service_ := serviceMock{}
	service_.On("CreateStudent", "example", "example@gmail.com").Return(errors.New("error"))

	h := NewHandler(&wrapper, &service_)
	h.CreateStudent()

	b := bytes.NewReader([]byte(`{
				"name": "example",
				"student_email": "example@gmail.com"
			}`))
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "whocares", b)

	// When
	err := wrapper.f(w, r)
	if err == nil {
		t.Fatal("test must fail")
	}

	// Then
	require.EqualError(t, err, "error")
}

func TestHandler_AssignStudentToCareer(t *testing.T) {
	// Given
	wrapper := wrapperMock{}
	service_ := serviceMock{}
	service_.On("AssignStudentToCareer", "example@gmail.com", "1").Return(nil)

	h := NewHandler(&wrapper, &service_)
	h.AssignStudentToCareer()

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "whocares", nil)
	r = mux.SetURLVars(r, map[string]string{
		"studentEmail": "example@gmail.com",
		"careerID":     "1",
	})

	// When
	err := wrapper.f(w, r)
	if err != nil {
		t.Fatal(err)
	}

	// Then
	require.Nil(t, err)
}

func TestHandler_AssignStudentToCareer_ParamsError(t *testing.T) {
	type expectedError struct {
		StatusCode int
		Code       string
		Message    string
	}

	tt := []struct {
		name        string
		params      map[string]string
		expectedErr expectedError
	}{
		{
			name: "student email is missing",
			params: map[string]string{
				"studentEmail": "",
				"careerID":     "1",
			},
			expectedErr: expectedError{
				StatusCode: http.StatusBadRequest,
				Code:       "bad_request",
				Message:    "student email is required",
			},
		},
		{
			name: "career id is missing",
			params: map[string]string{
				"studentEmail": "example@gmail.com",
				"careerID":     "",
			},
			expectedErr: expectedError{
				StatusCode: http.StatusBadRequest,
				Code:       "bad_request",
				Message:    "career id is required",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// Given
			wrapper := wrapperMock{}
			service_ := serviceMock{}
			h := NewHandler(&wrapper, &service_)
			h.AssignStudentToCareer()

			w := httptest.NewRecorder()
			r, _ := http.NewRequest("POST", "whocares", nil)
			r = mux.SetURLVars(r, tc.params)

			// When
			err := wrapper.f(w, r)
			if err == nil {
				t.Fatal("test must fail")
			}

			// Then
			hErr := err.(*server.Error)
			require.Equal(t, tc.expectedErr.StatusCode, hErr.StatusCode)
			require.Equal(t, tc.expectedErr.Code, hErr.Code)
			require.Equal(t, tc.expectedErr.Message, hErr.Message)
		})
	}
}

func TestHandler_AssignStudentToCareer_ServiceError(t *testing.T) {
	tt := []struct {
		name          string
		expectedError string
		returnedError error
	}{
		{
			name:          "unknown error",
			expectedError: "error",
			returnedError: errors.New("error"),
		},
		{
			name:          "not found error",
			expectedError: "404 not_found: service: resource not found",
			returnedError: service.ErrNotFound,
		},
		{
			name:          "career already assigned error",
			expectedError: "409 conflict: service: career already assigned",
			returnedError: service.ErrCareerAlreadyAssigned,
		},
		{
			name:          "max career reached error",
			expectedError: "409 conflict: service: student already has maximum careers assigned",
			returnedError: service.ErrMaxCareerReached,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// Given
			wrapper := wrapperMock{}
			service_ := serviceMock{}
			service_.On("AssignStudentToCareer", "example@gmail.com", "1").Return(tc.returnedError)

			h := NewHandler(&wrapper, &service_)
			h.AssignStudentToCareer()

			w := httptest.NewRecorder()
			r, _ := http.NewRequest("POST", "whocares", nil)
			r = mux.SetURLVars(r, map[string]string{
				"studentEmail": "example@gmail.com",
				"careerID":     "1",
			})

			// When
			err := wrapper.f(w, r)
			if err == nil {
				t.Fatal("test must fail")
			}

			// Then
			require.EqualError(t, err, tc.expectedError)
		})
	}
}

func TestHandler_GetStudentSubjects(t *testing.T) {
	// Given
	wrapper := wrapperMock{}
	service_ := serviceMock{}
	service_.On("GetStudentSubjects", "example@gmail.com", "1").Return([]byte(`{
		"correlatives": null,
		"subjects": null
	}`), nil)

	h := NewHandler(&wrapper, &service_)
	h.GetStudentSubjects()

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "whocares", nil)
	r = mux.SetURLVars(r, map[string]string{
		"studentEmail": "example@gmail.com",
		"careerID":     "1",
	})

	// When
	err := wrapper.f(w, r)
	if err != nil {
		t.Fatal(err)
	}

	// Then
	var responseBody struct {
		Correlatives interface{} `json:"correlatives"`
		Subjects     interface{} `json:"subjects"`
	}

	if err := json.NewDecoder(w.Body).Decode(&responseBody); err != nil {
		t.Fatal(err)
	}

	require.Equal(t, http.StatusOK, w.Code)
	require.Nil(t, responseBody.Correlatives)
	require.Nil(t, responseBody.Subjects)
}

func TestHandler_GetStudentSubjects_ParamsError(t *testing.T) {
	type expectedError struct {
		StatusCode int
		Code       string
		Message    string
	}

	tt := []struct {
		name        string
		params      map[string]string
		expectedErr expectedError
	}{
		{
			name: "student email is missing",
			params: map[string]string{
				"studentEmail": "",
				"careerID":     "1",
			},
			expectedErr: expectedError{
				StatusCode: http.StatusBadRequest,
				Code:       "bad_request",
				Message:    "student email is required",
			},
		},
		{
			name: "career id is missing",
			params: map[string]string{
				"studentEmail": "example@gmail.com",
				"careerID":     "",
			},
			expectedErr: expectedError{
				StatusCode: http.StatusBadRequest,
				Code:       "bad_request",
				Message:    "career id is required",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// Given
			wrapper := wrapperMock{}
			service_ := serviceMock{}
			h := NewHandler(&wrapper, &service_)
			h.GetStudentSubjects()

			w := httptest.NewRecorder()
			r, _ := http.NewRequest("GET", "whocares", nil)
			r = mux.SetURLVars(r, tc.params)

			// When
			err := wrapper.f(w, r)
			if err == nil {
				t.Fatal("test must fail")
			}

			// Then
			hErr := err.(*server.Error)
			require.Equal(t, tc.expectedErr.StatusCode, hErr.StatusCode)
			require.Equal(t, tc.expectedErr.Code, hErr.Code)
			require.Equal(t, tc.expectedErr.Message, hErr.Message)
		})
	}
}

func TestHandler_GetStudentSubjects_ServiceError(t *testing.T) {
	// Given
	wrapper := wrapperMock{}
	service_ := serviceMock{}
	service_.On("GetStudentSubjects", "example@gmail.com", "1").Return([]byte{}, errors.New("error"))

	h := NewHandler(&wrapper, &service_)
	h.GetStudentSubjects()

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "whocares", nil)
	r = mux.SetURLVars(r, map[string]string{
		"studentEmail": "example@gmail.com",
		"careerID":     "1",
	})

	// When
	err := wrapper.f(w, r)
	if err == nil {
		t.Fatal("test must fail")
	}

	// Then
	require.EqualError(t, err, "error")
}

func TestHandler_GetStudentSubjects_ServiceNotFoundError(t *testing.T) {
	// Given
	wrapper := wrapperMock{}
	service_ := serviceMock{}
	service_.On("GetStudentSubjects", "example@gmail.com", "1").Return([]byte{}, service.ErrNotFound)

	h := NewHandler(&wrapper, &service_)
	h.GetStudentSubjects()

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "whocares", nil)
	r = mux.SetURLVars(r, map[string]string{
		"studentEmail": "example@gmail.com",
		"careerID":     "1",
	})

	// When
	err := wrapper.f(w, r)
	if err == nil {
		t.Fatal("test must fail")
	}

	// Then
	hErr := err.(*server.Error)
	require.Equal(t, http.StatusNotFound, hErr.StatusCode)
	require.Equal(t, "not_found", hErr.Code)
	require.Equal(t, "service: resource not found", hErr.Message)
}

func TestHandler_UpdateStudentSubject(t *testing.T) {
	// Given
	wrapper := wrapperMock{}
	service_ := serviceMock{}

	service_.On("UpdateStudentSubject", service.UpdateStudentSubjectRequest{
		StudentEmail: "test@gmail.com",
		CareerID:     "2",
		SubjectID:    "1",
		Status:       "APROBADA",
		Description:  "Aprobé!",
	}).Return(nil)

	h := NewHandler(&wrapper, &service_)
	h.UpdateStudentSubject()

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PUT", "whocares", bytes.NewReader([]byte(`{"status":"APROBADA","description":"Aprobé!"}`)))
	r = mux.SetURLVars(r, map[string]string{
		"studentEmail": "test@gmail.com",
		"careerID":     "2",
		"subjectID":    "1",
	})

	// When
	err := wrapper.f(w, r)
	if err != nil {
		t.Fatal(err)
	}

	// Then
	require.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_UpdateStudentSubject_ParamsError(t *testing.T) {
	type expectedError struct {
		StatusCode int
		Code       string
		Message    string
	}

	tt := []struct {
		name        string
		params      map[string]string
		expectedErr expectedError
	}{
		{
			name: "student email is missing",
			params: map[string]string{
				"studentEmail": "",
				"careerID":     "2",
				"subjectID":    "1",
			},
			expectedErr: expectedError{
				StatusCode: http.StatusBadRequest,
				Code:       "bad_request",
				Message:    "student email is required",
			},
		},
		{
			name: "career id is missing",
			params: map[string]string{
				"studentEmail": "test@gmail.com",
				"careerID":     "",
				"subjectID":    "1",
			},
			expectedErr: expectedError{
				StatusCode: http.StatusBadRequest,
				Code:       "bad_request",
				Message:    "career id is required",
			},
		},
		{
			name: "subject id is missing",
			params: map[string]string{
				"studentEmail": "test@gmail.com",
				"careerID":     "2",
				"subjectID":    "",
			},
			expectedErr: expectedError{
				StatusCode: http.StatusBadRequest,
				Code:       "bad_request",
				Message:    "subject id is required",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// Given
			wrapper := wrapperMock{}

			h := NewHandler(&wrapper, nil)
			h.UpdateStudentSubject()

			w := httptest.NewRecorder()
			r, _ := http.NewRequest("PUT", "whocares", bytes.NewReader([]byte(`{"status":"APROBADA","description":"Aprobé!"}`)))
			r = mux.SetURLVars(r, tc.params)

			// When
			err := wrapper.f(w, r)
			if err == nil {
				t.Fatal("test must fail")
			}

			// Then
			hErr := err.(*server.Error)
			require.Equal(t, tc.expectedErr.StatusCode, hErr.StatusCode)
			require.Equal(t, tc.expectedErr.Code, hErr.Code)
			require.Equal(t, tc.expectedErr.Message, hErr.Message)
		})
	}
}

func TestHandler_UpdateStudentSubject_BodyError(t *testing.T) {
	// Given
	wrapper := wrapperMock{}

	h := NewHandler(&wrapper, nil)
	h.UpdateStudentSubject()

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PUT", "whocares", bytes.NewReader([]byte(`{"status":1,"description":"Aprobé!"}`)))
	r = mux.SetURLVars(r, map[string]string{
		"studentEmail": "test@gmail.com",
		"careerID":     "2",
		"subjectID":    "1",
	})

	// When
	err := wrapper.f(w, r)
	if err == nil {
		t.Fatal("test must fail")
	}

	// Then
	hErr := err.(*server.Error)
	require.Equal(t, 422, hErr.StatusCode)
	require.Equal(t, "unprocessable_entity", hErr.Code)
	require.Equal(t, "json: cannot unmarshal number into Go struct field .status of type string", hErr.Message)
}

func TestHandler_UpdateStudentSubject_BodyValidationError(t *testing.T) {
	type expectedError struct {
		StatusCode int
		Code       string
		Message    string
	}

	tt := []struct {
		name          string
		body          io.Reader
		expectedError expectedError
	}{
		{
			name: "status is missing",
			body: bytes.NewReader([]byte(`{"status":"","description":"Aprobé!"}`)),
			expectedError: expectedError{
				StatusCode: 400,
				Code:       "bad_request",
				Message:    "Key: 'Status' Error:Field validation for 'Status' failed on the 'required' tag",
			},
		},
		{
			name: "status value is invalid",
			body: bytes.NewReader([]byte(`{"status":"INVALID","description":"Aprobé!"}`)),
			expectedError: expectedError{
				StatusCode: 400,
				Code:       "bad_request",
				Message:    "Key: 'Status' Error:Field validation for 'Status' failed on the 'oneof' tag",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// Given
			wrapper := wrapperMock{}

			h := NewHandler(&wrapper, nil)
			h.UpdateStudentSubject()

			w := httptest.NewRecorder()
			r, _ := http.NewRequest("PUT", "whocares", tc.body)
			r = mux.SetURLVars(r, map[string]string{
				"studentEmail": "test@gmail.com",
				"careerID":     "2",
				"subjectID":    "1",
			})

			// When
			err := wrapper.f(w, r)
			if err == nil {
				t.Fatal("test must fail")
			}

			// Then
			hErr := err.(*server.Error)
			require.Equal(t, tc.expectedError.StatusCode, hErr.StatusCode)
			require.Equal(t, tc.expectedError.Code, hErr.Code)
			require.Equal(t, tc.expectedError.Message, hErr.Message)
		})
	}
}

func TestHandler_UpdateStudentSubject_ServiceError(t *testing.T) {
	// Given
	wrapper := wrapperMock{}
	service_ := serviceMock{}
	service_.On("UpdateStudentSubject", service.UpdateStudentSubjectRequest{
		StudentEmail: "test@gmail.com",
		CareerID:     "2",
		SubjectID:    "1",
		Status:       "APROBADA",
		Description:  "Aprobé!",
	}).Return(errors.New("error"))

	h := NewHandler(&wrapper, &service_)
	h.UpdateStudentSubject()

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PUT", "whocares", bytes.NewReader([]byte(`{"status":"APROBADA","description":"Aprobé!"}`)))
	r = mux.SetURLVars(r, map[string]string{
		"studentEmail": "test@gmail.com",
		"careerID":     "2",
		"subjectID":    "1",
	})

	// When
	err := wrapper.f(w, r)
	if err == nil {
		t.Fatal("test must fail")
	}

	// Then
	require.EqualError(t, err, "error")
}

func TestHandler_UpdateStudentSubject_ServiceNotFoundError(t *testing.T) {
	// Given
	wrapper := wrapperMock{}
	service_ := serviceMock{}
	service_.On("UpdateStudentSubject", service.UpdateStudentSubjectRequest{
		StudentEmail: "test@gmail.com",
		CareerID:     "2",
		SubjectID:    "1",
		Status:       "APROBADA",
		Description:  "Aprobé!",
	}).Return(service.ErrNotFound)

	h := NewHandler(&wrapper, &service_)
	h.UpdateStudentSubject()

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PUT", "whocares", bytes.NewReader([]byte(`{"status":"APROBADA","description":"Aprobé!"}`)))
	r = mux.SetURLVars(r, map[string]string{
		"studentEmail": "test@gmail.com",
		"careerID":     "2",
		"subjectID":    "1",
	})

	// When
	err := wrapper.f(w, r)
	if err == nil {
		t.Fatal("test must fail")
	}

	// Then
	hErr := err.(*server.Error)
	require.Equal(t, http.StatusNotFound, hErr.StatusCode)
	require.Equal(t, "not_found", hErr.Code)
	require.Equal(t, "service: resource not found", hErr.Message)
}

func TestHandler_GetSubjectDetails(t *testing.T) {
	// Given
	wrapper := wrapperMock{}
	service_ := serviceMock{}
	service_.On("GetSubjectDetails", "1", "2").Return([]byte(`{"id": 3}`), nil)

	h := NewHandler(&wrapper, &service_)
	h.GetSubjectDetails()

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "whocares", nil)
	r = mux.SetURLVars(r, map[string]string{
		"subjectID": "1",
		"careerID":  "2",
	})

	// When
	err := wrapper.f(w, r)
	if err != nil {
		t.Fatal(err)
	}

	// Then
	var responseBody struct {
		SubjectID int `json:"id"`
	}

	if err := json.NewDecoder(w.Body).Decode(&responseBody); err != nil {
		t.Fatal(err)
	}

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, 3, responseBody.SubjectID)
}

func TestHandler_GetSubjectDetails_ParamsError(t *testing.T) {
	type expectedError struct {
		StatusCode int
		Code       string
		Message    string
	}

	tt := []struct {
		name        string
		params      map[string]string
		expectedErr expectedError
	}{
		{
			name: "subject id is missing",
			params: map[string]string{
				"subjectID": "",
				"careerID":  "1",
			},
			expectedErr: expectedError{
				StatusCode: http.StatusBadRequest,
				Code:       "bad_request",
				Message:    "subject id is required",
			},
		},
		{
			name: "career id is missing",
			params: map[string]string{
				"subjectID": "1",
				"careerID":  "",
			},
			expectedErr: expectedError{
				StatusCode: http.StatusBadRequest,
				Code:       "bad_request",
				Message:    "career id is required",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// Given
			wrapper := wrapperMock{}
			service_ := serviceMock{}
			h := NewHandler(&wrapper, &service_)
			h.GetSubjectDetails()

			w := httptest.NewRecorder()
			r, _ := http.NewRequest("GET", "whocares", nil)
			r = mux.SetURLVars(r, tc.params)

			// When
			err := wrapper.f(w, r)
			if err == nil {
				t.Fatal("test must fail")
			}

			// Then
			hErr := err.(*server.Error)
			require.Equal(t, tc.expectedErr.StatusCode, hErr.StatusCode)
			require.Equal(t, tc.expectedErr.Code, hErr.Code)
			require.Equal(t, tc.expectedErr.Message, hErr.Message)
		})
	}
}

func TestHandler_GetSubjectDetails_ServiceError(t *testing.T) {
	// Given
	wrapper := wrapperMock{}
	service_ := serviceMock{}
	service_.On("GetSubjectDetails", "1", "2").Return([]byte{}, errors.New("error"))

	h := NewHandler(&wrapper, &service_)
	h.GetSubjectDetails()

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "whocares", nil)
	r = mux.SetURLVars(r, map[string]string{
		"subjectID": "1",
		"careerID":  "2",
	})

	// When
	err := wrapper.f(w, r)
	if err == nil {
		t.Fatal("test must fail")
	}

	// Then
	require.EqualError(t, err, "error")
}

func TestHandler_GetSubjectDetails_ServiceNotFoundError(t *testing.T) {
	// Given
	wrapper := wrapperMock{}
	service_ := serviceMock{}
	service_.On("GetSubjectDetails", "1", "2").Return([]byte{}, service.ErrNotFound)

	h := NewHandler(&wrapper, &service_)
	h.GetSubjectDetails()

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "whocares", nil)
	r = mux.SetURLVars(r, map[string]string{
		"subjectID": "1",
		"careerID":  "2",
	})

	// When
	err := wrapper.f(w, r)
	if err == nil {
		t.Fatal("test must fail")
	}

	// Then
	hErr := err.(*server.Error)
	require.Equal(t, http.StatusNotFound, hErr.StatusCode)
	require.Equal(t, "not_found", hErr.Code)
	require.Equal(t, "service: resource not found", hErr.Message)
}

func TestHandler_GetProfessorships(t *testing.T) {
	// Given
	wrapper := wrapperMock{}
	service_ := serviceMock{}
	service_.On("GetProfessorships", "1", "2").Return([]byte(`{"professorship": null}`), nil)

	h := NewHandler(&wrapper, &service_)
	h.GetProfessorships()

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "whocares", nil)
	r = mux.SetURLVars(r, map[string]string{
		"subjectID": "1",
		"careerID":  "2",
	})

	// When
	err := wrapper.f(w, r)
	if err != nil {
		t.Fatal(err)
	}

	// Then
	var responseBody struct {
		Professorship interface{} `json:"professorship"`
	}

	if err := json.NewDecoder(w.Body).Decode(&responseBody); err != nil {
		t.Fatal(err)
	}

	require.Equal(t, http.StatusOK, w.Code)
	require.Nil(t, responseBody.Professorship)
}

func TestHandler_GetProfessorships_ParamsError(t *testing.T) {
	type expectedError struct {
		StatusCode int
		Code       string
		Message    string
	}

	tt := []struct {
		name        string
		params      map[string]string
		expectedErr expectedError
	}{
		{
			name: "subject id is missing",
			params: map[string]string{
				"subjectID": "",
				"careerID":  "1",
			},
			expectedErr: expectedError{
				StatusCode: http.StatusBadRequest,
				Code:       "bad_request",
				Message:    "subject id is required",
			},
		},
		{
			name: "career id is missing",
			params: map[string]string{
				"subjectID": "1",
				"careerID":  "",
			},
			expectedErr: expectedError{
				StatusCode: http.StatusBadRequest,
				Code:       "bad_request",
				Message:    "career id is required",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// Given
			wrapper := wrapperMock{}
			service_ := serviceMock{}
			h := NewHandler(&wrapper, &service_)
			h.GetProfessorships()

			w := httptest.NewRecorder()
			r, _ := http.NewRequest("GET", "whocares", nil)
			r = mux.SetURLVars(r, tc.params)

			// When
			err := wrapper.f(w, r)
			if err == nil {
				t.Fatal("test must fail")
			}

			// Then
			hErr := err.(*server.Error)
			require.Equal(t, tc.expectedErr.StatusCode, hErr.StatusCode)
			require.Equal(t, tc.expectedErr.Code, hErr.Code)
			require.Equal(t, tc.expectedErr.Message, hErr.Message)
		})
	}
}

func TestHandler_GetProfessorships_ServiceError(t *testing.T) {
	// Given
	wrapper := wrapperMock{}
	service_ := serviceMock{}
	service_.On("GetProfessorships", "1", "2").Return([]byte{}, errors.New("error"))

	h := NewHandler(&wrapper, &service_)
	h.GetProfessorships()

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "whocares", nil)
	r = mux.SetURLVars(r, map[string]string{
		"subjectID": "1",
		"careerID":  "2",
	})

	// When
	err := wrapper.f(w, r)
	if err == nil {
		t.Fatal("test must fail")
	}

	// Then
	require.EqualError(t, err, "error")
}

func TestHandler_GetProfessorships_ServiceNotFoundError(t *testing.T) {
	// Given
	wrapper := wrapperMock{}
	service_ := serviceMock{}
	service_.On("GetProfessorships", "1", "2").Return([]byte{}, service.ErrNotFound)

	h := NewHandler(&wrapper, &service_)
	h.GetProfessorships()

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "whocares", nil)
	r = mux.SetURLVars(r, map[string]string{
		"subjectID": "1",
		"careerID":  "2",
	})

	// When
	err := wrapper.f(w, r)
	if err == nil {
		t.Fatal("test must fail")
	}

	// Then
	hErr := err.(*server.Error)
	require.Equal(t, http.StatusNotFound, hErr.StatusCode)
	require.Equal(t, "not_found", hErr.Code)
	require.Equal(t, "service: resource not found", hErr.Message)
}

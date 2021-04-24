package internal

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/mateoferrari97/Kit/web/server"
	"github.com/mateoferrari97/my-path/cmd/app/internal/service"
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

func (s *serviceMock) GetStudentSubjects(studentEmail, careerID string) ([]byte, error) {
	args := s.Called(studentEmail, careerID)
	return args.Get(0).([]byte), args.Error(1)
}

func (s *serviceMock) GetSubjectDetails(subjectID, careerID string) ([]byte, error) {
	args := s.Called(subjectID, careerID)
	return args.Get(0).([]byte), args.Error(1)
}

func (s *serviceMock) GetProfessorships(subjectID, careerID string) ([]byte, error) {
	args := s.Called(subjectID, careerID)
	return args.Get(0).([]byte), args.Error(1)
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

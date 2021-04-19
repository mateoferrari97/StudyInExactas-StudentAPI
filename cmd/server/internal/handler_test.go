package internal

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mateoferrari97/my-path/internal/server"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

type wrapperMock struct {
	router *mux.Router
}

func (w *wrapperMock) Wrap(method, pattern string, f server.HandlerFunc) {
	wrapH := func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err == nil {
			return
		}

		hErr := err.(*server.Error)
		w.WriteHeader(hErr.StatusCode)
		_, _ = w.Write([]byte(hErr.Message))
	}

	w.router.HandleFunc(pattern, wrapH).Methods(method)
}

type serviceMock struct {
	mock.Mock
}

func (s *serviceMock) GetStudentSubjects(studentEmail, careerID string) ([]byte, error) {
	args := s.Called(studentEmail, careerID)
	return args.Get(0).([]byte), args.Error(1)
}

func (s *serviceMock) GetSubjectDetails(subjectID, careerID string) ([]byte, error) {
	panic("implement me")
}

func (s *serviceMock) GetProfessorships(subjectID, careerID string) ([]byte, error) {
	panic("implement me")
}

func TestHandler_GetStudentSubjects(t *testing.T) {
	// Given
	wrapper := wrapperMock{router: mux.NewRouter()}
	service := serviceMock{}
	service.On("GetStudentSubjects", "example@gmail.com", "1").Return([]byte(`{}`), nil)

	h := NewHandler(&wrapper, &service)
	h.GetStudentSubjects()

	ts := httptest.NewServer(wrapper.router)
	defer ts.Close()

	// When
	resp, err := http.Get(fmt.Sprintf("%s/students/example@gmail.com/careers/1/subjects", ts.URL))
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	// Then
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

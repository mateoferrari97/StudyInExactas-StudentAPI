package storage

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestStorage_GetStudentSubjects(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("could not start sql mock: %v", err)
	}

	defer db.Close()

	storage_ := NewStorage(sqlx.NewDb(db, ""))

	q := `SELECT cs.subject_id,
       s.name,
       cs.correlative_id,
       cs.type,
       scs.status,
       scs.description
FROM student AS st
         INNER JOIN student_career sc ON sc.student_id = st.id
         INNER JOIN career_subject cs ON cs.career_id = sc.career_id
         INNER JOIN career c ON c.id = cs.career_id
         INNER JOIN student_career_subject scs ON scs.career_subject_id = cs.id AND st.id = scs.student_id
         INNER JOIN subject s on s.id = cs.subject_id
WHERE st.email = ? AND c.id = ?`
	mock.ExpectPrepare(q).WillReturnError(nil)
	mock.ExpectQuery(q).
		WithArgs("example@gmail.com", "1").
		WillReturnError(nil).
		WillReturnRows(
			sqlmock.NewRows([]string{"subject_id", "name", "correlative_id", "type", "status", "description"}).
				AddRow(2, "Subject 2", nil, "REQUIRED", "PENDING", nil))

	// When
	subjects, err := storage_.GetStudentSubjects("example@gmail.com", "1")
	if err != nil {
		t.Fatal(err)
	}

	// Then
	require.NotNil(t, subjects)
	require.Len(t, subjects, 1)

	for _, subject := range subjects {
		require.Equal(t, 2, subject.ID)
		require.Equal(t, "Subject 2", subject.Name)
		require.Equal(t, 0, subject.CorrelativeID)
		require.Equal(t, "REQUIRED", subject.Type)
		require.Equal(t, "PENDING", subject.Status)
		require.Nil(t, subject.Description)
	}
}

func TestStorage_GetStudentSubjects_SubjectHasCorrelative(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("could not start sql mock: %v", err)
	}

	defer db.Close()

	storage_ := NewStorage(sqlx.NewDb(db, ""))

	q := `SELECT cs.subject_id,
       s.name,
       cs.correlative_id,
       cs.type,
       scs.status,
       scs.description
FROM student AS st
         INNER JOIN student_career sc ON sc.student_id = st.id
         INNER JOIN career_subject cs ON cs.career_id = sc.career_id
         INNER JOIN career c ON c.id = cs.career_id
         INNER JOIN student_career_subject scs ON scs.career_subject_id = cs.id AND st.id = scs.student_id
         INNER JOIN subject s on s.id = cs.subject_id
WHERE st.email = ? AND c.id = ?`
	mock.ExpectPrepare(q).WillReturnError(nil)
	mock.ExpectQuery(q).
		WithArgs("example@gmail.com", "1").
		WillReturnError(nil).
		WillReturnRows(
			sqlmock.NewRows([]string{"subject_id", "name", "correlative_id", "type", "status", "description"}).
				AddRow(2, "Subject 2", 1, "REQUIRED", "PENDING", nil))

	// When
	subjects, err := storage_.GetStudentSubjects("example@gmail.com", "1")
	if err != nil {
		t.Fatal(err)
	}

	// Then
	require.NotNil(t, subjects)
	require.Len(t, subjects, 1)

	for _, subject := range subjects {
		require.Equal(t, 2, subject.ID)
		require.Equal(t, "Subject 2", subject.Name)
		require.Equal(t, 1, subject.CorrelativeID)
		require.Equal(t, "REQUIRED", subject.Type)
		require.Equal(t, "PENDING", subject.Status)
		require.Nil(t, subject.Description)
	}
}

func TestStorage_GetStudentSubjects_SubjectHasDescription(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("could not start sql mock: %v", err)
	}

	defer db.Close()

	storage_ := NewStorage(sqlx.NewDb(db, ""))

	q := `SELECT cs.subject_id,
       s.name,
       cs.correlative_id,
       cs.type,
       scs.status,
       scs.description
FROM student AS st
         INNER JOIN student_career sc ON sc.student_id = st.id
         INNER JOIN career_subject cs ON cs.career_id = sc.career_id
         INNER JOIN career c ON c.id = cs.career_id
         INNER JOIN student_career_subject scs ON scs.career_subject_id = cs.id AND st.id = scs.student_id
         INNER JOIN subject s on s.id = cs.subject_id
WHERE st.email = ? AND c.id = ?`
	mock.ExpectPrepare(q).WillReturnError(nil)
	mock.ExpectQuery(q).
		WithArgs("example@gmail.com", "1").
		WillReturnError(nil).
		WillReturnRows(
			sqlmock.NewRows([]string{"subject_id", "name", "correlative_id", "type", "status", "description"}).
				AddRow(2, "Subject 2", nil, "REQUIRED", "PENDING", "..."))

	// When
	subjects, err := storage_.GetStudentSubjects("example@gmail.com", "1")
	if err != nil {
		t.Fatal(err)
	}

	// Then
	require.NotNil(t, subjects)
	require.Len(t, subjects, 1)

	for _, subject := range subjects {
		require.Equal(t, 2, subject.ID)
		require.Equal(t, "Subject 2", subject.Name)
		require.Equal(t, 0, subject.CorrelativeID)
		require.Equal(t, "REQUIRED", subject.Type)
		require.Equal(t, "PENDING", subject.Status)
		require.NotNil(t, subject.Description)
		require.Equal(t, "...", *subject.Description)
	}
}

func TestStorage_GetStudentSubjects_PrepareStmtError(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("could not start sql mock: %v", err)
	}

	defer db.Close()

	storage_ := NewStorage(sqlx.NewDb(db, ""))

	q := `SELECT cs.subject_id,
       s.name,
       cs.correlative_id,
       cs.type,
       scs.status,
       scs.description
FROM student AS st
         INNER JOIN student_career sc ON sc.student_id = st.id
         INNER JOIN career_subject cs ON cs.career_id = sc.career_id
         INNER JOIN career c ON c.id = cs.career_id
         INNER JOIN student_career_subject scs ON scs.career_subject_id = cs.id AND st.id = scs.student_id
         INNER JOIN subject s on s.id = cs.subject_id
WHERE st.email = ? AND c.id = ?`
	mock.ExpectPrepare(q).WillReturnError(errors.New("error"))

	// When
	_, err = storage_.GetStudentSubjects("example@gmail.com", "1")
	if err == nil {
		t.Fatal("test must fail")
	}

	// Then
	require.EqualError(t, err, "error")
}

func TestStorage_GetStudentSubjects_ExecuteStmtError(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("could not start sql mock: %v", err)
	}

	defer db.Close()

	storage_ := NewStorage(sqlx.NewDb(db, ""))

	q := `SELECT cs.subject_id,
       s.name,
       cs.correlative_id,
       cs.type,
       scs.status,
       scs.description
FROM student AS st
         INNER JOIN student_career sc ON sc.student_id = st.id
         INNER JOIN career_subject cs ON cs.career_id = sc.career_id
         INNER JOIN career c ON c.id = cs.career_id
         INNER JOIN student_career_subject scs ON scs.career_subject_id = cs.id AND st.id = scs.student_id
         INNER JOIN subject s on s.id = cs.subject_id
WHERE st.email = ? AND c.id = ?`
	mock.ExpectPrepare(q).WillReturnError(nil)
	mock.ExpectQuery(q).
		WithArgs("example@gmail.com", "1").
		WillReturnError(errors.New("error"))

	// When
	_, err = storage_.GetStudentSubjects("example@gmail.com", "1")
	if err == nil {
		t.Fatal("test must fail")
	}

	// Then
	require.EqualError(t, err, "error")
}

func TestStorage_GetSubjectDetails(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("could not start sql mock: %v", err)
	}

	defer db.Close()

	storage_ := NewStorage(sqlx.NewDb(db, ""))

	q := `SELECT s.id, s.name, s.uri, s.meet, cs.type, cs.hours, cs.points
FROM career_subject cs
         INNER JOIN subject s ON cs.subject_id = s.id
WHERE s.id = ?
  AND career_id = ?
LIMIT 1;`
	mock.ExpectPrepare(q).WillReturnError(nil)
	mock.ExpectQuery(q).
		WithArgs("1", "2").
		WillReturnError(nil).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "hours", "points", "name", "type", "uri", "meet"}).
				AddRow(1, 240, 8, "Subject 1", "REQUIRED", nil, nil))

	// When
	subject, err := storage_.GetSubjectDetails("1", "2")
	if err != nil {
		t.Fatal(err)
	}

	// Then
	require.Equal(t, 1, subject.ID)
	require.Equal(t, 240, subject.Hours)
	require.Equal(t, 8, subject.Points)
	require.Equal(t, "Subject 1", subject.Name)
	require.Equal(t, "REQUIRED", subject.Type)
	require.Nil(t, subject.URI)
	require.Nil(t, subject.Meet)

}

func TestStorage_GetSubjectDetails_PrepareStmtError(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("could not start sql mock: %v", err)
	}

	defer db.Close()

	storage_ := NewStorage(sqlx.NewDb(db, ""))

	q := `SELECT s.id, s.name, s.uri, s.meet, cs.type, cs.hours, cs.points
FROM career_subject cs
         INNER JOIN subject s ON cs.subject_id = s.id
WHERE s.id = ?
  AND career_id = ?
LIMIT 1;`
	mock.ExpectPrepare(q).WillReturnError(errors.New("error"))

	// When
	_, err = storage_.GetSubjectDetails("1", "2")
	if err == nil {
		t.Fatal("test must fail")
	}

	// Then
	require.EqualError(t, err, "error")
}

func TestStorage_GetSubjectDetails_ExecuteStmtError(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("could not start sql mock: %v", err)
	}

	defer db.Close()

	storage_ := NewStorage(sqlx.NewDb(db, ""))

	q := `SELECT s.id, s.name, s.uri, s.meet, cs.type, cs.hours, cs.points
FROM career_subject cs
         INNER JOIN subject s ON cs.subject_id = s.id
WHERE s.id = ?
  AND career_id = ?
LIMIT 1;`
	mock.ExpectPrepare(q).WillReturnError(nil)
	mock.ExpectQuery(q).WillReturnError(errors.New("error"))

	// When
	_, err = storage_.GetSubjectDetails("1", "2")
	if err == nil {
		t.Fatal("test must fail")
	}

	// Then
	require.EqualError(t, err, "error")
}

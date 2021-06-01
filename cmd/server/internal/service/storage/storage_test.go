package storage

import (
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestStorage_CreateStudent(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("could not start sql mock: %v", err)
	}

	defer db.Close()

	storage_ := NewStorage(sqlx.NewDb(db, ""))

	q := `INSERT INTO student (name, email) VALUES (?, ?)`
	mock.ExpectPrepare(q).WillReturnError(nil)
	mock.ExpectExec(q).
		WithArgs("example", "example@gmail.com").
		WillReturnError(nil).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// When
	err = storage_.CreateStudent("example", "example@gmail.com")
	if err != nil {
		t.Fatal(err)
	}

	// Then
	require.Nil(t, err)
}

func TestStorage_CreateStudent_PrepareStmtError(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("could not start sql mock: %v", err)
	}

	defer db.Close()

	storage_ := NewStorage(sqlx.NewDb(db, ""))

	q := `INSERT INTO student (name, email) VALUES (?, ?)`
	mock.ExpectPrepare(q).WillReturnError(errors.New("error"))

	// When
	err = storage_.CreateStudent("example", "example@gmail.com")
	if err == nil {
		t.Fatal("test must fail")
	}

	// Then
	require.EqualError(t, err, "error")
}

func TestStorage_CreateStudent_ExecuteStmtError(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("could not start sql mock: %v", err)
	}

	defer db.Close()

	storage_ := NewStorage(sqlx.NewDb(db, ""))

	q := `INSERT INTO student (name, email) VALUES (?, ?)`
	mock.ExpectPrepare(q).WillReturnError(nil)
	mock.ExpectExec(q).
		WithArgs("example", "example@gmail.com").
		WillReturnError(errors.New("error"))

	// When
	err = storage_.CreateStudent("example", "example@gmail.com")
	if err == nil {
		t.Fatal("test must fail")
	}

	// Then
	require.EqualError(t, err, "error")
}

func TestStorage_CreateStudent_DuplicateEntryError(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("could not start sql mock: %v", err)
	}

	defer db.Close()

	storage_ := NewStorage(sqlx.NewDb(db, ""))

	q := `INSERT INTO student (name, email) VALUES (?, ?)`
	mock.ExpectPrepare(q).WillReturnError(nil)
	mock.ExpectExec(q).
		WithArgs("example", "example@gmail.com").
		WillReturnError(&mysql.MySQLError{Number: 1062})

	// When
	err = storage_.CreateStudent("example", "example@gmail.com")
	if err == nil {
		t.Fatal("test must fail")
	}

	// Then
	require.EqualError(t, err, "storage: resource already exist")
}

func TestStorage_GetStudentCareerIDs(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("could not start sql mock: %v", err)
	}

	defer db.Close()

	storage_ := NewStorage(sqlx.NewDb(db, ""))

	q := `SELECT career_id FROM student_career sc
    INNER JOIN student s ON sc.student_id = s.id
WHERE s.email = ?;`
	mock.ExpectPrepare(q).WillReturnError(nil)
	mock.ExpectQuery(q).
		WithArgs("example@gmail.com").
		WillReturnError(nil).
		WillReturnRows(
			sqlmock.NewRows([]string{"career_id"}).
				AddRow(2))

	// When
	careersIDs, err := storage_.GetStudentCareerIDs("example@gmail.com")
	if err != nil {
		t.Fatal(err)
	}

	// Then
	require.NotNil(t, careersIDs)
	require.Len(t, careersIDs, 1)

	for _, careerID := range careersIDs {
		require.Equal(t, 2, careerID)
	}
}

func TestStorage_GetStudentCareerIDs_PrepareStmtError(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("could not start sql mock: %v", err)
	}

	defer db.Close()

	storage_ := NewStorage(sqlx.NewDb(db, ""))

	q := `SELECT career_id FROM student_career sc
    INNER JOIN student s ON sc.student_id = s.id
WHERE s.email = ?;`
	mock.ExpectPrepare(q).WillReturnError(errors.New("error"))

	// When
	_, err = storage_.GetStudentCareerIDs("example@gmail.com")
	if err == nil {
		t.Fatal("test must fail")
	}

	// Then
	require.EqualError(t, err, "error")
}

func TestStorage_GetStudentCareerIDs_ExecuteStmtError(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("could not start sql mock: %v", err)
	}

	defer db.Close()

	storage_ := NewStorage(sqlx.NewDb(db, ""))

	q := `SELECT career_id FROM student_career sc
    INNER JOIN student s ON sc.student_id = s.id
WHERE s.email = ?;`
	mock.ExpectPrepare(q).WillReturnError(nil)
	mock.ExpectQuery(q).
		WithArgs("example@gmail.com").
		WillReturnError(errors.New("error"))

	// When
	_, err = storage_.GetStudentCareerIDs("example@gmail.com")
	if err == nil {
		t.Fatal("test must fail")
	}

	// Then
	require.EqualError(t, err, "error")
}

func TestStorage_GetStudentCareerIDs_NotFoundError(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("could not start sql mock: %v", err)
	}

	defer db.Close()

	storage_ := NewStorage(sqlx.NewDb(db, ""))

	q := `SELECT career_id FROM student_career sc
    INNER JOIN student s ON sc.student_id = s.id
WHERE s.email = ?;`
	mock.ExpectPrepare(q).WillReturnError(nil)
	mock.ExpectQuery(q).
		WithArgs("example@gmail.com").
		WillReturnError(sql.ErrNoRows)

	// When
	_, err = storage_.GetStudentCareerIDs("example@gmail.com")
	if err == nil {
		t.Fatal("test must fail")
	}

	// Then
	require.EqualError(t, err, "storage: resource not found")
}

func TestStorage_AssignStudentToCareer(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("could not start sql mock: %v", err)
	}

	defer db.Close()

	storage_ := NewStorage(sqlx.NewDb(db, ""))

	mock.ExpectBegin()

	q := `SELECT id FROM student WHERE email = ?;`
	mock.ExpectQuery(q).
		WithArgs("example@gmail.com").
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	q = `SELECT COUNT(1) FROM career WHERE id = ?;`
	mock.ExpectQuery(q).
		WithArgs("1").
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"COUNT(1)"}).AddRow(2))

	q = `INSERT INTO student_career (student_id, career_id) VALUES (?, ?);`
	mock.ExpectExec(q).
		WithArgs(1, "1").
		WillReturnError(nil).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	// When
	err = storage_.AssignStudentToCareer("example@gmail.com", "1")
	if err != nil {
		t.Fatal(err)
	}

	// Then
	require.Nil(t, err)
}

func TestStorage_AssignStudentToCareer_BeginTxError(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("could not start sql mock: %v", err)
	}

	defer db.Close()

	storage_ := NewStorage(sqlx.NewDb(db, ""))

	mock.ExpectBegin().WillReturnError(errors.New("error"))

	// When
	err = storage_.AssignStudentToCareer("example@gmail.com", "1")
	if err == nil {
		t.Fatal("test must fail")
	}

	// Then
	require.EqualError(t, err, "could not begin tx: error")
}

func TestStorage_AssignStudentToCareer_GetStudentIDError(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("could not start sql mock: %v", err)
	}

	defer db.Close()

	storage_ := NewStorage(sqlx.NewDb(db, ""))

	mock.ExpectBegin()

	q := `SELECT id FROM student WHERE email = ?;`
	mock.ExpectQuery(q).
		WithArgs("example@gmail.com").
		WillReturnError(errors.New("error"))

	// When
	err = storage_.AssignStudentToCareer("example@gmail.com", "1")
	if err == nil {
		t.Fatal("test must fail")
	}

	// Then
	require.EqualError(t, err, "error")
}

func TestStorage_AssignStudentToCareer_GetStudentIDNotFound(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("could not start sql mock: %v", err)
	}

	defer db.Close()

	storage_ := NewStorage(sqlx.NewDb(db, ""))

	mock.ExpectBegin()

	q := `SELECT id FROM student WHERE email = ?;`
	mock.ExpectQuery(q).
		WithArgs("example@gmail.com").
		WillReturnError(sql.ErrNoRows)

	mock.ExpectCommit()

	// When
	err = storage_.AssignStudentToCareer("example@gmail.com", "1")
	if err == nil {
		t.Fatal("test must fail")
	}

	// Then
	require.EqualError(t, err, "could not find student: storage: resource not found")
}

func TestStorage_AssignStudentToCareer_FindCareerWithIDError(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("could not start sql mock: %v", err)
	}

	defer db.Close()

	storage_ := NewStorage(sqlx.NewDb(db, ""))

	mock.ExpectBegin()

	q := `SELECT id FROM student WHERE email = ?;`
	mock.ExpectQuery(q).
		WithArgs("example@gmail.com").
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	q = `SELECT COUNT(1) FROM career WHERE id = ?;`
	mock.ExpectQuery(q).
		WithArgs("1").
		WillReturnError(errors.New("error"))

	// When
	err = storage_.AssignStudentToCareer("example@gmail.com", "1")
	if err == nil {
		t.Fatal("test must fail")
	}

	// Then
	require.EqualError(t, err, "error")
}

func TestStorage_AssignStudentToCareer_FindCareerWithIDNotFound(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("could not start sql mock: %v", err)
	}

	defer db.Close()

	storage_ := NewStorage(sqlx.NewDb(db, ""))

	mock.ExpectBegin()

	q := `SELECT id FROM student WHERE email = ?;`
	mock.ExpectQuery(q).
		WithArgs("example@gmail.com").
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	q = `SELECT COUNT(1) FROM career WHERE id = ?;`
	mock.ExpectQuery(q).
		WithArgs("1").
		WillReturnError(sql.ErrNoRows)

	// When
	err = storage_.AssignStudentToCareer("example@gmail.com", "1")
	if err == nil {
		t.Fatal("test must fail")
	}

	// Then
	require.EqualError(t, err, "could not find career: storage: resource not found")
}

func TestStorage_AssignStudentToCareer_CreateStudentWithCareerError(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("could not start sql mock: %v", err)
	}

	defer db.Close()

	storage_ := NewStorage(sqlx.NewDb(db, ""))

	mock.ExpectBegin()

	q := `SELECT id FROM student WHERE email = ?;`
	mock.ExpectQuery(q).
		WithArgs("example@gmail.com").
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	q = `SELECT COUNT(1) FROM career WHERE id = ?;`
	mock.ExpectQuery(q).
		WithArgs("1").
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"COUNT(1)"}).AddRow(2))

	q = `INSERT INTO student_career (student_id, career_id) VALUES (?, ?);`
	mock.ExpectExec(q).
		WithArgs(1, "1").
		WillReturnError(errors.New("error"))

	// When
	err = storage_.AssignStudentToCareer("example@gmail.com", "1")
	if err == nil {
		t.Fatal("test must fail")
	}

	// Then
	require.EqualError(t, err, "error")
}

func TestStorage_AssignStudentToCareer_CommitTxError(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("could not start sql mock: %v", err)
	}

	defer db.Close()

	storage_ := NewStorage(sqlx.NewDb(db, ""))

	mock.ExpectBegin()

	q := `SELECT id FROM student WHERE email = ?;`
	mock.ExpectQuery(q).
		WithArgs("example@gmail.com").
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	q = `SELECT COUNT(1) FROM career WHERE id = ?;`
	mock.ExpectQuery(q).
		WithArgs("1").
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"COUNT(1)"}).AddRow(2))

	q = `INSERT INTO student_career (student_id, career_id) VALUES (?, ?);`
	mock.ExpectExec(q).
		WithArgs(1, "1").
		WillReturnError(nil).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit().WillReturnError(errors.New("error"))

	// When
	err = storage_.AssignStudentToCareer("example@gmail.com", "1")
	if err == nil {
		t.Fatal("test must fail")
	}

	// Then
	require.EqualError(t, err, "could not commit tx: error")
}

func TestStorage_UpdateStudentSubject(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("could not start sql mock: %v", err)
	}

	defer db.Close()

	storage_ := NewStorage(sqlx.NewDb(db, ""))

	mock.ExpectBegin()

	q := `SELECT id FROM student WHERE email = ?;`
	mock.ExpectQuery(q).
		WithArgs("example@gmail.com").
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	q = `SELECT COUNT(1) FROM student_career WHERE student_id = ? AND career_id = ?`
	mock.ExpectQuery(q).
		WithArgs(1, "1").
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"COUNT(1)"}).AddRow(1))

	q = `SELECT id FROM career_subject WHERE career_id = ? AND subject_id = ?`
	mock.ExpectQuery(q).
		WithArgs("1", "1").
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))

	q = `INSERT INTO student_career_subject (student_id, career_subject_id, status, description) VALUES (?, ?, ?, ?) ON DUPLICATE KEY UPDATE status = ?, description = ?;`
	mock.ExpectExec(q).
		WithArgs(1, 2, "PENDIENTE", nil, "PENDIENTE", nil).
		WillReturnError(nil).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	// When
	err = storage_.UpdateStudentSubject(UpdateStudentSubjectRequest{
		StudentEmail: "example@gmail.com",
		CareerID:     "1",
		SubjectID:    "1",
		Status:       "PENDIENTE",
		Description:  nil,
	})

	if err != nil {
		t.Fatal(err)
	}

	// Then
	require.Nil(t, err)
}

func TestStorage_UpdateStudentSubject_BeginTxError(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("could not start sql mock: %v", err)
	}

	defer db.Close()

	storage_ := NewStorage(sqlx.NewDb(db, ""))

	mock.ExpectBegin().WillReturnError(errors.New("error"))

	// When
	err = storage_.UpdateStudentSubject(UpdateStudentSubjectRequest{
		StudentEmail: "example@gmail.com",
		CareerID:     "1",
		SubjectID:    "1",
		Status:       "PENDIENTE",
		Description:  nil,
	})

	if err == nil {
		t.Fatal("test must fail")
	}

	// Then
	require.EqualError(t, err, "could not begin tx: error")
}

func TestStorage_UpdateStudentSubject_GetStudentEmailError(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("could not start sql mock: %v", err)
	}

	defer db.Close()

	storage_ := NewStorage(sqlx.NewDb(db, ""))

	mock.ExpectBegin()

	q := `SELECT id FROM student WHERE email = ?;`
	mock.ExpectQuery(q).
		WithArgs("example@gmail.com").
		WillReturnError(errors.New("error"))

	// When
	err = storage_.UpdateStudentSubject(UpdateStudentSubjectRequest{
		StudentEmail: "example@gmail.com",
		CareerID:     "1",
		SubjectID:    "1",
		Status:       "PENDIENTE",
		Description:  nil,
	})

	if err == nil {
		t.Fatal("test must fail")
	}

	// Then
	require.EqualError(t, err, "error")
}

func TestStorage_UpdateStudentSubject_GetStudentEmailNotFoundError(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("could not start sql mock: %v", err)
	}

	defer db.Close()

	storage_ := NewStorage(sqlx.NewDb(db, ""))

	mock.ExpectBegin()

	q := `SELECT id FROM student WHERE email = ?;`
	mock.ExpectQuery(q).
		WithArgs("example@gmail.com").
		WillReturnError(sql.ErrNoRows)

	// When
	err = storage_.UpdateStudentSubject(UpdateStudentSubjectRequest{
		StudentEmail: "example@gmail.com",
		CareerID:     "1",
		SubjectID:    "1",
		Status:       "PENDIENTE",
		Description:  nil,
	})

	if err == nil {
		t.Fatal("test must fail")
	}

	// Then
	require.EqualError(t, err, "could not find student: storage: resource not found")
}

func TestStorage_UpdateStudentSubject_CheckStudentAssignedToCareerError(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("could not start sql mock: %v", err)
	}

	defer db.Close()

	storage_ := NewStorage(sqlx.NewDb(db, ""))

	mock.ExpectBegin()

	q := `SELECT id FROM student WHERE email = ?;`
	mock.ExpectQuery(q).
		WithArgs("example@gmail.com").
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	q = `SELECT COUNT(1) FROM student_career WHERE student_id = ? AND career_id = ?`
	mock.ExpectQuery(q).
		WithArgs(1, "1").
		WillReturnError(errors.New("error"))

	// When
	err = storage_.UpdateStudentSubject(UpdateStudentSubjectRequest{
		StudentEmail: "example@gmail.com",
		CareerID:     "1",
		SubjectID:    "1",
		Status:       "PENDIENTE",
		Description:  nil,
	})

	if err == nil {
		t.Fatal("test must fail")
	}

	// Then
	require.EqualError(t, err, "error")
}

func TestStorage_UpdateStudentSubject_CheckStudentAssignedToCareerNotFoundError(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("could not start sql mock: %v", err)
	}

	defer db.Close()

	storage_ := NewStorage(sqlx.NewDb(db, ""))

	mock.ExpectBegin()

	q := `SELECT id FROM student WHERE email = ?;`
	mock.ExpectQuery(q).
		WithArgs("example@gmail.com").
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	q = `SELECT COUNT(1) FROM student_career WHERE student_id = ? AND career_id = ?`
	mock.ExpectQuery(q).
		WithArgs(1, "1").
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"COUNT(1)"}).AddRow(0))

	// When
	err = storage_.UpdateStudentSubject(UpdateStudentSubjectRequest{
		StudentEmail: "example@gmail.com",
		CareerID:     "1",
		SubjectID:    "1",
		Status:       "PENDIENTE",
		Description:  nil,
	})

	if err == nil {
		t.Fatal("test must fail")
	}

	// Then
	require.EqualError(t, err, "could not find student assigned to career: storage: resource not found")
}

func TestStorage_UpdateStudentSubject_GetCareerSubjectByIDsError(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("could not start sql mock: %v", err)
	}

	defer db.Close()

	storage_ := NewStorage(sqlx.NewDb(db, ""))

	mock.ExpectBegin()

	q := `SELECT id FROM student WHERE email = ?;`
	mock.ExpectQuery(q).
		WithArgs("example@gmail.com").
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	q = `SELECT COUNT(1) FROM student_career WHERE student_id = ? AND career_id = ?`
	mock.ExpectQuery(q).
		WithArgs(1, "1").
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"COUNT(1)"}).AddRow(1))

	q = `SELECT id FROM career_subject WHERE career_id = ? AND subject_id = ?`
	mock.ExpectQuery(q).
		WithArgs("1", "1").
		WillReturnError(errors.New("error"))

	mock.ExpectCommit()

	// When
	err = storage_.UpdateStudentSubject(UpdateStudentSubjectRequest{
		StudentEmail: "example@gmail.com",
		CareerID:     "1",
		SubjectID:    "1",
		Status:       "PENDIENTE",
		Description:  nil,
	})

	if err == nil {
		t.Fatal("test must fail")
	}

	// Then
	require.EqualError(t, err, "error")
}

func TestStorage_UpdateStudentSubject_GetCareerSubjectByIDsNotFoundError(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("could not start sql mock: %v", err)
	}

	defer db.Close()

	storage_ := NewStorage(sqlx.NewDb(db, ""))

	mock.ExpectBegin()

	q := `SELECT id FROM student WHERE email = ?;`
	mock.ExpectQuery(q).
		WithArgs("example@gmail.com").
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	q = `SELECT COUNT(1) FROM student_career WHERE student_id = ? AND career_id = ?`
	mock.ExpectQuery(q).
		WithArgs(1, "1").
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"COUNT(1)"}).AddRow(1))

	q = `SELECT id FROM career_subject WHERE career_id = ? AND subject_id = ?`
	mock.ExpectQuery(q).
		WithArgs("1", "1").
		WillReturnError(sql.ErrNoRows)

	mock.ExpectCommit()

	// When
	err = storage_.UpdateStudentSubject(UpdateStudentSubjectRequest{
		StudentEmail: "example@gmail.com",
		CareerID:     "1",
		SubjectID:    "1",
		Status:       "PENDIENTE",
		Description:  nil,
	})

	if err == nil {
		t.Fatal("test must fail")
	}

	// Then
	require.EqualError(t, err, "could not find career and subject: storage: resource not found")
}

func TestStorage_UpdateStudentSubject_UpdateStudentSubjectError(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("could not start sql mock: %v", err)
	}

	defer db.Close()

	storage_ := NewStorage(sqlx.NewDb(db, ""))

	mock.ExpectBegin()

	q := `SELECT id FROM student WHERE email = ?;`
	mock.ExpectQuery(q).
		WithArgs("example@gmail.com").
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	q = `SELECT COUNT(1) FROM student_career WHERE student_id = ? AND career_id = ?`
	mock.ExpectQuery(q).
		WithArgs(1, "1").
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"COUNT(1)"}).AddRow(1))

	q = `SELECT id FROM career_subject WHERE career_id = ? AND subject_id = ?`
	mock.ExpectQuery(q).
		WithArgs("1", "1").
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))

	q = `INSERT INTO student_career_subject (student_id, career_subject_id, status, description) VALUES (?, ?, ?, ?) ON DUPLICATE KEY UPDATE status = ?, description = ?;`
	mock.ExpectExec(q).
		WithArgs(1, 2, "PENDIENTE", nil, "PENDIENTE", nil).
		WillReturnError(errors.New("error"))

	mock.ExpectCommit()

	// When
	err = storage_.UpdateStudentSubject(UpdateStudentSubjectRequest{
		StudentEmail: "example@gmail.com",
		CareerID:     "1",
		SubjectID:    "1",
		Status:       "PENDIENTE",
		Description:  nil,
	})

	if err == nil {
		t.Fatal("test must fail")
	}

	// Then
	require.EqualError(t, err, "error")
}

func TestStorage_UpdateStudentSubject_CommitTxError(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("could not start sql mock: %v", err)
	}

	defer db.Close()

	storage_ := NewStorage(sqlx.NewDb(db, ""))

	mock.ExpectBegin()

	q := `SELECT id FROM student WHERE email = ?;`
	mock.ExpectQuery(q).
		WithArgs("example@gmail.com").
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	q = `SELECT COUNT(1) FROM student_career WHERE student_id = ? AND career_id = ?`
	mock.ExpectQuery(q).
		WithArgs(1, "1").
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"COUNT(1)"}).AddRow(1))

	q = `SELECT id FROM career_subject WHERE career_id = ? AND subject_id = ?`
	mock.ExpectQuery(q).
		WithArgs("1", "1").
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))

	q = `INSERT INTO student_career_subject (student_id, career_subject_id, status, description) VALUES (?, ?, ?, ?) ON DUPLICATE KEY UPDATE status = ?, description = ?;`
	mock.ExpectExec(q).
		WithArgs(1, 2, "PENDIENTE", nil, "PENDIENTE", nil).
		WillReturnError(nil).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit().WillReturnError(errors.New("error"))

	// When
	err = storage_.UpdateStudentSubject(UpdateStudentSubjectRequest{
		StudentEmail: "example@gmail.com",
		CareerID:     "1",
		SubjectID:    "1",
		Status:       "PENDIENTE",
		Description:  nil,
	})

	if err == nil {
		t.Fatal("test must fail")
	}

	// Then
	require.EqualError(t, err, "could not commit tx: error")
}

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
       IFNULL(scs.status, 'PENDIENTE') status,
       scs.description
FROM student AS st
         INNER JOIN career_subject cs ON cs.career_id = ?
         INNER JOIN subject s on s.id = cs.subject_id
         LEFT JOIN student_career_subject scs ON scs.student_id = st.id AND scs.career_subject_id = cs.id
WHERE st.email = ?`
	mock.ExpectPrepare(q).WillReturnError(nil)
	mock.ExpectQuery(q).
		WithArgs("1", "example@gmail.com").
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
       IFNULL(scs.status, 'PENDIENTE') status,
       scs.description
FROM student AS st
         INNER JOIN career_subject cs ON cs.career_id = ?
         INNER JOIN subject s on s.id = cs.subject_id
         LEFT JOIN student_career_subject scs ON scs.student_id = st.id AND scs.career_subject_id = cs.id
WHERE st.email = ?`
	mock.ExpectPrepare(q).WillReturnError(nil)
	mock.ExpectQuery(q).
		WithArgs("1", "example@gmail.com").
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
       IFNULL(scs.status, 'PENDIENTE') status,
       scs.description
FROM student AS st
         INNER JOIN career_subject cs ON cs.career_id = ?
         INNER JOIN subject s on s.id = cs.subject_id
         LEFT JOIN student_career_subject scs ON scs.student_id = st.id AND scs.career_subject_id = cs.id
WHERE st.email = ?`
	mock.ExpectPrepare(q).WillReturnError(nil)
	mock.ExpectQuery(q).
		WithArgs("1", "example@gmail.com").
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
       IFNULL(scs.status, 'PENDIENTE') status,
       scs.description
FROM student AS st
         INNER JOIN career_subject cs ON cs.career_id = ?
         INNER JOIN subject s on s.id = cs.subject_id
         LEFT JOIN student_career_subject scs ON scs.student_id = st.id AND scs.career_subject_id = cs.id
WHERE st.email = ?`
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
       IFNULL(scs.status, 'PENDIENTE') status,
       scs.description
FROM student AS st
         INNER JOIN career_subject cs ON cs.career_id = ?
         INNER JOIN subject s on s.id = cs.subject_id
         LEFT JOIN student_career_subject scs ON scs.student_id = st.id AND scs.career_subject_id = cs.id
WHERE st.email = ?`
	mock.ExpectPrepare(q).WillReturnError(nil)
	mock.ExpectQuery(q).
		WithArgs("1", "example@gmail.com").
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
	require.Equal(t, 240, *subject.Hours)
	require.Equal(t, 8, *subject.Points)
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

func TestStorage_GetProfessorships(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("could not start sql mock: %v", err)
	}

	defer db.Close()

	storage_ := NewStorage(sqlx.NewDb(db, ""))

	q := `SELECT p.name, s.day, s.start, s.end
FROM professorship p
         INNER JOIN schedule s on p.id = s.professorship_id
         INNER JOIN career_subject cs on p.career_subject_id = cs.id
WHERE cs.subject_id = ? AND cs.career_id = ?
ORDER BY day;`
	mock.ExpectPrepare(q).WillReturnError(nil)
	mock.ExpectQuery(q).
		WithArgs("1", "2").
		WillReturnError(nil).
		WillReturnRows(
			sqlmock.NewRows([]string{"name", "day", "start", "end"}).
				AddRow("Professorship 1", 1, "17:00:00", "21:00:00"))

	// When
	professorships, err := storage_.GetProfessorships("1", "2")
	if err != nil {
		t.Fatal(err)
	}

	// Then
	require.Len(t, professorships, 1)

	for _, professorship := range professorships {
		require.Equal(t, "Professorship 1", professorship.Name)
		require.Equal(t, 1, professorship.Day)
		require.Equal(t, "17:00:00", professorship.Start)
		require.Equal(t, "21:00:00", professorship.End)
	}
}

func TestStorage_GetProfessorships_PrepareStmtError(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("could not start sql mock: %v", err)
	}

	defer db.Close()

	storage_ := NewStorage(sqlx.NewDb(db, ""))

	q := `SELECT p.name, s.day, s.start, s.end
FROM professorship p
         INNER JOIN schedule s on p.id = s.professorship_id
         INNER JOIN career_subject cs on p.career_subject_id = cs.id
WHERE cs.subject_id = ? AND cs.career_id = ?
ORDER BY day;`
	mock.ExpectPrepare(q).WillReturnError(errors.New("error"))

	// When
	_, err = storage_.GetProfessorships("1", "2")
	if err == nil {
		t.Fatal("test must fail")
	}

	// Then
	require.EqualError(t, err, "error")
}

func TestStorage_GetProfessorships_ExecuteStmtError(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("could not start sql mock: %v", err)
	}

	defer db.Close()

	storage_ := NewStorage(sqlx.NewDb(db, ""))

	q := `SELECT p.name, s.day, s.start, s.end
FROM professorship p
         INNER JOIN schedule s on p.id = s.professorship_id
         INNER JOIN career_subject cs on p.career_subject_id = cs.id
WHERE cs.subject_id = ? AND cs.career_id = ?
ORDER BY day;`
	mock.ExpectPrepare(q).WillReturnError(nil)
	mock.ExpectQuery(q).WillReturnError(errors.New("error"))

	// When
	_, err = storage_.GetProfessorships("1", "2")
	if err == nil {
		t.Fatal("test must fail")
	}

	// Then
	require.EqualError(t, err, "error")
}
